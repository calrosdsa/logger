// Copyright (c) 2019 The Jaeger Authors.
// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package app

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	"logger/cmd/collector/app/processor"
	"logger/cmd/collector/app/sanitizer"
	"logger/model"
	"logger/pkg/queue"

	"logger/pkg/tenancy"
	"logger/storage/logstore"
)

const (
	// if this proves to be too low, we can increase it
	maxQueueSize = 1_000_000

	// if the new queue size isn't 20% bigger than the previous one, don't change
	minRequiredChange = 1.2
)

type logProcessor struct {
	queue              *queue.BoundedQueue
	queueResizeMu      sync.Mutex
	// metrics            *SpanProcessorMetrics
	preProcessSpans    ProcessLogs
	filterSpan         FilterLog             // filter is called before the sanitizer but after preProcessSpans
	sanitizer          sanitizer.SanitizeSpan // sanitizer is called before processSpan
	processSpan        ProcessLog
	logger             *zap.Logger
	logWriter         logstore.Writer
	reportBusy         bool
	numWorkers         int
	collectorTags      map[string]string
	dynQueueSizeWarmup uint
	dynQueueSizeMemory uint
	bytesProcessed     atomic.Uint64
	spansProcessed     atomic.Uint64
	stopCh             chan struct{}
}


type queueItem struct {
	queuedTime time.Time
	span       *model.LogRecord
	tenant     string
}

// NewSpanProcessor returns a SpanProcessor that preProcesses, filters, queues, sanitizes, and processes spans.
func NewSpanProcessor(
	logWriter logstore.Writer,
	additional []ProcessLog,
	opts ...Option,
) processor.LogProcessor {
	sp := newLogProcessor(logWriter, additional, opts...)
	log.Println("LOGPROCESSOR",sp)
	sp.queue.StartConsumers(sp.numWorkers, func(item interface{}) {
		value := item.(*queueItem)
		sp.processItemFromQueue(value)
	})

	// sp.background(1*time.Second, sp.updateGauges)

	// if sp.dynQueueSizeMemory > 0 {
	// 	sp.background(1*time.Minute, sp.updateQueueSize)
	// }

	return sp
}

func newLogProcessor(logWriter logstore.Writer, additional []ProcessLog, opts ...Option) *logProcessor {
	options := Options.apply(opts...)
	handlerMetrics := NewLogProcessorMetrics(
		options.serviceMetrics,
		options.hostMetrics,
		options.extraFormatTypes)
	droppedItemHandler := func(item interface{}) {
		handlerMetrics.SpansDropped.Inc(1)
		if options.onDroppedSpan != nil {
			options.onDroppedSpan(item.(*queueItem).span)
		}
	}
	boundedQueue := queue.NewBoundedQueue(options.queueSize, droppedItemHandler)

	sanitizers := sanitizer.NewStandardSanitizers()
	if options.sanitizer != nil {
		sanitizers = append(sanitizers, options.sanitizer)
	}

	sp := logProcessor{
		queue:              boundedQueue,
		// metrics:            handlerMetrics,
		logger:             options.logger,
		preProcessSpans:    options.preProcessSpans,
		filterSpan:         options.logFilter,
		sanitizer:          sanitizer.NewChainedSanitizer(sanitizers...),
		reportBusy:         options.reportBusy,
		numWorkers:         options.numWorkers,
		logWriter:         logWriter,
		collectorTags:      options.collectorTags,
		stopCh:             make(chan struct{}),
		dynQueueSizeMemory: options.dynQueueSizeMemory,
		dynQueueSizeWarmup: options.dynQueueSizeWarmup,
	}

	processLogFuncs := []ProcessLog{options.preSave, sp.saveLog}
	if options.dynQueueSizeMemory > 0 {
		options.logger.Info("Dynamically adjusting the queue size at runtime.",
			zap.Uint("memory-mib", options.dynQueueSizeMemory/1024/1024),
			zap.Uint("queue-size-warmup", options.dynQueueSizeWarmup))
	}
	if options.dynQueueSizeMemory > 0 || options.spanSizeMetricsEnabled {
		// add to processLogFuncs
		processLogFuncs = append(processLogFuncs, sp.countSpan)
	}

	processLogFuncs = append(processLogFuncs, additional...)

	sp.processSpan = ChainedProcessLog(processLogFuncs...)
	return &sp
}

func (sp *logProcessor) Close() error {
	close(sp.stopCh)
	sp.queue.Stop()

	return nil
}

func (sp *logProcessor) saveLog(log *model.LogRecord, tenant string) {
	if nil == log.Process {
		sp.logger.Error("process is empty for the log")
		return
	}
	fmt.Println("SAVE LOG",log)

	ctx := tenancy.WithTenant(context.Background(), tenant)
	if err := sp.logWriter.WriteLog(ctx,log);err != nil {
		sp.logger.Error("Failed to save log", zap.Error(err))
	}
	// if err := sp.logWriter.WriteLog(ctx, log); err != nil {
	// 	sp.logger.Error("Failed to save log", zap.Error(err))
	// 	// sp.metrics.SavedErrBySvc.ReportServiceNameForSpan(log)
	// } else {
	// 	sp.logger.Debug("LogRecord written to the storage by the collector")
	// }
}

func (sp *logProcessor) countSpan(log *model.LogRecord, tenant string) {
	// sp.bytesProcessed.Add(uint64(log.Size))
	sp.spansProcessed.Add(1)
}

func (sp *logProcessor) ProcessLogs(mSpans []*model.LogRecord, options processor.LogOptions) ([]bool, error) {
	sp.preProcessSpans(mSpans, options.Tenant)
	// sp.metrics.BatchSize.Update(int64(len(mSpans)))
	retMe := make([]bool, len(mSpans))

	// Note: this is not the ideal place to do this because collector tags are added to Process.Tags,
	// and Process can be shared between different spans in the batch, but we no longer know that,
	// the relation is lost upstream and it's impossible in Go to dedupe pointers. But at least here
	// we have a single thread updating all spans that may share the same Process, before concurrency
	// kicks in.
	for _, span := range mSpans {
		sp.addCollectorTags(span)
	}

	for i, mSpan := range mSpans {
		ok := sp.enqueueSpan(mSpan, options.LogFormat, options.InboundTransport, options.Tenant)
		if !ok && sp.reportBusy {
			return nil, processor.ErrBusy
		}
		retMe[i] = ok
	}
	return retMe, nil
}

func (sp *logProcessor) processItemFromQueue(item *queueItem) {
	sp.processSpan(sp.sanitizer(item.span), item.tenant)
	// sp.metrics.InQueueLatency.Record(time.Since(item.queuedTime))
}

func (sp *logProcessor) addCollectorTags(span *model.LogRecord) {
	if len(sp.collectorTags) == 0 {
		return
	}
	// dedupKey := make(map[string]struct{})
	// for _, tag := range span.Process.Tags {
	// 	if value, ok := sp.collectorTags[tag.Key]; ok && value == tag.AsString() {
	// 		sp.logger.Debug("ignore collector process tags", zap.String("key", tag.Key), zap.String("value", value))
	// 		dedupKey[tag.Key] = struct{}{}
	// 	}
	// }

	// ignore collector tags if has the same key-value in spans
	// for k, v := range sp.collectorTags {
	// 	if _, ok := dedupKey[k]; !ok {
	// 		span.Process.Tags = append(span.Process.Tags, model.String(k, v))
	// 	}
	// }
	// typedTags := model.KeyValues(span.Process.Tags)
	// typedTags.Sort()
}

// Note: spans may share the Process object, so no changes should be made to Process
// in this function as it may cause race conditions.
func (sp *logProcessor) enqueueSpan(span *model.LogRecord, originalFormat processor.LogFormat, transport processor.InboundTransport, tenant string) bool {
	// spanCounts := sp.metrics.GetCountsForFormat(originalFormat, transport)
	// spanCounts.ReceivedBySvc.ReportServiceNameForSpan(span)

	// if !sp.filterSpan(span) {
	// 	spanCounts.RejectedBySvc.ReportServiceNameForSpan(span)
	// 	return true // as in "not dropped", because it's actively rejected
	// }

	// add format tag
	// span.Tags = append(span.Tags, model.String("internal.span.format", string(originalFormat)))

	item := &queueItem{
		queuedTime: time.Now(),
		span:       span,
		tenant:     tenant,
	}
	return sp.queue.Produce(item)
}

func (sp *logProcessor) background(reportPeriod time.Duration, callback func()) {
	go func() {
		ticker := time.NewTicker(reportPeriod)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				callback()
			case <-sp.stopCh:
				return
			}
		}
	}()
}

func (sp *logProcessor) updateQueueSize() {
	if sp.dynQueueSizeWarmup == 0 {
		return
	}

	if sp.dynQueueSizeMemory == 0 {
		return
	}

	if sp.spansProcessed.Load() < uint64(sp.dynQueueSizeWarmup) {
		return
	}

	sp.queueResizeMu.Lock()
	defer sp.queueResizeMu.Unlock()

	// first, we get the average size of a span, by dividing the bytes processed by num of spans
	average := sp.bytesProcessed.Load() / sp.spansProcessed.Load()

	// finally, we divide the available memory by the average size of a span
	idealQueueSize := float64(sp.dynQueueSizeMemory / uint(average))

	// cap the queue size, just to be safe...
	if idealQueueSize > maxQueueSize {
		idealQueueSize = maxQueueSize
	}

	var diff float64
	current := float64(sp.queue.Capacity())
	if idealQueueSize > current {
		diff = idealQueueSize / current
	} else {
		diff = current / idealQueueSize
	}

	// resizing is a costly operation, we only perform it if we are at least n% apart from the desired value
	if diff > minRequiredChange {
		s := int(idealQueueSize)
		sp.logger.Info("Resizing the internal span queue", zap.Int("new-size", s), zap.Uint64("average-span-size-bytes", average))
		sp.queue.Resize(s)
	}
}

func (sp *logProcessor) updateGauges() {
	// sp.metrics.SpansBytes.Update(int64(sp.bytesProcessed.Load()))
	// sp.metrics.QueueLength.Update(int64(sp.queue.Size()))
	// sp.metrics.QueueCapacity.Update(int64(sp.queue.Capacity()))
}
