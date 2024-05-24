package proto

import (
	"logger/model"
	pbL "logger/model/proto/logs/v1"
	common "logger/model/proto/common/v1"
)

func ToDomainLog(log *pbL.LogRecord, processs *model.Process) *model.LogRecord {
	return toDomain{}.ToDomainLog(log, processs)
}

type toDomain struct{}

func (t toDomain) ToDomainLog(log *pbL.LogRecord, process *model.Process) *model.LogRecord {
	return t.transformToLog(log,process)
}

func (t toDomain) transformToLog(log *pbL.LogRecord, process *model.Process) *model.LogRecord {
	return &model.LogRecord{
		TimeUnixNano:           log.GetTimeUnixNano(),
		ObservedTimeUnixNano:   log.GetObservedTimeUnixNano(),
		SeverityNumber:         log.GetSeverityNumber(),
		SeverityText:           log.GetSeverityText(),
		Body:                   log.GetBody().GetStringValue(),
		Attributes:             t.toDomainAtrributes(log.Attributes),
		DroppedAttributesCount: log.GetDroppedAttributesCount(),
		Flags:                  log.GetFlags(),
		TraceId:                log.GetTraceId(),
		SpanId:                 log.GetSpanId(),
		Process: process,
	}
}

func (t toDomain)toDomainAtrributes(attributes []*common.KeyValue)[]model.KeyValue{
	res := make([]model.KeyValue,len(attributes))
	for _,v:= range attributes {
		res = append(res, model.KeyValue{
			Key: v.Key,
			Value: v.Value,
		})
	}
	return res
}