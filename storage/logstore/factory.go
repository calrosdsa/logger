package logstore

import (
	"logger/storage"
)

type Factory interface {
	storage.FactoryBase
	CreateLogReader()(Reader,error)
	CreateLogWriter()(Writer,error)
}