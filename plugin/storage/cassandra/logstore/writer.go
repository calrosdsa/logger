package logstore

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"go.uber.org/zap"

	"logger/model"
	"logger/pkg/cassandra"
	casMetrics "logger/pkg/cassandra/metrics"
	"logger/pkg/metrics"
	"logger/plugin/storage/cassandra/logstore/dbmodel"
)

const (
	insertSpan = `
		INSERT
		INTO traces(trace_id, span_id, span_hash, parent_id, operation_name, flags,
				    start_time, duration, tags, logs, refs, process)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	serviceNameIndex = `
		INSERT
		INTO service_name_index(service_name, bucket, start_time, trace_id)
		VALUES (?, ?, ?, ?)`

	serviceOperationIndex = `
		INSERT
		INTO
		service_operation_index(service_name, operation_name, start_time, trace_id)
		VALUES (?, ?, ?, ?)`

	tagIndex = `
		INSERT
		INTO tag_index(trace_id, span_id, service_name, start_time, tag_key, tag_value)
		VALUES (?, ?, ?, ?, ?, ?)`

	durationIndex = `
		INSERT
		INTO duration_index(service_name, operation_name, bucket, duration, start_time, trace_id)
		VALUES (?, ?, ?, ?, ?, ?)`

	maximumTagKeyOrValueSize = 256

	// DefaultNumBuckets Number of buckets for bucketed keys
	defaultNumBuckets = 10

	durationBucketSize = time.Hour
)

const (
	storeFlag = storageMode(1 << iota)
	indexFlag
)

type (
	storageMode          uint8
	serviceNamesWriter   func(serviceName string) error
	// operationNamesWriter func(operation dbmodel.Operation) error
)

type spanWriterMetrics struct {
	traces                *casMetrics.Table
	tagIndex              *casMetrics.Table
	serviceNameIndex      *casMetrics.Table
	serviceOperationIndex *casMetrics.Table
	durationIndex         *casMetrics.Table
}

// LogWriter handles all writes to Cassandra for the Jaeger data model
type LogWriter struct {
	session              cassandra.Session
	serviceNamesWriter   serviceNamesWriter
	// operationNamesWriter operationNamesWriter
	writerMetrics        spanWriterMetrics
	logger               *zap.Logger
	// tagIndexSkipped      metrics.Counter
	// tagFilter            dbmodel.TagFilter
	storageMode storageMode
	indexFilter dbmodel.IndexFilter
}

// NewLogWriter returns a LogWriter
func NewLogWriter(
	session cassandra.Session,
	writeCacheTTL time.Duration,
	metricsFactory metrics.Factory,
	logger *zap.Logger,
	options ...Option,
) *LogWriter {
	serviceNamesStorage := NewServiceNamesStorage(session, writeCacheTTL, metricsFactory, logger)
	// operationNamesStorage := NewOperationNamesStorage(session, writeCacheTTL, metricsFactory, logger)
	// tagIndexSkipped := metricsFactory.Counter(metrics.Options{Name: "tag_index_skipped", Tags: nil})
	opts := applyOptions(options...)
	return &LogWriter{
		session:              session,
		serviceNamesWriter:   serviceNamesStorage.Write,
		// operationNamesWriter: operationNamesStorage.Write,
		writerMetrics: spanWriterMetrics{
			traces:                casMetrics.NewTable(metricsFactory, "traces"),
			tagIndex:              casMetrics.NewTable(metricsFactory, "tag_index"),
			serviceNameIndex:      casMetrics.NewTable(metricsFactory, "service_name_index"),
			serviceOperationIndex: casMetrics.NewTable(metricsFactory, "service_operation_index"),
			durationIndex:         casMetrics.NewTable(metricsFactory, "duration_index"),
		},
		logger:          logger,
		// tagIndexSkipped: tagIndexSkipped,
		// tagFilter:       opts.tagFilter,
		storageMode: opts.storageMode,
		indexFilter: opts.indexFilter,
	}
}

// Close closes LogWriter
func (s *LogWriter) Close() error {
	s.session.Close()
	return nil
}

// WriteLog saves the span into Cassandra
func (s *LogWriter) WriteLog(ctx context.Context, span *model.LogRecord) error {
	ds := dbmodel.FromDomain(span)
	if s.storageMode&storeFlag == storeFlag {
		if err := s.writeSpan(span, ds); err != nil {
			return err
		}
	}
	if s.storageMode&indexFlag == indexFlag {
		if err := s.writeIndexes(span, ds); err != nil {
			return err
		}
	}
	return nil
}

func (s *LogWriter) writeSpan(span *model.LogRecord, ds *dbmodel.LogRecord) error {
	fmt.Println("WRITE SPAN",span)
	// mainQuery := s.session.Query(
	// 	insertSpan,
	// 	ds.TraceID,
	// 	ds.SpanID,
	// 	ds.SpanHash,
	// 	ds.ParentID,
	// 	ds.OperationName,
	// 	ds.Flags,
	// 	ds.StartTime,
	// 	ds.Duration,
	// 	ds.Tags,
	// 	ds.Logs,
	// 	ds.Refs,
	// 	ds.Process,
	// )
	// if err := s.writerMetrics.traces.Exec(mainQuery, s.logger); err != nil {
	// 	return s.logError(ds, err, "Failed to insert span", s.logger)
	// }
	return nil
}

func (s *LogWriter) writeIndexes(span *model.LogRecord, ds *dbmodel.LogRecord) error {
	// spanKind, _ := span.GetSpanKind()
	// if err := s.saveServiceNameAndOperationName(dbmodel.Operation{
	// 	ServiceName:   ds.ServiceName,
	// 	SpanKind:      spanKind.String(),
	// 	OperationName: ds.OperationName,
	// }); err != nil {
	// 	// should this be a soft failure?
	// 	return s.logError(ds, err, "Failed to insert service name and operation name", s.logger)
	// }

	// if s.indexFilter(ds, dbmodel.ServiceIndex) {
	// 	if err := s.indexByService(ds); err != nil {
	// 		return s.logError(ds, err, "Failed to index service name", s.logger)
	// 	}
	// }

	// if s.indexFilter(ds, dbmodel.OperationIndex) {
	// 	if err := s.indexByOperation(ds); err != nil {
	// 		return s.logError(ds, err, "Failed to index operation name", s.logger)
	// 	}
	// }

	// if span.Flags.IsFirehoseEnabled() {
	// 	return nil // skipping expensive indexing
	// }

	// if err := s.indexByTags(span, ds); err != nil {
	// 	return s.logError(ds, err, "Failed to index tags", s.logger)
	// }

	// if s.indexFilter(ds, dbmodel.DurationIndex) {
	// 	if err := s.indexByDuration(ds, span.StartTime); err != nil {
	// 		return s.logError(ds, err, "Failed to index duration", s.logger)
	// 	}
	// }
	return nil
}

func (s *LogWriter) indexByTags(span *model.LogRecord, ds *dbmodel.LogRecord) error {
	// for _, v := range dbmodel.GetAllUniqueTags(span, s.tagFilter) {
	// 	// we should introduce retries or just ignore failures imo, retrying each individual tag insertion might be better
	// 	// we should consider bucketing.
	// 	if s.shouldIndexTag(v) {
	// 		insertTagQuery := s.session.Query(tagIndex, ds.TraceID, ds.SpanID, v.ServiceName, ds.StartTime, v.TagKey, v.TagValue)
	// 		if err := s.writerMetrics.tagIndex.Exec(insertTagQuery, s.logger); err != nil {
	// 			withTagInfo := s.logger.
	// 				With(zap.String("tag_key", v.TagKey)).
	// 				With(zap.String("tag_value", v.TagValue)).
	// 				With(zap.String("service_name", v.ServiceName))
	// 			return s.logError(ds, err, "Failed to index tag", withTagInfo)
	// 		}
	// 	} else {
	// 		s.tagIndexSkipped.Inc(1)
	// 	}
	// }
	return nil
}

func (s *LogWriter) indexByDuration(span *dbmodel.LogRecord, startTime time.Time) error {
	// query := s.session.Query(durationIndex)
	// timeBucket := startTime.Round(durationBucketSize)
	// var err error
	// indexByOperationName := func(operationName string) {
	// 	q1 := query.Bind(span.Process.ServiceName, operationName, timeBucket, span.Duration, span.StartTime, span.TraceID)
	// 	if err2 := s.writerMetrics.durationIndex.Exec(q1, s.logger); err2 != nil {
	// 		_ = s.logError(span, err2, "Cannot index duration", s.logger)
	// 		err = err2
	// 	}
	// }
	// indexByOperationName("")                 // index by service name alone
	// indexByOperationName(span.OperationName) // index by service name and operation name
	return nil
}

func (s *LogWriter) indexByService(span *dbmodel.LogRecord) error {
	// bucketNo := uint64(span.SpanHash) % defaultNumBuckets
	// query := s.session.Query(serviceNameIndex)
	// q := query.Bind(span.Process.ServiceName, bucketNo, span.StartTime, span.TraceID)
	// return s.writerMetrics.serviceNameIndex.Exec(q, s.logger)
	return nil
}

func (s *LogWriter) indexByOperation(span *dbmodel.LogRecord) error {
	// query := s.session.Query(serviceOperationIndex)
	// q := query.Bind(span.Process.ServiceName, span.OperationName, span.StartTime, span.TraceID)
	// return s.writerMetrics.serviceOperationIndex.Exec(q, s.logger)
	return nil
}

// shouldIndexTag checks to see if the tag is json or not, if it's UTF8 valid and it's not too large
func (s *LogWriter) shouldIndexTag(tag dbmodel.TagInsertion) bool {
	isJSON := func(s string) bool {
		var js json.RawMessage
		// poor man's string-is-a-json check shortcircuits full unmarshalling
		return strings.HasPrefix(s, "{") && json.Unmarshal([]byte(s), &js) == nil
	}

	return len(tag.TagKey) < maximumTagKeyOrValueSize &&
		len(tag.TagValue) < maximumTagKeyOrValueSize &&
		utf8.ValidString(tag.TagValue) &&
		utf8.ValidString(tag.TagKey) &&
		!isJSON(tag.TagValue)
}

func (s *LogWriter) logError(span *dbmodel.LogRecord, err error, msg string, logger *zap.Logger) error {
	// logger.
	// 	With(zap.String("trace_id", span.TraceID.String())).
	// 	With(zap.Int64("span_id", span.SpanID)).
	// 	With(zap.Error(err)).
	// 	Error(msg)
	return fmt.Errorf("%s: %w", msg, err)
}

// func (s *LogWriter) saveServiceNameAndOperationName(operation dbmodel.Operation) error {
// 	if err := s.serviceNamesWriter(operation.ServiceName); err != nil {
// 		return err
// 	}
// 	return s.operationNamesWriter(operation)
// }
