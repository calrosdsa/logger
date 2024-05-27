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
	queryLogs = `SELECT severity_number,body, time_unix_nano, observed_time_unix_nano,service_name,service_attributes,attributes
	FROM logs`
	// queryLogs = `SELECT severity_number,body, time_unix_nano, observed_time_unix_nano, attributes, process
	// FROM logs`
)

type LogReader struct {
	session cassandra.Session
	logger  *zap.Logger
}

func NewLogReader(
	session cassandra.Session,
	logger *zap.Logger,
) logstore.Reader {
	return &LogReader{
		session: session,
		logger:  logger,
	}
}

func (l *LogReader) GetLogs(ctx context.Context) ([]*model.LogRecord, error) {
	return l.getLogs(ctx)
}

func (l *LogReader) getLogs(ctx context.Context) ([]*model.LogRecord, error) {
	q := l.session.Query(queryLogs)
	i := q.Iter()
	// var dbProcess dbmodel.Process
	var timeUnixNano, observedTimeUnixNano uint64
	var severityNumber uint32
	var body,serviceName string
	var attributes,serviceAttributes []dbmodel.KeyValue
	res := make([]*model.LogRecord, 0)
	for i.Scan(&severityNumber, &body, &timeUnixNano, &observedTimeUnixNano, &serviceName,&serviceAttributes,&attributes) {
		dbLog := dbmodel.LogRecord{
			SeverityNumber:       severityNumber,
			Body:                 body,
			TimeUnixNano:         timeUnixNano,
			ObservedTimeUnixNano: observedTimeUnixNano,
			ServiceName: serviceName,
			ServiceAttributes: serviceAttributes,
			Attributes:           attributes,
		}
		fmt.Println("dbLog---",dbLog)
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
