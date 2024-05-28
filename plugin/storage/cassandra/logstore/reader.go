package logstore

import (
	"context"
	"errors"
	"fmt"
	"logger/pkg/cassandra"
	"logger/plugin/storage/cassandra/logstore/dbmodel"

	"logger/model"
	"logger/storage/logstore"

	"go.uber.org/zap"
)

// attributes
const (
	queryLogs = `SELECT severity_number,body, start_time, observed_time_unix_nano,service_name,operation_name,service_attributes,attributes
	FROM logs where service_name = ? and operation_name = ?  AND start_time > ? AND start_time < ? LIMIT ?`
	defaultNumTraces = 100
	// queryLogs = `SELECT severity_number,body, start_time, observed_time_unix_nano, attributes, process
	// FROM logs`
)

var (
	// ErrServiceNameNotSet occurs when attempting to query with an empty service name
	ErrServiceNameNotSet = errors.New("service Name must be set")

	// ErrOperationNameNotSet occurs when attempting to query with an empty service name
	ErrOperationNameNotSet = errors.New("operation Name must be set")

	// ErrStartTimeMinGreaterThanMax occurs when start time min is above start time max
	ErrStartTimeMinGreaterThanMax = errors.New("start Time Minimum is above Maximum")

	// ErrMalformedRequestObject occurs when a request object is nil
	ErrMalformedRequestObject = errors.New("malformed request object")

	// ErrDurationAndTagQueryNotSupported occurs when duration and tags are both set
	// ErrDurationAndTagQueryNotSupported = errors.New("cannot query for duration and tags simultaneously")

	// ErrStartAndEndTimeNotSet occurs when start time and end time are not set
	ErrStartAndEndTimeNotSet = errors.New("start and End Time must be set")
)

type serviceNamesReader func() ([]string, error)

type operationNamesReader func(query logstore.OperationQueryParameters) ([]logstore.Operation, error)

type LogReader struct {
	session              cassandra.Session
	logger               *zap.Logger
	serviceNamesReader   serviceNamesReader
	operationNamesReader operationNamesReader
}

func NewLogReader(
	session cassandra.Session,
	logger *zap.Logger,
) logstore.Reader {
	serviceNamesStorage := NewServiceNamesStorage(session, 0, logger)
	operationNamesStorage := NewOperationNamesStorage(session, 0, logger)
	return &LogReader{
		session:              session,
		logger:               logger,
		serviceNamesReader:   serviceNamesStorage.GetServices,
		operationNamesReader: operationNamesStorage.GetOperations,
	}
}

func (l *LogReader) GetServices(ctx context.Context) ([]string, error) {
	return l.serviceNamesReader()
}

func (l *LogReader) GetOperations(ctx context.Context, p logstore.OperationQueryParameters) ([]logstore.Operation, error) {
	return l.operationNamesReader(p)
}


func (l *LogReader) GetLogs(ctx context.Context, p logstore.LogQueryParameters) ([]*model.LogRecord, error) {

	return l.getLogs(ctx, p)
}

func (l *LogReader) getLogs(ctx context.Context, p logstore.LogQueryParameters) ([]*model.LogRecord, error) {
	if err := validateQuery(&p); err != nil {
		return nil, err
	}
	if p.NumTraces == 0 {
		p.NumTraces = defaultNumTraces
	}
	// query := l.buildQuery(p)
	q := l.session.Query(queryLogs,
		p.ServiceName,
		p.OperationName,
		model.TimeAsEpochMicroseconds(p.StartTimeMin),
		model.TimeAsEpochMicroseconds(p.StartTimeMax),
		p.NumTraces,
	)
	i := q.Iter()
	var timeUnixNano, observedTimeUnixNano uint64
	var severityNumber uint32
	var body, serviceName, methodName string
	var attributes, serviceAttributes []dbmodel.KeyValue
	res := make([]*model.LogRecord, 0)
	for i.Scan(&severityNumber, &body, &timeUnixNano, &observedTimeUnixNano, &serviceName, &methodName, &serviceAttributes, &attributes) {
		dbLog := dbmodel.LogRecord{
			SeverityNumber:       severityNumber,
			Body:                 body,
			TimeUnixNano:         timeUnixNano,
			ObservedTimeUnixNano: observedTimeUnixNano,
			ServiceName:          serviceName,
			OperationName:        methodName,
			ServiceAttributes:    serviceAttributes,
			Attributes:           attributes,
		}
		logModel, err := dbmodel.ToDomain(&dbLog)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		res = append(res, logModel)
	}

	err := i.Close()
	if err != nil {
		return nil, fmt.Errorf("error reading logs from storage: %w", err)
	}
	return res, nil
}


func validateQuery(p *logstore.LogQueryParameters) error {
	if p == nil {
		return ErrMalformedRequestObject
	}
	if p.ServiceName == "" {
		return ErrServiceNameNotSet
	}
	if p.OperationName == "" {
		return ErrOperationNameNotSet
	}
	if p.StartTimeMin.IsZero() || p.StartTimeMax.IsZero() {
		return ErrStartAndEndTimeNotSet
	}
	if !p.StartTimeMin.IsZero() && !p.StartTimeMax.IsZero() && p.StartTimeMax.Before(p.StartTimeMin) {
		return ErrStartTimeMinGreaterThanMax
	}
	return nil
}

// func (l *LogReader) buildQuery( p logstore.LogQueryParameters) string {
// 	var partitionQuery string
// 	// if p.ShouldFetchAll {
// 	// 	partitionQuery = ""
// 	// }else{
// 	partitionQuery = fmt.Sprintf("service_name = '%s' and operation_name = '%s' and",p.ServiceName,p.OperationName)
// 	// }
// 	// if p.SeverityNumber != 0 {
// 	// 	severityNumberQuery = fmt.Sprintf("and severity_number = %d",p.SeverityNumber)
// 	// }
// 	query := fmt.Sprintf(`SELECT severity_number,body, start_time, observed_time_unix_nano,service_name,
// 	operation_name,service_attributes,attributes FROM logs where %s 
// 	start_time > ? AND start_time < ? limit ?`,partitionQuery)
// 	fmt.Println("QUERY BUILDER",query)
// 	return query
// }