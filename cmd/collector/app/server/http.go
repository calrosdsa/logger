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

package server

import (
	"logger/cmd/collector/app/handler"

	"github.com/savsgio/atreugo/v11"
	"go.uber.org/zap"
	// "go.uber.org/zap/zapcore"
)

type HttpServerParams struct {
	Handler handler.BatchesHandler
	Logger *zap.Logger
	HostPort string
}

func StartHTTPServer(params *HttpServerParams)(error){
	params.Logger.Info("Starting collector HTTP server", zap.String("http host-port", params.HostPort))
	// errorLog, _ := zap.NewStdLogAt(params.Logger, zapcore.ErrorLevel)

	config := atreugo.Config{
		Addr:      "0.0.0.0:8000",
		TLSEnable: false,
	}
	server := atreugo.New(config)
	serveHttp(server,params)
	err := server.ListenAndServe();
	if err != nil {
     	params.Logger.Info("Fail to start collector HTTP server", zap.Error(err))
	}
	return err
}

func serveHttp(server *atreugo.Atreugo,params *HttpServerParams){
	apiHandler := handler.NewAPIHandler(params.Handler)
	apiHandler.RegisterRoutes(server)
}