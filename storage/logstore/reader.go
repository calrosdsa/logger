package logstore

import (
	"context"
	"logger/model"
	"time"
)

type Reader interface {
	// GetServices(ctx context.Context) ([]string, error)
	GetLogs(ctx context.Context, p LogQueryParameters) ([]*model.LogRecord, error)
	GetServices(ctx context.Context) ([]string, error)
	GetOperations(ctx context.Context, p OperationQueryParameters) ([]Operation, error)
}

// LogQueryParameters contains parameters of a log query.
type LogQueryParameters struct {
	ServiceName   string `json:"service_name"`
	OperationName string `json:"operation_name"`
	// Tags          map[string]string
	StartTimeMin   time.Time `json:"start_time_min"`
	StartTimeMax   time.Time `json:"start_time_max"`
	NumTraces      int       `json:"num_traces"`
	SeverityNumber int       `json:"severity_number"`
	ShouldFetchAll bool    `json:"should_fetch_all"`
}

// OperationQueryParameters contains parameters of query operations, empty spanKind means get operations for all kinds of span.
type OperationQueryParameters struct {
	ServiceName string `json:"service_name"`
}

// Operation contains operation name and span kind
type Operation struct {
	Name string `json:"name"`
}
