package main

import "landzero.net/x/net/web"

func main() {
	m := web.Classic()
	m.Use(web.Renderer())
	m.Get("/", func(ctx *web.Context) {
		ctx.PlainText(200, []byte("Hello, World"))
	})
	m.RunIvy()
}
