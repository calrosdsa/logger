// Copyright (c) 2020 The Jaeger Authors.
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

package processor

import (
	"errors"
	"io"

	model "logger/model"
)

// ErrBusy signalizes that processor cannot process incoming data
var ErrBusy = errors.New("server busy")

// LogOptions additional options passed to processor along with the spans.
type LogOptions struct {
	LogFormat        LogFormat
	InboundTransport InboundTransport
	Tenant           string
}

// LogProcessor handles model spans
type LogProcessor interface {
	// ProcessLogs processes model spans and return with either a list of true/false success or an error
	ProcessLogs(mSpans []*model.LogRecord, options LogOptions) ([]bool, error)
	io.Closer
}

// InboundTransport identifies the transport used to receive spans.
type InboundTransport string

const (
	// GRPCTransport indicates spans received over gRPC.
	GRPCTransport InboundTransport = "grpc"
	// HTTPTransport indicates spans received over HTTP.
	HTTPTransport InboundTransport = "http"
	// UnknownTransport is the fallback/catch-all category.
	UnknownTransport InboundTransport = "unknown"
)

// LogFormat identifies the data format in which the span was originally received.
type LogFormat string

const (
	// JaegerSpanFormat is for Jaeger Thrift spans.
	// JaegerSpanFormat LogFormat = "jaeger"
	// ZipkinSpanFormat is for Zipkin Thrift spans.
	// ZipkinSpanFormat LogFormat = "zipkin"
	// ProtoLogFormat is for Jaeger protobuf Spans.
	ProtoLogFormat LogFormat = "proto"
	// OTLPLogFormat is for OpenTelemetry OTLP format.
	OTLPLogFormat LogFormat = "otlp"
	// UnknownSpanFormat is the fallback/catch-all category.
	// UnknownSpanFormat LogFormat = "unknown"
)
