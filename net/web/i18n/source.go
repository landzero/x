package i18n

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"

	"landzero.net/x/com"

	"landzero.net/x/runtime/binfs"

	"landzero.net/x/encoding/yaml"
	"landzero.net/x/log"
)

// Source i18n source
type Source struct {
	loaded bool
	data   map[string]string
	dir    string
	binfs  bool
	l      *sync.RWMutex
}

func flatten(pfx string, in map[interface{}]interface{}, out map[string]string) {
	for key, val := range in {
		keyStr := com.ToStr(key)
		if valMap, ok := val.(map[interface{}]interface{}); ok {
			flatten(pfx+keyStr+".", valMap, out)
		} else {
			out[pfx+keyStr] = com.ToStr(val)
		}
	}
}

func (s *Source) loadYaml(pfx string, buf []byte) {
	d := map[interface{}]interface{}{}
	if err := yaml.Unmarshal(buf, &d); err != nil {
		log.Println("i18n: failed to load yaml", pfx, err)
		return
	}
	flatten(pfx, d, s.data)
}

func (s *Source) load() {
	var f http.File
	var err error
	// create http.File for directory
	if s.binfs {
		if f, err = binfs.Open(s.dir); err != nil {
			log.Println("i18n: failed to open binfs dir", s.dir, err)
			return
		}
	} else {
		if f, err = os.Open(s.dir); err != nil {
			log.Println("i18n: failed to open dir", s.dir, err)
			return
		}
	}
	var st os.FileInfo
	if st, err = f.Stat(); err != nil {
		log.Println("i18n: failed to stat dir", s.dir, err)
		return
	}
	if !st.IsDir() {
		log.Println("i18n:", s.dir, "is not a directory")
		return
	}
	// iterate files
	var fs []os.FileInfo
	if fs, err = f.Readdir(0); err != nil {
		log.Println("i18n: failed to readdir", s.dir)
	}
	for _, fi := range fs {
		if fi.IsDir() {
			continue
		}
		name := fi.Name()
		ext := path.Ext(name)
		if ext != ".yml" && ext == ".yaml" {
			continue
		}
		var f http.File
		if s.binfs {
			fp := path.Join(s.dir, name)
			if f, err = binfs.Open(fp); err != nil {
				log.Println("i18n: failed to open binfs", fp)
				continue
			}
		} else {
			fp := filepath.Join(s.dir, name)
			if f, err = os.Open(fp); err != nil {
				log.Println("i18n: failed to open", fp)
				continue
			}
		}
		defer f.Close()
		buf, err := ioutil.ReadAll(f)
		if err != nil {
			log.Println("i18n: failed to read", name)
			continue
		}
		s.loadYaml(name[:len(name)-len(ext)]+".", buf)
	}
}

func (s *Source) loadIfNeeded() {
	if s.loaded {
		return
	}
	s.l.Lock()
	defer s.l.Unlock()
	s.load()
	s.loaded = true
}

// Get get a value by key
func (s *Source) Get(key string) string {
	s.loadIfNeeded()
	s.l.RLock()
	defer s.l.RUnlock()
	return s.data[key]
}

// Reload reload all data
func (s *Source) Reload() {
	s.l.Lock()
	defer s.l.Unlock()
	s.loaded = false
	s.data = map[string]string{}
}
