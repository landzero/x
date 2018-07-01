package i18n

import (
	"sync"
	"testing"

	"landzero.net/x/log"
	"landzero.net/x/runtime/binfs"
)

func init() {
	binfs.Load(&binfs.Chunk{
		Path: []string{"locales", "en-US.yml"},
		Data: []byte("hello: world"),
	})
}

func TestSourceBinFS(t *testing.T) {
	s := Source{
		dir:   "locales",
		binfs: true,
		data:  map[string]string{},
		l:     &sync.RWMutex{},
	}
	log.Println(s.Get("en-US.hello"))
}
