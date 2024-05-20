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
