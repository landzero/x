package arcorm

import (
	"landzero.net/x/database/orm"
	"landzero.net/x/net/arc"
	"net/url"
	"os"
)

// SackKey key in sack.Sack
const SackKey = "arc.orm"

// Options orm options
type Options struct {
	Dialect     string
	DatabaseURL string
}

func sanitizeOptions(opts ...Options) (opt Options) {
	if len(opts) > 0 {
		opt = opts[0]
	}
	if len(opt.DatabaseURL) == 0 {
		opt.DatabaseURL = os.Getenv("DATABASE_URL")
	}
	if len(opt.Dialect) == 0 {
		opt.Dialect = os.Getenv("DATABASE_DIALECT")
	}
	if len(opt.Dialect) == 0 {
		if pu, err := url.Parse(opt.DatabaseURL); err == nil {
			opt.Dialect = pu.Scheme
		}
	}
	return
}

// Install install *orm.DB to *arc.Context
func Install(opts ...Options) arc.Handler {
	opt := sanitizeOptions(opts...)
	db, _ := orm.Open(opt.Dialect, opt.DatabaseURL)
	return func(ctx *arc.Context) {
		ctx.Set(SackKey, db)
		ctx.Next()
	}
}

// Extract extract *orm.DB from *arc.Context
func Extract(ctx *arc.Context) *orm.DB {
	return ctx.Value(SackKey).(*orm.DB)
}
