package logstore

import (
	"context"
	"time"
	"logger/model"

)



type Reader interface {
	// GetServices(ctx context.Context) ([]string, error)
	GetLogs(ctx context.Context) ([]*model.LogRecord,error)
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

// OperationQueryParameters contains parameters of query operations, empty spanKind means get operations for all kinds of span.
type OperationQueryParameters struct {
	ServiceName string
	SpanKind    string
}

// Operation contains operation name and span kind
type Operation struct {
	Name     string
	SpanKind string
}

