package cassandra

import (
	"flag"
	"logger/pkg/cassandra"
	"logger/pkg/cassandra/config"
	"logger/pkg/metrics"
	"logger/plugin"
	"logger/storage"

	"io"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	ls "logger/storage/logstore"
)

const (
	primaryStorageConfig = "cassandra"
	archiveStorageConfig = "cassandra-archive"
)


var ( // interface comformance checks
	_ storage.FactoryBase              = (*Factory)(nil)
	// _ storage.Purger               = (*Factory)(nil)
	// _ storage.ArchiveFactory       = (*Factory)(nil)
	// _ storage.SamplingStoreFactory = (*Factory)(nil)
	_ io.Closer                    = (*Factory)(nil)
	_ plugin.Configurable          = (*Factory)(nil)
)

type Factory struct {
	Options *Options

	primaryMetricsFactory metrics.Factory
	archiveMetricsFactory metrics.Factory
	logger                *zap.Logger
	// tracer                trace.TracerProvider

	primaryConfig  config.SessionBuilder
	primarySession cassandra.Session
	archiveConfig  config.SessionBuilder
	archiveSession cassandra.Session
}


func NewFactory() *Factory {
	return &Factory{
		// tracer:  otel.GetTracerProvider(),
		Options: NewOptions(primaryStorageConfig, archiveStorageConfig),
	}
}



func (f *Factory) Initialiace(metricFactory metrics.Factory, logger *zap.Logger) error {
	return nil
}

func (f *Factory) Close() error {
	return nil
}

func (f *Factory) CreateLogReader() (ls.Reader, error) {
	return nil, nil
}
func (f *Factory) CreateLogWriter() (ls.Writer, error) {
	return nil, nil
}

func (f *Factory) AddFlags(flagSet *flag.FlagSet) {
}

func (f *Factory) InitFromViper(v *viper.Viper, logger *zap.Logger) {
}
