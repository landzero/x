package database

import (
	"landzero.net/x/database/orm"
	"landzero.net/x/net/arc"
	"net/http"
	"net/url"
	"os"
)

// SackKey key in sack.Sack
const SackKey = "arc.orm"

// Options orm options
type Options struct {
	Dialect string
	URL     string
}

func sanitizeOptions(opts ...Options) (opt Options) {
	if len(opts) > 0 {
		opt = opts[0]
	}
	if len(opt.URL) == 0 {
		opt.URL = os.Getenv("DATABASE_URL")
	}
	if len(opt.Dialect) == 0 {
		opt.Dialect = os.Getenv("DATABASE_DIALECT")
	}
	if len(opt.Dialect) == 0 {
		if pu, err := url.Parse(opt.URL); err == nil {
			opt.Dialect = pu.Scheme
		}
	}
	return
}

// Installer arc.Installer, provide global orm.DB instance, panic if failed
func Installer(opts ...Options) arc.Installer {
	opt := sanitizeOptions(opts...)
	db, err := orm.Open(opt.Dialect, opt.URL)
	if err != nil {
		panic(err)
	}
	return func(a *arc.Arc) {
		a.Set(SackKey, db)
	}
}

// Handler arc.Handler, provider per request orm.DB instance, call http.Error if failed
func Handler(opts ...Options) arc.Handler {
	opt := sanitizeOptions(opts...)
	return func(ctx *arc.Context) {
		var err error
		var db *orm.DB
		if db, err = orm.Open(opt.Dialect, opt.URL); err != nil {
			http.Error(ctx.Res, "failed to initialize database", http.StatusInternalServerError)
			return
		}
		defer db.Close()
		ctx.Set(SackKey, db)
		ctx.Next()
	}
}

// Extract extract *orm.DB from *arc.Context
func Extract(ctx *arc.Context) *orm.DB {
	return ctx.Value(SackKey).(*orm.DB)
}
