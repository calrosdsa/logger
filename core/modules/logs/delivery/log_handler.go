package delivery

import (
	"log"
	"net/http"

	"github.com/savsgio/atreugo/v11"

	pb "logger/model/proto"
	"fmt"
	"github.com/golang/protobuf/proto"

	// "encoding/json"
	// "fmt"
)

type logHandler struct {}

func New(router *atreugo.Atreugo){
	handler := logHandler{}
	router.POST("/v1/logs",handler.Logs)
}

func (h logHandler) Logs(c *atreugo.RequestCtx) (err error) {
	data := c.PostBody()
	// data := []byte{
    //     10, 168, 1, 10, 56, 10, 24, 10, 16, 108, 105, 98, 114, 97, 114, 121, 46, 108, 97, 110, 103, 117, 97, 103, 101, 18, 4, 10, 2, 103, 111, 10, 28, 10, 12, 115, 101, 114, 118, 105, 99, 101, 46, 110, 97, 109, 101, 18, 12, 10, 10, 109, 121, 45, 115, 101, 114, 118, 105, 99, 101, 18, 108, 10, 61, 10, 46, 103, 111, 46, 111, 112, 101, 110, 116, 101, 108, 101, 109, 101, 116, 114, 121, 46, 105, 111, 47, 99, 111, 110, 116, 114, 105, 98, 47, 98, 114, 105, 100, 103, 101, 115, 47, 111, 116, 101, 108, 108, 111, 103, 114, 117, 115, 18, 11, 48, 46, 48, 46, 49, 45, 97, 108, 112, 104, 97, 18, 43, 9, 252, 49, 189, 137, 127, 165, 208, 23, 16, 4, 42, 21, 10, 19, 73, 78, 73, 84, 32, 76, 79, 71, 71, 69, 82, 32, 83, 69, 82, 86, 73, 67, 69, 89, 252, 49, 189, 137, 127, 165, 208, 23,
    // }
    fmt.Printf("FETCHING DATA ------------------")
	var example pb.ExportLogsServiceRequest
    err = proto.Unmarshal(data, &example)
    if err != nil {
        log.Fatalf("Failed to parse data: %v", err)
    }
    // Print the decoded data
    // fmt.Printf("Library Language: %s\n", example.GetLibrary().GetLanguage())
    // fmt.Printf("Service Name: %s\n", example.GetService().GetName())
    // fmt.Printf("URL: %s\n", example.GetUrl())
    // fmt.Printf("Init Logger Service: %v\n", example.GetInitLoggerService())
    // fmt.Printf("Logger Service URL: %s\n", example.GetLoggerServiceUrl())
    fmt.Printf("Scheme URL: %s\n", example.GetResourceLogs())
	return c.JSONResponse("SUCESS",http.StatusOK)
}
