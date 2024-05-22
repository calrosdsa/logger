package cassandra

import (
	"flag"
	"fmt"
	"logger/pkg/cassandra"
	"logger/pkg/cassandra/config"
	"logger/pkg/metrics"
	"logger/plugin"
	"logger/storage"

	"io"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	ls "logger/storage/logstore"

	cLogStore "logger/plugin/storage/cassandra/logstore"
)

const (
	primaryStorageConfig = "cassandra"
	archiveStorageConfig = "cassandra-archive"
)

var ( // interface comformance checks
	_ storage.FactoryBase = (*Factory)(nil)
	// _ storage.Purger               = (*Factory)(nil)
	// _ storage.ArchiveFactory       = (*Factory)(nil)
	// _ storage.SamplingStoreFactory = (*Factory)(nil)
	_ io.Closer           = (*Factory)(nil)
	_ plugin.Configurable = (*Factory)(nil)
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
	cs := config.Configuration{
		Servers: []string{"localhost"},
		Port: 9042,
		Keyspace: "kspace",
		Authenticator: config.Authenticator{
			Basic: config.BasicAuthenticator{
				Username: "admin",
				Password: "admin",
			},
		},
		
	}
	return &Factory{
		// tracer:  otel.GetTracerProvider(),
		Options: NewOptions(primaryStorageConfig, archiveStorageConfig),
		primaryConfig: &cs,
	}
}

func (f *Factory) Initialize(metricsFactory metrics.Factory, logger *zap.Logger) error {
	fmt.Println("INITIALIZE DB CASSANDRA")
	f.primaryMetricsFactory = metricsFactory.Namespace(metrics.NSOptions{Name: "cassandra", Tags: nil})
	f.archiveMetricsFactory = metricsFactory.Namespace(metrics.NSOptions{Name: "cassandra-archive", Tags: nil})
	f.logger = logger
	fmt.Println("Primary config",f.primaryConfig,"LOGGER--",logger)
	primarySession, err := f.primaryConfig.NewSession(logger)
	if err != nil {
		return err
	}
	f.primarySession = primarySession
	if f.archiveConfig != nil {
		if archiveSession, err := f.archiveConfig.NewSession(logger); err == nil {
			f.archiveSession = archiveSession
		} else {
			return err
		}
	} else {
		logger.Info("Cassandra archive storage configuration is empty, skipping")
	}
	return nil
}

func (f *Factory) Close() error {
	return nil
}

func (f *Factory) CreateLogReader() (ls.Reader, error) {
	return nil, nil
}
func (f *Factory) CreateLogWriter() (ls.Writer, error) {
	fmt.Println("CRATEING LOG WRITER CASSANDRA")
	options, err := writerOptions(f.Options)
	if err != nil {
		return nil, err
	}
	return cLogStore.NewLogWriter(
		f.primarySession, f.Options.SpanStoreWriteCacheTTL, f.primaryMetricsFactory, f.logger, options...), nil
}

func (f *Factory) AddFlags(flagSet *flag.FlagSet) {
}

func (f *Factory) InitFromViper(v *viper.Viper, logger *zap.Logger) {
}

func writerOptions(opts *Options) ([]cLogStore.Option, error) {
	// var tagFilters []dbmodel.TagFilter

	// // drop all tag filters
	// if !opts.Index.Tags || !opts.Index.ProcessTags || !opts.Index.Logs {
	// 	tagFilters = append(tagFilters, dbmodel.NewTagFilterDropAll(!opts.Index.Tags, !opts.Index.ProcessTags, !opts.Index.Logs))
	// }

	// black/white list tag filters
	// tagIndexBlacklist := opts.TagIndexBlacklist()
	// tagIndexWhitelist := opts.TagIndexWhitelist()
	// if len(tagIndexBlacklist) > 0 && len(tagIndexWhitelist) > 0 {
	// 	return nil, errors.New("only one of TagIndexBlacklist and TagIndexWhitelist can be specified")
	// }
	// if len(tagIndexBlacklist) > 0 {
	// 	tagFilters = append(tagFilters, dbmodel.NewBlacklistFilter(tagIndexBlacklist))
	// } else if len(tagIndexWhitelist) > 0 {
	// 	tagFilters = append(tagFilters, dbmodel.NewWhitelistFilter(tagIndexWhitelist))
	// }

	// if len(tagFilters) == 0 {
	// 	return nil, nil
	// } else if len(tagFilters) == 1 {
	// 	return []cLogStore.Option{cLogStore.TagFilter(tagFilters[0])}, nil
	// }

	return []cLogStore.Option{}, nil
}
