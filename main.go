// package main
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


package main

import (
    "fmt"
    "time"
)

func main() {
    // The given timestamp in nanoseconds
    timestamp := int64(1716837192181714500)

    // Convert the nanoseconds timestamp to a time.Time object
    t := time.Unix(0, timestamp)

    // Print the time in a readable format
    fmt.Println("The time is:", t)

	// r := EpochMicrosecondsAsTime(uint64(timestamp))
	// fmt.Println(r.Format(time.RFC3339))

	// r2 := TimeAsEpochMicroseconds(r)
	fmt.Println(t.UnixNano())
}

func EpochMicrosecondsAsTime(ts uint64) time.Time {
	seconds := ts / 1000000
	nanos := 1000 * (ts % 1000000)
	return time.Unix(int64(seconds), int64(nanos)).UTC()
}

func TimeAsEpochMicroseconds(t time.Time) uint64 {
	return uint64(t.UnixNano() / 1000)
}
