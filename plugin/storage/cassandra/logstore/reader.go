package logstore

import (
	"context"
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
	FROM logs where service_name = ? and operation_name = ?  AND start_time > ? AND start_time < ?`
	// queryLogs = `SELECT severity_number,body, start_time, observed_time_unix_nano, attributes, process
	// FROM logs`
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
// 1716833724638359100

// 1716708050520000

func (l *LogReader) getLogs(ctx context.Context, p logstore.LogQueryParameters) ([]*model.LogRecord, error) {
	fmt.Println("TIME TO EPOCH",model.TimeAsEpochMicroseconds(p.StartTimeMin),model.TimeAsEpochMicroseconds(p.StartTimeMin))
	q := l.session.Query(queryLogs,p.ServiceName,p.OperationName,
	model.TimeAsEpochMicroseconds(p.StartTimeMin),
	model.TimeAsEpochMicroseconds(p.StartTimeMax),
)
	i := q.Iter()
	// var dbProcess dbmodel.Process
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
		fmt.Println("dbLog---", dbLog)
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
