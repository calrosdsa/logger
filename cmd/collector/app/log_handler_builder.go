package app

import (
	"logger/storage/logstore"
	"os"

	"logger/cmd/collector/app/flags"
	"logger/cmd/collector/app/handler"
	"logger/cmd/collector/app/processor"

	"logger/pkg/metrics"
	"logger/pkg/tenancy"

	"go.uber.org/zap"
)

type LogHandlerBuilder struct {
	LogWriter      logstore.Writer
	CollectorOpts  *flags.CollectorOptions
	Logger         *zap.Logger
	MetricsFactory metrics.Factory
	TenancyMgr     *tenancy.Manager
}

type LogHandlers struct {
	BatchesHandler handler.BatchesHandler
}

func (b *LogHandlerBuilder) BuildLogProcessor(additional ...ProcessLog) processor.LogProcessor {
	hostname,_ := os.Hostname()
	svcMetrics := b.metricsFactory()
	hostMetrics := svcMetrics.Namespace(metrics.NSOptions{Tags: map[string]string{"host": hostname}})
	return NewSpan
}

func (b *LogHandlerBuilder) logger() *zap.Logger {
	if b.Logger == nil {
		return zap.NewNop()
	}
	return b.Logger
}

func (b *LogHandlerBuilder) metricsFactory() metrics.Factory {
	if b.MetricsFactory == nil {
		return metrics.NullFactory
	}
	return b.MetricsFactory
}
