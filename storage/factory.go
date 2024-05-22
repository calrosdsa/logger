package storage

import (
	"logger/pkg/metrics"
	"go.uber.org/zap"
	"logger/storage/logstore"
)

type FactoryBase interface {
	Initialize(metricFactory metrics.Factory, logger *zap.Logger) error
	Close() error	
	CreateLogReader()(logstore.Reader,error)
	CreateLogWriter()(logstore.Writer,error)
}