// Copyright (c) 2019 The Jaeger Authors.
// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dbmodel

import (
	"bytes"
	// "encoding/binary"

	"logger/model"
	common "logger/model/proto/common/v1"
	logs "logger/model/proto/logs/v1"
)

const (
	childOf     = "child-of"
	followsFrom = "follows-from"

	stringType  = "string"
	boolType    = "bool"
	int64Type   = "int64"
	float64Type = "float64"
	binaryType  = "binary"
)

// TraceID is a serializable form of model.TraceID
type TraceID [16]byte

// Span is the database representation of a span.

type LogRecord struct {
	TimeUnixNano           uint64              `protobuf:"fixed64,1,opt,name=time_unix_nano,json=timeUnixNano,proto3" json:"time_unix_nano,omitempty"`
	ObservedTimeUnixNano   uint64              `protobuf:"fixed64,11,opt,name=observed_time_unix_nano,json=observedTimeUnixNano,proto3" json:"observed_time_unix_nano,omitempty"`
	SeverityNumber         logs.SeverityNumber `protobuf:"varint,2,opt,name=severity_number,json=severityNumber,proto3,enum=opentelemetry.proto.logs.v1.SeverityNumber" json:"severity_number,omitempty"`
	SeverityText           string              `protobuf:"bytes,3,opt,name=severity_text,json=severityText,proto3" json:"severity_text,omitempty"`
	Body                   *common.AnyValue    `protobuf:"bytes,5,opt,name=body,proto3" json:"body,omitempty"`
	Attributes             []*common.KeyValue  `protobuf:"bytes,6,rep,name=attributes,proto3" json:"attributes,omitempty"`
	DroppedAttributesCount uint32              `protobuf:"varint,7,opt,name=dropped_attributes_count,json=droppedAttributesCount,proto3" json:"dropped_attributes_count,omitempty"`
	Flags                  uint32              `protobuf:"fixed32,8,opt,name=flags,proto3" json:"flags,omitempty"`
	TraceId                []byte              `protobuf:"bytes,9,opt,name=trace_id,json=traceId,proto3" json:"trace_id,omitempty"`
	SpanId                 []byte              `protobuf:"bytes,10,opt,name=span_id,json=spanId,proto3" json:"span_id,omitempty"`
	Process                *model.Process
}

// KeyValue is the UDT representation of a Jaeger KeyValue.
type KeyValue struct {
	Key          string  `cql:"key"`
	ValueType    string  `cql:"value_type"`
	ValueString  string  `cql:"value_string"`
	ValueBool    bool    `cql:"value_bool"`
	ValueInt64   int64   `cql:"value_long"`   // using more natural column name for Cassandra
	ValueFloat64 float64 `cql:"value_double"` // using more natural column name for Cassandra
	ValueBinary  []byte  `cql:"value_binary"`
}

// Log is the UDT representation of a Jaeger Log.
type Log struct {
	Timestamp int64      `cql:"ts"` // microseconds since epoch
	Fields    []KeyValue `cql:"fields"`
}

// SpanRef is the UDT representation of a Jaeger Span Reference.
type SpanRef struct {
	RefType string  `cql:"ref_type"`
	TraceID TraceID `cql:"trace_id"`
	SpanID  int64   `cql:"span_id"`
}

// Process is the UDT representation of a Jaeger Process.
type Process struct {
	ServiceName string     `cql:"service_name"`
	Atributtes  []KeyValue `cql:"attributes"`
}

// TagInsertion contains the items necessary to insert a tag for a given span
type TagInsertion struct {
	ServiceName string
	TagKey      string
	TagValue    string
}

func (t TagInsertion) String() string {
	const uniqueTagDelimiter = ":"
	var buffer bytes.Buffer
	buffer.WriteString(t.ServiceName)
	buffer.WriteString(uniqueTagDelimiter)
	buffer.WriteString(t.TagKey)
	buffer.WriteString(uniqueTagDelimiter)
	buffer.WriteString(t.TagValue)
	return buffer.String()
}

// TraceIDFromDomain converts domain TraceID into serializable DB representation.
// func TraceIDFromDomain(traceID model.TraceID) TraceID {
// 	dbTraceID := TraceID{}
// 	binary.BigEndian.PutUint64(dbTraceID[:8], uint64(traceID.High))
// 	binary.BigEndian.PutUint64(dbTraceID[8:], uint64(traceID.Low))
// 	return dbTraceID
// }

// // ToDomain converts trace ID from db-serializable form to domain TradeID
// func (dbTraceID TraceID) ToDomain() model.TraceID {
// 	traceIDHigh := binary.BigEndian.Uint64(dbTraceID[:8])
// 	traceIDLow := binary.BigEndian.Uint64(dbTraceID[8:])
// 	return model.NewTraceID(traceIDHigh, traceIDLow)
// }

// String returns hex string representation of the trace ID.
// func (dbTraceID TraceID) String() string {
// 	return dbTraceID.ToDomain().String()
// }
