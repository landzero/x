package main

import (
	"fmt"
	"landzero.net/x/database/orm"
	_ "landzero.net/x/database/orm/dialects/postgres"
	"landzero.net/x/net/arc"
	"landzero.net/x/net/arc/database"
	"net/http"
)

type User struct {
	orm.Model
}

func main() {
	a := arc.New()
	// global orm.DB instance, using DATABASE_URL
	a.Install(database.Installer())
	// per request orm.DB instance
	// a.Use(database.Handler())
	a.Route().Get("/").Use(func(ctx *arc.Context) {
		db := database.Extract(ctx)
		var count int
		db.Model(&User{}).Count(&count)
		ctx.Res.Write([]byte("hello, static " + fmt.Sprintf("%d", count)))
	})
	a.Route().Get("/hello/:name").Use(func(ctx *arc.Context) {
		ctx.Res.Write([]byte("hello, " + ctx.Req.PathParams.Get("name")))
	})
	http.ListenAndServe("127.0.0.1:9999", a)
}
