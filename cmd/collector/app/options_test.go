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
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"logger/cmd/collector/app/flags"
	"logger/cmd/collector/app/processor"
	"logger/model"
	"logger/pkg/metrics"
)

func TestAllOptionSet(t *testing.T) {
	types := []processor.LogFormat{processor.LogFormat("sneh")}
	opts := Options.apply(
		Options.ReportBusy(true),
		Options.BlockingSubmit(true),
		Options.ExtraFormatTypes(types),
		Options.LogFilter(func(span *model.LogRecord) bool { return true }),
		Options.HostMetrics(metrics.NullFactory),
		Options.ServiceMetrics(metrics.NullFactory),
		Options.Logger(zap.NewNop()),
		Options.NumWorkers(5),
		Options.PreProcessSpans(func(spans []*model.LogRecord, tenant string) {}),
		// Options.Sanitizer(func(span *model.LogRecord) *model.LogRecord { return span }),
		Options.QueueSize(10),
		Options.DynQueueSizeWarmup(1000),
		Options.DynQueueSizeMemory(1024),
		Options.PreSave(func(span *model.LogRecord, tenant string) {}),
		Options.CollectorTags(map[string]string{"extra": "tags"}),
		Options.SpanSizeMetricsEnabled(true),
		Options.OnDroppedSpan(func(span *model.LogRecord) {}),
	)
	assert.EqualValues(t, 5, opts.numWorkers)
	assert.EqualValues(t, 10, opts.queueSize)
	assert.EqualValues(t, map[string]string{"extra": "tags"}, opts.collectorTags)
	assert.EqualValues(t, 1000, opts.dynQueueSizeWarmup)
	assert.EqualValues(t, 1024, opts.dynQueueSizeMemory)
	assert.True(t, opts.spanSizeMetricsEnabled)
	assert.NotNil(t, opts.onDroppedSpan)
}

func TestNoOptionsSet(t *testing.T) {
	opts := Options.apply()
	assert.EqualValues(t, flags.DefaultNumWorkers, opts.numWorkers)
	assert.EqualValues(t, 0, opts.queueSize)
	assert.Nil(t, opts.collectorTags)
	assert.False(t, opts.reportBusy)
	assert.False(t, opts.blockingSubmit)
	assert.NotPanics(t, func() { opts.preProcessSpans(nil, "") })
	assert.NotPanics(t, func() { opts.preSave(nil, "") })
	assert.True(t, opts.logFilter(nil))
	// log := model.LogRecord{}
	// assert.EqualValues(t, &span, opts.sanitizer(&span))
	assert.EqualValues(t, 0, opts.dynQueueSizeWarmup)
	assert.False(t, opts.spanSizeMetricsEnabled)
	assert.Nil(t, opts.onDroppedSpan)
}