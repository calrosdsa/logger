package handler

import (
	"log"
	"net/http"

	"github.com/savsgio/atreugo/v11"

	"fmt"
	"logger/cmd/collector/app/processor"
	pb "logger/model/proto/v1"

	"github.com/golang/protobuf/proto"
	// "encoding/json"
	// "fmt"
)

type logHandler struct {
	jaegerBatchesHandler BatchesHandler
}

func New(router *atreugo.Atreugo) {
	handler := logHandler{}
	router.POST("/v1/logs", handler.Logs)
}

func (h logHandler) Logs(c *atreugo.RequestCtx) (err error) {
	fmt.Printf("FETCHING DATA ------------------")
	var data pb.ExportLogsServiceRequest
	err = proto.Unmarshal(c.PostBody(), &data)
	if err != nil {
		log.Fatalf("Failed to parse data: %v", err)
	}
	batches := data.GetResourceLogs()
	opts := SubmitBatchOptions{InboundTransport: processor.HTTPTransport}
	if _, err = h.jaegerBatchesHandler.SubmitBatches(batches, opts); err != nil {
		return c.JSONResponse("ERROR", http.StatusBadRequest)
	}

	fmt.Printf("Scheme URL: %s\n", data.GetResourceLogs())
	return c.JSONResponse("SUCESS", http.StatusOK)
}
