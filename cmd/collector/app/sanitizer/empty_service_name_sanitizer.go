// Copyright (c) 2022 The Jaeger Authors.
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

package sanitizer

import (
	"logger/model"
)

const (
	serviceNameReplacement = "empty-service-name"
	nullProcessServiceName = "null-process-and-service-name"
)

// NewEmptyServiceNameSanitizer returns a function that replaces empty service name
// with a string "empty-service-name".
// If the whole log.Process is null, it creates one with "null-process-and-service-name".
func NewEmptyServiceNameSanitizer() SanitizeSpan {
	return sanitizeEmptyServiceName
}

// Sanitize sanitizes the service names in the log annotations.
func sanitizeEmptyServiceName(log *model.LogRecord) *model.LogRecord {
	if log.Process == nil {
		log.Process = &model.Process{ServiceName: nullProcessServiceName}
	} else if log.Process.ServiceName == "" {
		log.Process.ServiceName = serviceNameReplacement
	}
	return log
}
