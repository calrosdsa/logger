package app

import (
	"fmt"
	"logger/cmd/query/app/querysvc"
	"net/http"

	"github.com/savsgio/atreugo/v11"
	"go.uber.org/zap"
)


type structuredError struct {
	Code    int        `json:"code,omitempty"`
	Msg     string     `json:"msg"`
}


type HttpHandler interface {
	RegisterRoutes(router *atreugo.Atreugo)
}



func NewRouter() *atreugo.Atreugo {
	config := atreugo.Config{
		Addr:      "0.0.0.0:8001",
		TLSEnable: false,
	}
	server := atreugo.New(config)
	return server
}


type APIHandler struct {
	queryService *querysvc.QueryService
	logger *zap.Logger
}

func NewAPIHandler(qsvc *querysvc.QueryService,options ...HandlerOption) *APIHandler{
	aH := &APIHandler{
		queryService: qsvc,
	}
	for _,option := range options {
		option(aH)
	}
	if aH.logger == nil {
		aH.logger = zap.NewNop()
	}

	return aH
}

func (aH *APIHandler) RegisterRoutes(router *atreugo.Atreugo){
	router.GET("/v1/logs/",aH.GetLogs)
}

func (aH *APIHandler) GetLogs(c *atreugo.RequestCtx)error {
	ctx := c.AttachedContext()
	logs,err := aH.queryService.GetLogs(ctx)
	if err != nil {
		aH.logger.Error("GerLogs",zap.Error(err))
		return c.JSONResponse(structuredError{
			Msg: err.Error(),
			Code: http.StatusBadRequest,
		})
	}
	fmt.Println("LOGS",logs)
	return c.JSONResponse(logs,http.StatusOK)
}