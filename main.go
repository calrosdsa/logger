package main

import (
	// "fmt"
	// "log"

	// "github.com/valyala/fasthttp"
	"github.com/savsgio/atreugo/v11"
	_log "logger/core/modules/logs/delivery"
	// "github.com/fasthttp/router"
)

func main(){
	config := atreugo.Config{
		Addr: "0.0.0.0:8000",
		TLSEnable: false,
	}
	server := atreugo.New(config)

	// server.GET("/", func(ctx *atreugo.RequestCtx) error {
	// 	res := struct {
	// 		Name string `json:"name"`
	// 	}{
	// 		Name: "Jorge",
	// 	}
	// 	return ctx.JSONResponse(res)
	// })

	// server.GET("/echo/{path:*}", func(ctx *atreugo.RequestCtx) error {
	// 	return ctx.TextResponse("Echo message: " + ctx.UserValue("path").(string))
	// })

	// v1 := server.NewGroupPath("/v1")
	_log.New(server)

	// v1.GET("/", func(ctx *atreugo.RequestCtx) error {
	// 	return ctx.TextResponse("Hello V1 Group")
	// })

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
	
}