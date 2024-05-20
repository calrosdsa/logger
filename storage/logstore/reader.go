package logstore

import (
	"context"
	"time"
	"logger/model/proto"

)



type Reader interface {
	GetServices(ctx context.Context) ([]string, error)
	GetLogs(ctx context.Context) ([]proto.LogRecord,error)
}

// LogQueryParameters contains parameters of a log query.
type TraceQueryParameters struct {
	ServiceName   string
	OperationName string
	// Tags          map[string]string
	StartTimeMin  time.Time
	StartTimeMax  time.Time
	NumTraces     int
}


