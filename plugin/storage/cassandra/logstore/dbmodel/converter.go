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
	"fmt"
	"logger/model"
	common "logger/model/proto/common/v1"
	logs "logger/model/proto/logs/v1"


)

const (
	// warningStringPrefix is a magic string prefix for tag names to store span warnings.
	warningStringPrefix = "$$span.warning."
)

var (
	// 	dbToDomainRefMap = map[string]model.LogRecordRefType{
	// 		childOf:     model.LogRecordRefType_CHILD_OF,
	// 		followsFrom: model.LogRecordRefType_FOLLOWS_FROM,
	// 	}

	// 	domainToDBRefMap = map[model.LogRecordRefType]string{
	// 		model.LogRecordRefType_CHILD_OF:     childOf,
	// 		model.LogRecordRefType_FOLLOWS_FROM: followsFrom,
	// 	}

	// domainToDBValueTypeMap = map[model.ValueType]string{
	// 	model.StringType:  stringType,
	// 	model.BoolType:    boolType,
	// 	model.Int64Type:   int64Type,
	// 	model.Float64Type: float64Type,
	// 	model.BinaryType:  binaryType,
	// }
)

// FromDomain converts a domain model.LogRecord to a database LogRecord
func FromDomain(span *model.LogRecord) *LogRecord {
	return converter{}.fromDomain(span)
}

// ToDomain converts a database LogRecord to a domain model.LogRecord
func ToDomain(dbSpan *LogRecord) (*model.LogRecord, error) {
	return converter{}.toDomain(dbSpan)
}

// converter converts Spans between domain and database representations.
// It primarily exists to namespace the conversion functions.
type converter struct{}

func (c converter) fromDomain(log *model.LogRecord) *LogRecord {
	attributes := c.toDBAttributes(log.Attributes)
	process := c.toDBProcess(log.Process)

	return &LogRecord{
		TimeUnixNano:           log.TimeUnixNano,
		ObservedTimeUnixNano:   log.ObservedTimeUnixNano,
		SeverityNumber:         uint32(log.SeverityNumber.Number()),
		SeverityText:           log.SeverityText,
		Body:                   log.Body,
		Attributes:             attributes,
		DroppedAttributesCount: log.DroppedAttributesCount,
		Flags:                  log.Flags,
		TraceId:                log.TraceId,
		SpanId:                 log.SpanId,
		Process:                process,
	}
}

func (c converter) toDomain(log *LogRecord) (*model.LogRecord, error) {
	attributes, err := c.fromDBAttrinutes(log.Attributes)
	if err != nil {
		return nil, err
	}
	process, err := c.fromDBProcess(log.Process)
	if err != nil {
		return nil, err
	}
	span := &model.LogRecord{
		TimeUnixNano:           log.TimeUnixNano,
		ObservedTimeUnixNano:   log.ObservedTimeUnixNano,
		SeverityNumber:         logs.SeverityNumber(log.SeverityNumber),
		SeverityText:           log.SeverityText,
		Body:                   log.Body,
		Attributes:             attributes,
		DroppedAttributesCount: log.DroppedAttributesCount,
		Flags:                  log.Flags,
		TraceId:                log.TraceId,
		SpanId:                 log.SpanId,
		Process:                process,
	}
	return span, nil
}

func (c converter) fromDBAttrinutes(tags []KeyValue) ([]model.KeyValue, error) {
	retMe := make([]model.KeyValue, 0, len(tags))
	for i := range tags {
		kv, err := c.fromDBAttribute(&tags[i])
		if err != nil {
			return nil, err
		}
		retMe = append(retMe, kv)
	}
	return retMe, nil
}

func (c converter) fromDBAttribute(attribute *KeyValue) (model.KeyValue, error) {
	switch attribute.ValueType {
	case stringType:
		return model.KeyValue{
			Key: attribute.Key,
			Value: &common.AnyValue{
				Value: &common.AnyValue_StringValue{
					StringValue: attribute.ValueString,
				},
			},
		}, nil
	case boolType:
		return model.KeyValue{
			Key: attribute.Key,
			Value: &common.AnyValue{
				Value: &common.AnyValue_BoolValue{
					BoolValue: attribute.ValueBool,
				},
			},
		}, nil
	case int64Type:
		return model.KeyValue{
			Key: attribute.Key,
			Value: &common.AnyValue{
				Value: &common.AnyValue_IntValue{
					IntValue: attribute.ValueInt64,
				},
			},
		}, nil
	case float64Type:
		return model.KeyValue{
			Key: attribute.Key,
			Value: &common.AnyValue{
				Value: &common.AnyValue_DoubleValue{
					DoubleValue: attribute.ValueFloat64,
				},
			},
		}, nil
	case binaryType:
		return model.KeyValue{
			Key: attribute.Key,
			Value: &common.AnyValue{
				Value: &common.AnyValue_BytesValue{
					BytesValue: attribute.ValueBinary,
				},
			},
		}, nil
	}
	return model.KeyValue{}, fmt.Errorf("invalid ValueType in %+v", attribute)
}



func (c converter) fromDBProcess(process Process) (*model.Process, error) {
	attributes, err := c.fromDBAttrinutes(process.Atributtes)
	if err != nil {
		return nil, err
	}
	return &model.Process{
		Attributes:        attributes,
		ServiceName: process.ServiceName,
	}, nil
}

func (c converter) toDBAttributes(attributes []model.KeyValue) []KeyValue {
	retMe := make([]KeyValue, len(attributes))
	for i, t := range attributes {
		// do we want to validate a jaeger tag here? Making sure that the type and value matches up?
		retMe[i] = KeyValue{
			Key:          t.Key,
			ValueType:    t.GetTypeValues(),
			ValueString:  t.Value.GetStringValue(),
			ValueBool:    t.Value.GetBoolValue(),
			ValueInt64:   t.Value.GetIntValue(),
			ValueFloat64: t.Value.GetDoubleValue(),
			ValueBinary:  t.Value.GetBytesValue(),
		}
	}
	return retMe
}


func (c converter) toDBProcess(process *model.Process) Process {
	return Process{
		ServiceName: process.ServiceName,
		Atributtes:        c.toDBAttributes(process.Attributes),
	}
}
