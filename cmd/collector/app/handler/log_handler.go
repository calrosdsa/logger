package handler

import (
	"logger/cmd/collector/app/processor"
	"logger/model"
	pb "logger/model/proto/v1"
	pbL "logger/model/proto/logs/v1"
	pbC "logger/model/proto/common/v1"

	pConv "logger/model/converter/proto"

	"go.uber.org/zap"
)

// SubmitBatchOptions are passed to Submit methods of the handlers.
type SubmitBatchOptions struct {
	InboundTransport processor.InboundTransport
}

type BatchSubmitResponse struct {
	Ok bool
}

type Batch pb.ExportLogsServiceRequest

// type Batch struct {
// 	Process *Process `thrift:"process,1,required" db:"process" json:"process"`
// 	Spans []*Span	 `thrift:"spans,2,required" db:"spans" json:"spans"`
// 	SeqNo *int64 `thrift:"seqNo,3" db:"seqNo" json:"seqNo,omitempty"`
// 	Stats *ClientStats `thrift:"stats,4" db:"stats" json:"stats,omitempty"`
//   }

type BatchesHandler interface {
	// SubmitBatches records a batch of spans in Jaeger Thrift format
	SubmitBatches(batches []*pbL.ResourceLogs, options SubmitBatchOptions) ([]*BatchSubmitResponse, error)
}

type batchesHandler struct {
	logger         *zap.Logger
	modelProcessor processor.LogProcessor
}

// NewJaegerSpanHandler returns a JaegerBatchesHandler
func NewJaegerLogHandler(logger *zap.Logger, modelProcessor processor.LogProcessor) BatchesHandler {
	return &batchesHandler{
		logger:         logger,
		modelProcessor: modelProcessor,
	}
}

func (h *batchesHandler) SubmitBatches(batches []*pbL.ResourceLogs, opts SubmitBatchOptions) ([]*BatchSubmitResponse, error) {
	responses := make([]*BatchSubmitResponse, 0, len(batches))
	for _, batch := range batches {
		proccess := h.getProcess(batch)

		for _, scopeLog := range batch.GetScopeLogs() {
			mLogs := make([]*model.LogRecord, 0, len(scopeLog.GetLogRecords()))
			for _, log := range scopeLog.GetLogRecords() {
				mLogs = append(mLogs, pConv.ToDomainLog(log, proccess))
			}
			oks, err := h.modelProcessor.ProcessLogs(mLogs, processor.LogOptions{
				InboundTransport: opts.InboundTransport,
				LogFormat:        processor.OTLPLogFormat,
			})

			if err != nil {
				h.logger.Error("Collector failed to process span batch", zap.Error(err))
				return nil, err
			}
			batchOk := true
			for _, ok := range oks {
				if !ok {
					batchOk = false
					break
				}
			}

			h.logger.Debug("Span batch processed by the collector.", zap.Bool("ok", batchOk))
			res := &BatchSubmitResponse{
				Ok: batchOk,
			}
			responses = append(responses, res)
		}
	}
	return responses, nil
}

func (h *batchesHandler) getProcess(r *pbL.ResourceLogs) *model.Process {
	if r == nil {
		return nil
	}
	resources := r.GetResource()
	var serviceName string
	tags := make([]*pbC.KeyValue, 0, len(resources.GetAttributes()))

	for _, attr := range resources.GetAttributes() {
		if attr == nil {
			continue
		}
		if attr.Key == "service.name" {
			serviceName = attr.GetValue().GetStringValue()
		} else {
			tags = append(tags, attr)
		}
	}
	return &model.Process{
		ServiceName: serviceName,
		Tags:        tags,
	}
}
