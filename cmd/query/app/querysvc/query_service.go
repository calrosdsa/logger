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

func (s *QueryService)GetLogs(ctx context.Context,query logstore.LogQueryParameters)([]*model.LogRecord,error) {
	return s.logReader.GetLogs(ctx,query)
}

func (s *QueryService)GetServices(ctx context.Context)([]string,error){
	return s.logReader.GetServices(ctx)
}

func (s *QueryService) GetOperations(ctx context.Context,query logstore.OperationQueryParameters)([]logstore.Operation,error){
	return s.logReader.GetOperations(ctx,query)
}
