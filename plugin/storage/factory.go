package storage

import (
	"flag"
	"fmt"
	"io"
	"logger/pkg/metrics"
	"logger/plugin"
	"logger/storage"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"logger/plugin/storage/cassandra"
	"logger/storage/logstore"
	ls "logger/storage/logstore"
)

const (
	cassandraStorageType = "cassandra"
	logStorageType       = "log-storage-type"

	// defaultDownsamplingRatio is the default downsampling ratio.
	// defaultDownsamplingHashSalt is the default downsampling hashsalt.
)

// AllStorageTypes defines all available storage backends
var AllStorageTypes = []string{
	cassandraStorageType,
}

var ( // interface comformance checks
	_ storage.FactoryBase = (*Factory)(nil)
	// _ storage.ArchiveFactory = (*Factory)(nil)
	_ io.Closer           = (*Factory)(nil)
	_ plugin.Configurable = (*Factory)(nil)
)

// Factory implements storage.Factory interface as a meta-factory for storage components.
type Factory struct {
	FactoryConfig
	metricsFactory         metrics.Factory
	factories              map[string]storage.FactoryBase
	downsamplingFlagsAdded bool
}

// NewFactory creates the meta-factory.
func NewFactory(config FactoryConfig) (*Factory, error) {
	f := &Factory{FactoryConfig: config}
	uniqueTypes := map[string]struct{}{
		f.LogReaderType:          {},
		f.DependenciesStorageType: {},
	}
	for _, storageType := range f.LogWriterTypes {
		uniqueTypes[storageType] = struct{}{}
	}
	// skip SamplingStorageType if it is empty. See CreateSamplingStoreFactory for details
	// if f.SamplingStorageType != "" {
	// 	uniqueTypes[f.SamplingStorageType] = struct{}{}
	// }
	f.factories = make(map[string]storage.FactoryBase)
	for t := range uniqueTypes {
		ff, err := f.getFactoryOfType(t)
		if err != nil {
			return nil, err
		}
		f.factories[t] = ff
	}
	return f, nil
}


func (f *Factory) getFactoryOfType(factoryType string) (storage.FactoryBase, error) {
	switch factoryType {
	case cassandraStorageType:
		return cassandra.NewFactory(), nil
	default:
		return nil, fmt.Errorf("unknown storage type %s. Valid types are %v", factoryType, AllStorageTypes)
	}
}

func (f *Factory) Initialize(metricsFactory metrics.Factory, logger *zap.Logger) error {
	f.metricsFactory = metricsFactory
	for _, factory := range f.factories {
		if err := factory.Initialize(metricsFactory, logger); err != nil {
			return err
		}
	}
	// f.publishOpts()

	return nil
}

// func (f *Factory) publishOpts() {
// 	internalFactory := f.metricsFactory.Namespace(metrics.NSOptions{Name: "internal"})
// 	internalFactory.Gauge(metrics.Options{Name: downsamplingRatio}).
// 		Update(int64(f.FactoryConfig.DownsamplingRatio))
// 	internalFactory.Gauge(metrics.Options{Name: spanStorageType + "-" + f.SpanReaderType}).
// 		Update(1)
// }


func (f *Factory) Close() error {
	return nil
}

func (f *Factory) CreateLogReader() (ls.Reader, error) {
	return nil, nil
}
func (f *Factory) CreateLogWriter() (ls.Writer, error) {
	fmt.Println("CREATING LOG WRITER")
	var writers []logstore.Writer
	for _, storageType := range f.LogWriterTypes {
		factory ,ok := f.factories[storageType]
		if !ok {
			return nil, fmt.Errorf("no %s backend registered for span store", storageType)
		}
		writer, err := factory.CreateLogWriter()
		if err != nil {
			return nil, err
		}
		writers = append(writers, writer)
	}
	var spanWriter logstore.Writer
	if len(f.LogWriterTypes) == 1 {
		spanWriter = writers[0]
	} else {
		fmt.Println("IS GREATER THAT 1")
		// spanWriter = spanstore.NewCompositeWriter(writers...)
	}
	return spanWriter,nil
	// Turn off DownsamplingWriter entirely if ratio == defaultDownsamplingRatio.
	// if f.DownsamplingRatio == defaultDownsamplingRatio {
	// 	return spanWriter, nil
	// }
	// return spanstore.NewDownsamplingWriter(spanWriter, spanstore.DownsamplingOptions{
	// 	Ratio:          f.DownsamplingRatio,
	// 	HashSalt:       f.DownsamplingHashSalt,
	// 	MetricsFactory: f.metricsFactory.Namespace(metrics.NSOptions{Name: "downsampling_writer"}),
	// }), nil
}

func (f *Factory) AddFlags(flagSet *flag.FlagSet) {
}

func (f *Factory) InitFromViper(v *viper.Viper, logger *zap.Logger) {
}
