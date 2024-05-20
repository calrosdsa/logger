package storage

import (
	"logger/pkg/metrics"
	"go.uber.org/zap"

)

type FactoryBase interface {
	Initialiace(metricFactory metrics.Factory, logger *zap.Logger) error
	Close() error	
}