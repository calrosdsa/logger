package app

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"

	"logger/cmd/collector/app/flags"
	"logger/cmd/collector/app/processor"
	"logger/cmd/collector/app/server"

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
	// grpcServer                 *grpc.Server
	otlpReceiver               receiver.Traces
	// zipkinReceiver             receiver.Traces
	// tlsGRPCCertWatcherCloser   io.Closer
	tlsHTTPCertWatcherCloser   io.Closer
	// tlsZipkinCertWatcherCloser io.Closer
}

type CollectorParams struct {
	ServiceName    string
	Logger         *zap.Logger
	MetricsFactory metrics.Factory
	LogWriter     logstore.Writer
	// StrategyStore  strategystore.StrategyStore
	// Aggregator     strategystore.Aggregator
	HealthCheck    *healthcheck.HealthCheck
	TenancyMgr     *tenancy.Manager
}

func New(params *CollectorParams) *Collector {
	return &Collector{
		serviceName:    params.ServiceName,
		logger:         params.Logger,
		metricsFactory: params.MetricsFactory,
		logWriter:     params.LogWriter,
		// strategyStore:  params.StrategyStore,
		// aggregator:     params.Aggregator,
		hCheck:         params.HealthCheck,
		tenancyMgr:     params.TenancyMgr,
	}
}



// Start the component and underlying dependencies
func (c *Collector) Start(options *flags.CollectorOptions) error {
	handlerBuilder := &LogHandlerBuilder{
		LogWriter:     c.logWriter,
		CollectorOpts:  options,
		Logger:         c.logger,
		MetricsFactory: c.metricsFactory,
		TenancyMgr:     c.tenancyMgr,
	}

	var additionalProcessors []ProcessLog
	// if c.aggregator != nil {
	// 	additionalProcessors = append(additionalProcessors, handleRootSpan(c.aggregator, c.logger))
	// }

	c.logProcessor = handlerBuilder.BuildLogProcessor(additionalProcessors...)
	c.logHandlers = handlerBuilder.BuildHandlers(c.logProcessor)

	err := server.StartHTTPServer(&server.HttpServerParams{
		Handler: c.logHandlers.BatchesHandler,
		Logger: c.logger,
		HostPort: options.HTTP.HostPort,
	})
	if err != nil {
		return fmt.Errorf("could not start HTTP server: %w", err)
	}

	// grpcServer, err := server.StartGRPCServer(&server.GRPCServerParams{
	// 	HostPort:                options.GRPC.HostPort,
	// 	Handler:                 c.logHandlers.GRPCHandler,
	// 	TLSConfig:               options.GRPC.TLS,
	// 	SamplingStore:           c.strategyStore,
	// 	Logger:                  c.logger,
	// 	MaxReceiveMessageLength: options.GRPC.MaxReceiveMessageLength,
	// 	MaxConnectionAge:        options.GRPC.MaxConnectionAge,
	// 	MaxConnectionAgeGrace:   options.GRPC.MaxConnectionAgeGrace,
	// })
	// if err != nil {
	// 	return fmt.Errorf("could not start gRPC server: %w", err)
	// }
	// c.grpcServer = grpcServer

	// httpServer, err := server.StartHTTPServer(&server.HTTPServerParams{
	// 	HostPort:       options.HTTP.HostPort,
	// 	Handler:        c.logHandlers.JaegerBatchesHandler,
	// 	TLSConfig:      options.HTTP.TLS,
	// 	HealthCheck:    c.hCheck,
	// 	MetricsFactory: c.metricsFactory,
	// 	SamplingStore:  c.strategyStore,
	// 	Logger:         c.logger,
	// })
	// if err != nil {
	// 	return fmt.Errorf("could not start HTTP server: %w", err)
	// }
	// c.hServer = httpServer

	// c.tlsGRPCCertWatcherCloser = &options.GRPC.TLS
	// c.tlsHTTPCertWatcherCloser = &options.HTTP.TLS
	// c.tlsZipkinCertWatcherCloser = &options.Zipkin.TLS

	
	// if options.OTLP.Enabled {
	// 	otlpReceiver, err := handler.StartOTLPReceiver(options, c.logger, c.logProcessor, c.tenancyMgr)
	// 	if err != nil {
	// 		return fmt.Errorf("could not start OTLP receiver: %w", err)
	// 	}
	// 	c.otlpReceiver = otlpReceiver
	// }

	// c.publishOpts(options)

	return nil
}

func (c *Collector) publishOpts(cOpts *flags.CollectorOptions) {
	internalFactory := c.metricsFactory.Namespace(metrics.NSOptions{Name: "internal"})
	internalFactory.Gauge(metrics.Options{Name: metricNumWorkers}).Update(int64(cOpts.NumWorkers))
	internalFactory.Gauge(metrics.Options{Name: metricQueueSize}).Update(int64(cOpts.QueueSize))
}

// Close the component and all its underlying dependencies
func (c *Collector) Close() error {
	// Stop gRPC server

	// Stop HTTP server
	if c.hServer != nil {
		timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := c.hServer.Shutdown(timeout); err != nil {
			c.logger.Fatal("failed to stop the main HTTP server", zap.Error(err))
		}
		defer cancel()
	}

	

	// Stop OpenTelemetry OTLP receiver
	if c.otlpReceiver != nil {
		timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := c.otlpReceiver.Shutdown(timeout); err != nil {
			c.logger.Fatal("failed to stop the OTLP receiver", zap.Error(err))
		}
		defer cancel()
	}

	if err := c.logProcessor.Close(); err != nil {
		c.logger.Error("failed to close span processor.", zap.Error(err))
	}

	// // aggregator does not exist for all strategy stores. only Close() if exists.
	// if c.aggregator != nil {
	// 	if err := c.aggregator.Close(); err != nil {
	// 		c.logger.Error("failed to close aggregator.", zap.Error(err))
	// 	}
	// }

	// // watchers actually never return errors from Close
	// if c.tlsGRPCCertWatcherCloser != nil {
	// 	_ = c.tlsGRPCCertWatcherCloser.Close()
	// }
	if c.tlsHTTPCertWatcherCloser != nil {
		_ = c.tlsHTTPCertWatcherCloser.Close()
	}
	// if c.tlsZipkinCertWatcherCloser != nil {
	// 	_ = c.tlsZipkinCertWatcherCloser.Close()
	// }

	return nil
}

// LogHandlers returns span handlers used by the Collector.
func (c *Collector) LogHandlers() *LogHandlers {
	return c.logHandlers
}
