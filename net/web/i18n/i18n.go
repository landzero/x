package i18n

import (
	"fmt"
	"strings"
	"sync"

	"golang.org/x/text/language"
	"landzero.net/x/log"
	"landzero.net/x/net/web"
)

// I18n the i18n interface
type I18n interface {
	// T get the translation value from key
	T(key string, args ...string) string
	// Locale get the active language locale
	Locale() string
	// LocaleName get the active language locale name
	LocaleName() string
	// Locales get all language locales
	Locales() []string
	// LocaleNames get all language locale names
	LocaleNames() []string
}

// Options i18n options
type Options struct {
	// Directory directory contains i18n files, default "locales"
	Directory string
	// BinFS enable binfs support
	BinFS bool
	// Locales locales, first is default
	Locales []string
	// LocaleNames locale names
	LocaleNames []string
	// CookieName cookie name for locale overriding, default to "lang"
	CookieName string
	// QueryName query name for locale overriding, default to "lang"
	QueryName string
}

type i18n struct {
	ctx *web.Context
	opt Options
	mch language.Matcher
	src *Source
	l   string
	ln  string
}

func (in *i18n) setup() {
	lang, _ := language.MatchStrings(
		in.mch,
		in.ctx.GetCookie(in.opt.CookieName),
		in.ctx.Query(in.opt.QueryName),
		in.ctx.Header().Get("Accept-Language"),
	)
	in.l = lang.String()
	in.ctx.MapTo(in, (*I18n)(nil))
	in.ctx.Data["I18n"] = in
	in.ctx.Data["Lang"] = in.l
	for i, n := range in.opt.Locales {
		if n == in.l {
			in.ln = in.opt.LocaleNames[i]
			in.ctx.Data["LangName"] = in.ln
			break
		}
	}
}

func (in *i18n) T(key string, args ...string) string {
	fk := in.l + "." + key
	v := in.src.Get(fk)
	if len(v) == 0 {
		log.Println("i18n: missing", fk)
		return "[i18n missing: " + fk + "]"
	}
	if len(args) > 0 {
		for i, a := range args {
			v = strings.Replace(v, fmt.Sprintf("{{%d}}", i), a, -1)
		}
	}
	return v
}

func (in *i18n) Locale() string {
	return in.l
}

func (in *i18n) LocaleName() string {
	return in.ln
}

func (in *i18n) Locales() []string {
	return in.opt.Locales
}

func (in *i18n) LocaleNames() []string {
	return in.opt.LocaleNames
}

func extractOptions(opts ...Options) (opt Options) {
	if len(opts) > 0 {
		opt = opts[0]
	}
	if len(opt.Directory) == 0 {
		opt.Directory = "locales"
	}
	if len(opt.Locales) != len(opt.LocaleNames) {
		panic("i18n.Options len(opt.Locales) != len(opt.LocaleNames)")
	}
	if len(opt.Locales) == 0 {
		opt.Locales = []string{"en-US"}
		opt.LocaleNames = []string{"English"}
	}
	if len(opt.QueryName) == 0 {
		opt.QueryName = "lang"
	}
	if len(opt.CookieName) == 0 {
		opt.CookieName = "lang"
	}
	return
}

// I18ner create a i18n middleware
func I18ner(opts ...Options) web.Handler {
	// create options
	opt := extractOptions(opts...)
	// create matcher
	tags := []language.Tag{}
	for _, l := range opt.Locales {
		tags = append(tags, language.MustParse(l))
	}
	mch := language.NewMatcher(tags)
	// create source
	src := &Source{
		binfs: opt.BinFS,
		dir:   opt.Directory,
		data:  map[string]string{},
		l:     &sync.RWMutex{},
	}
	return func(ctx *web.Context) {
		if ctx.IsDevelopment() {
			src.Reload()
		}
		in := &i18n{ctx: ctx, opt: opt, src: src, mch: mch}
		in.setup()
		ctx.Next()
	}
}
