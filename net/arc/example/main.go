package main

import (
	"landzero.net/x/net/arc"
	"net/http"
)

func main() {
	a := arc.New()
	r1 := a.Route().Get("/")
	r1.Use(func(ctx *arc.Context) {
		ctx.Res.Write([]byte("hello"))
		ctx.Next()
	}, func(ctx *arc.Context) {
		ctx.Res.Write([]byte(", static"))
	})
	r2 := a.Route().Get("/hello/:name")
	r2.Use(func(ctx *arc.Context) {
		ctx.Res.Write([]byte("hello, " + ctx.Req.PathParams.Get("name")))
	})
	http.ListenAndServe("127.0.0.1:9999", a)
}
