package model

import (
	common "logger/model/proto/common/v1"
	logs "logger/model/proto/logs/v1"
)


type Process struct {
	ServiceName          string     `protobuf:"bytes,1,opt,name=service_name,json=serviceName,proto3" json:"service_name,omitempty"`
	Tags                 []*common.KeyValue `protobuf:"bytes,2,rep,name=tags,proto3" json:"tags"`
}

type ValueType int32

const (
	ValueType_STRING  ValueType = 0
	ValueType_BOOL    ValueType = 1
	ValueType_INT64   ValueType = 2
	ValueType_FLOAT64 ValueType = 3
	ValueType_BINARY  ValueType = 4
)


type KeyValue struct {
	Key                  string    `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	VType                ValueType `protobuf:"varint,2,opt,name=v_type,json=vType,proto3,enum=jaeger.api_v2.ValueType" json:"v_type,omitempty"`
	VStr                 string    `protobuf:"bytes,3,opt,name=v_str,json=vStr,proto3" json:"v_str,omitempty"`
	VBool                bool      `protobuf:"varint,4,opt,name=v_bool,json=vBool,proto3" json:"v_bool,omitempty"`
	VInt64               int64     `protobuf:"varint,5,opt,name=v_int64,json=vInt64,proto3" json:"v_int64,omitempty"`
	VFloat64             float64   `protobuf:"fixed64,6,opt,name=v_float64,json=vFloat64,proto3" json:"v_float64,omitempty"`
	VBinary              []byte    `protobuf:"bytes,7,opt,name=v_binary,json=vBinary,proto3" json:"v_binary,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}


type LogRecord struct {
	TimeUnixNano uint64 `protobuf:"fixed64,1,opt,name=time_unix_nano,json=timeUnixNano,proto3" json:"time_unix_nano,omitempty"`
	ObservedTimeUnixNano uint64 `protobuf:"fixed64,11,opt,name=observed_time_unix_nano,json=observedTimeUnixNano,proto3" json:"observed_time_unix_nano,omitempty"`
	SeverityNumber logs.SeverityNumber `protobuf:"varint,2,opt,name=severity_number,json=severityNumber,proto3,enum=opentelemetry.proto.logs.v1.SeverityNumber" json:"severity_number,omitempty"`
	SeverityText string `protobuf:"bytes,3,opt,name=severity_text,json=severityText,proto3" json:"severity_text,omitempty"`
	Body *common.AnyValue `protobuf:"bytes,5,opt,name=body,proto3" json:"body,omitempty"`
	Attributes             []*common.KeyValue `protobuf:"bytes,6,rep,name=attributes,proto3" json:"attributes,omitempty"`
	DroppedAttributesCount uint32          `protobuf:"varint,7,opt,name=dropped_attributes_count,json=droppedAttributesCount,proto3" json:"dropped_attributes_count,omitempty"`
	Flags uint32 `protobuf:"fixed32,8,opt,name=flags,proto3" json:"flags,omitempty"`
	TraceId []byte `protobuf:"bytes,9,opt,name=trace_id,json=traceId,proto3" json:"trace_id,omitempty"`
	SpanId []byte `protobuf:"bytes,10,opt,name=span_id,json=spanId,proto3" json:"span_id,omitempty"`
	Process *Process

}