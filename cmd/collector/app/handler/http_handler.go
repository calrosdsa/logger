package handler

import (
	"log"
	"net/http"

	"github.com/savsgio/atreugo/v11"

	// "fmt"
	"logger/cmd/collector/app/processor"
	pb "logger/model/proto/v1"

	"github.com/golang/protobuf/proto"
	// "encoding/json"
	// "fmt"
)

type APIHandler struct {
	BatchesHandler BatchesHandler
}

func NewAPIHandler(
	BatchesHandler BatchesHandler,
) *APIHandler {
	return &APIHandler{
		BatchesHandler: BatchesHandler,
	}
}


func (h *APIHandler)RegisterRoutes(router *atreugo.Atreugo) {
	router.POST("/v1/logs", h.Logs)
}




func (h *APIHandler) Logs(c *atreugo.RequestCtx) (err error) {
	var data pb.ExportLogsServiceRequest
	err = proto.Unmarshal(c.PostBody(), &data)
	if err != nil {
		log.Fatalf("Failed to parse data: %v", err)
	}
	// log.Println(data.GetResourceLogs())
	
	batches := data.GetResourceLogs()
	// for _,b := range batches {
	// 	for _,v := range b.ScopeLogs {
	// 		for _,l := range v.LogRecords {
	// 			l.Body.GetStringValue()
	// 			fmt.Println("Number",l.SeverityNumber.Number(),
	// 			"String",l.SeverityNumber.String())
	// 		}
	// 	}
	// }
	opts := SubmitBatchOptions{InboundTransport: processor.HTTPTransport}
	if _, err = h.BatchesHandler.SubmitBatches(batches, opts); err != nil {
		return c.JSONResponse("ERROR", http.StatusBadRequest)
	}

	// fmt.Printf("Scheme URL: %s\n", data.GetResourceLogs())
	return c.JSONResponse("SUCESS", http.StatusOK)
}
