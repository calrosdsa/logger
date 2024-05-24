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

const (
	queryLogs = `SELECT severity_number,body, time_unix_nano, observed_time_unix_nano, attributes, process
	FROM logs`
)

type LogReader struct{
	session cassandra.Session
	logger *zap.Logger
}

func NewLogReader(
	session cassandra.Session,
	logger *zap.Logger,
)logstore.Reader{
	return &LogReader{
		session: session,
		logger: logger,
	}
}

func (l *LogReader) GetLogs(ctx context.Context)([]*model.LogRecord,error){
	return l.getLogs(ctx)
}

func (l *LogReader) getLogs(ctx context.Context)([]*model.LogRecord,error){
	q := l.session.Query(queryLogs)
	i := q.Iter()
	var dbProcess dbmodel.Process
	var timeUnixNano,observedTimeUnixNano uint64
	var severityNumber uint32
	var body string
	var attributes []dbmodel.KeyValue
	res := make([]*model.LogRecord,0)
	for i.Scan(&severityNumber,&body,&timeUnixNano,&observedTimeUnixNano,&attributes,&dbProcess){
		dbLog := dbmodel.LogRecord{
			SeverityNumber: severityNumber,
			Body: body,
			TimeUnixNano: timeUnixNano,
			ObservedTimeUnixNano: observedTimeUnixNano,
			Attributes: attributes,
			Process: dbProcess,
		}
		logModel,err := dbmodel.ToDomain(&dbLog)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		res = append(res, logModel)
	}

	err := i.Close()
	if err != nil {
		return nil, fmt.Errorf("error reading traces from storage: %w", err)
	}
	return res,nil
}

