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

package app

import (
	"logger/model"
)

// ProcessLog processes a Domain Model Span
type ProcessLog func(log *model.LogRecord, tenant string)

// ProcessLogs processes a batch of Domain Model Spans
type ProcessLogs func(logs []*model.LogRecord, tenant string)

// FilterLog decides whether to allow or disallow a log
type FilterLog func(log *model.LogRecord) bool

// ChainedProcessLog chains logProcessors as a single ProcessSpan call
func ChainedProcessLog(logProcessors ...ProcessLog) ProcessLog {
	return func(log *model.LogRecord, tenant string) {
		for _, processor := range logProcessors {
			processor(log, tenant)
		}
	}
}
