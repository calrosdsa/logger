package app

import (
	"logger/cmd/query/app/querysvc"
	"logger/pkg/healthcheck"

	"github.com/savsgio/atreugo/v11"
	"go.uber.org/zap"
)

type Server struct {
	logger *zap.Logger
	querySvc *querysvc.QueryService
	healthCheck  *healthcheck.HealthCheck
	queryOptions QueryOptions
	server *atreugo.Atreugo
}

func NewServer(
	logger *zap.Logger,
	healtcheck *healthcheck.HealthCheck,
	querySvc *querysvc.QueryService,
	options *QueryOptions,
)(*Server,error){
	server :=  createHttpServer(logger,querySvc)

	return &Server{
		logger: logger,
		querySvc: querySvc,
		healthCheck: healtcheck,
		queryOptions: *options,
		server: server,
	},nil
}

func createHttpServer(
	logger *zap.Logger,
	querySvc *querysvc.QueryService,
)(*atreugo.Atreugo){
	apiHandler := NewAPIHandler(querySvc,HandlerOptions.Logger(logger))
	r := NewRouter()
	apiHandler.RegisterRoutes(r)
	return r
}

func (aH *Server) Start() error{
	if err := aH.server.ListenAndServe(); err != nil {
		aH.logger.Error("Fail to start server",zap.Error(err))
	    return err
	}
	return nil
}

func (aH *Server) Close() error {
	return aH.server.Shutdown()
}

