package app

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"logger/cmd/collector/app/flags"
	"logger/cmd/collector/app/handler"
	"logger/cmd/collector/app/processor"
	// "logger/cmd/collector/app/sampling/strategystore"
	// "logger/cmd/collector/app/server"
	"logger/pkg/healthcheck"
	"logger/pkg/metrics"
	"logger/pkg/tenancy"
	"logger/storage/logstore"
)

const (
	metricNumWorkers = "collector.num-workers"
	metricQueueSize  = "collector.queue-size"
)


type Collector struct {
	serviceName    string
	logger         *zap.Logger
	metricsFactory metrics.Factory
	logWriter     logstore.Writer
	// strategyStore  strategystore.StrategyStore
	// aggregator     strategystore.Aggregator
	hCheck         *healthcheck.HealthCheck
	logProcessor  processor.LogProcessor
	logHandlers   *LogHandlers
	tenancyMgr     *tenancy.Manager

	// state, read only
	hServer                    *http.Server
	grpcServer                 *grpc.Server
	otlpReceiver               receiver.Traces
	zipkinReceiver             receiver.Traces
	tlsGRPCCertWatcherCloser   io.Closer
	tlsHTTPCertWatcherCloser   io.Closer
	tlsZipkinCertWatcherCloser io.Closer
}