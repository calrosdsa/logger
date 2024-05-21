package main

// import (
// 	// "fmt"
// 	// "log"

// 	// "github.com/valyala/fasthttp"
// 	"log"
// 	_log "logger/core/modules/logs/delivery"

// 	"github.com/savsgio/atreugo/v11"
// 	// "github.com/fasthttp/router"
// )

// type options struct {
// 	id   int
// 	uuid string
// }

// type Option struct {
// 	f func(c *options)
// }

// var Options options

// func (o options) Id(id int) Option {
// 	return Option{f: func(c *options) {
// 		c.id = id
// 	}}
// }

// func (o options) Uuid(uuid string) Option {
// 	return Option{f: func(c *options) {
// 		c.uuid = uuid
// 	}}
// }

// func (o options) apply(opts ...Option) options {
// 	ret := options{}
// 	for _, opt := range opts {
// 		opt.f(&ret)
// 	}
// 	if ret.id == 0 {
// 		ret.id = 1
// 	}
// 	if ret.uuid == "" {
// 		ret.uuid = "DEFAULT VALUE"
// 	}
// 	return ret
// }

// func main() {
// 	config := atreugo.Config{
// 		Addr:      "0.0.0.0:8000",
// 		TLSEnable: false,
// 	}
// 	server := atreugo.New(config)

// 	options := Options.apply(Options.Id(1), Options.Uuid("281818`"))
// 	log.Println(options)

// 	// server.GET("/", func(ctx *atreugo.RequestCtx) error {
// 	// 	res := struct {
// 	// 		Name string `json:"name"`
// 	// 	}{
// 	// 		Name: "Jorge",
// 	// 	}
// 	// 	return ctx.JSONResponse(res)
// 	// })

// 	// server.GET("/echo/{path:*}", func(ctx *atreugo.RequestCtx) error {
// 	// 	return ctx.TextResponse("Echo message: " + ctx.UserValue("path").(string))
// 	// })

// 	// v1 := server.NewGroupPath("/v1")
// 	_log.New(server)

// 	// v1.GET("/", func(ctx *atreugo.RequestCtx) error {
// 	// 	return ctx.TextResponse("Hello V1 Group")
// 	// })

// 	if err := server.ListenAndServe(); err != nil {
// 		panic(err)
// 	}

// }
