package proto

import (
	"logger/model"
	pbL "logger/model/proto/logs/v1"
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
		Body:                   log.GetBody(),
		Attributes:             log.GetAttributes(),
		DroppedAttributesCount: log.GetDroppedAttributesCount(),
		Flags:                  log.GetFlags(),
		TraceId:                log.GetTraceId(),
		SpanId:                 log.GetSpanId(),
		Process: process,
	}
}
