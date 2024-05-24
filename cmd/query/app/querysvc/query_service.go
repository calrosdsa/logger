package querysvc

import (
	"context"
	"logger/model"
	"logger/storage/logstore"
)

type QueryService struct {
	logReader logstore.Reader
}

func NewQueryService(logReader logstore.Reader) *QueryService {
	qsvc := &QueryService{
		logReader: logReader,
	}
	return qsvc
}

func (s *QueryService)GetLogs(ctx context.Context)([]*model.LogRecord,error) {
	return s.logReader.GetLogs(ctx)
}

