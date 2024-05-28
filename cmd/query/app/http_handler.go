package app

import (
	"encoding/json"
	"fmt"
	"logger/cmd/query/app/querysvc"
	"logger/storage/logstore"
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
	server.UseBefore(func(rc *atreugo.RequestCtx) error {
		rc.Response.Header.Set("Access-Control-Allow-Origin", "*")
		return rc.Next()
	})

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
	router.POST("/v1/logs/",aH.GetLogs)
	router.GET("/v1/services/",aH.GetServices)
	router.POST("/v1/operations/",aH.GetOperations)
}


func (aH *APIHandler) GetOperations(c *atreugo.RequestCtx) error {
	ctx := c.AttachedContext()
	data := c.PostBody()
	var query logstore.OperationQueryParameters
	err := json.Unmarshal(data,&query)
	if err != nil {
		aH.logger.Error("GetOperations",zap.Error(err))
		return c.JSONResponse(structuredError{
			Msg: err.Error(),
			Code: http.StatusUnprocessableEntity,
		})
	}
	operations,err := aH.queryService.GetOperations(ctx,query)
	if err != nil {
		aH.logger.Error("GetOperations",zap.Error(err))
		return c.JSONResponse(structuredError{
			Msg: err.Error(),
			Code: http.StatusBadRequest,
		})
	}
	return c.JSONResponse(operations,http.StatusOK)
}

func (aH *APIHandler) GetServices(c *atreugo.RequestCtx) error {
	ctx := c.AttachedContext()
	// c. .Header.Set("Access-Control-Allow-Origin", "*")
	services,err := aH.queryService.GetServices(ctx)
	if err != nil {
		aH.logger.Error("GetServices",zap.Error(err))
		return c.JSONResponse(structuredError{
			Msg: err.Error(),
			Code: http.StatusBadRequest,
		})
	}
	return c.JSONResponse(services,http.StatusOK)
}

func (aH *APIHandler) GetLogs(c *atreugo.RequestCtx)error {
	ctx := c.AttachedContext()
	data := c.PostBody()
	var query logstore.LogQueryParameters
	err := json.Unmarshal(data,&query)
	if err != nil {
		aH.logger.Error("GetOperations",zap.Error(err))
		return c.JSONResponse(structuredError{
			Msg: err.Error(),
			Code: http.StatusUnprocessableEntity,
		})
	}
	logs,err := aH.queryService.GetLogs(ctx,query)
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