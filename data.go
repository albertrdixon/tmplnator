package tmplnator

import (
	"fmt"
	"os"
	"strings"

	l "github.com/Sirupsen/logrus"
	"github.com/albertrdixon/tmplnator/backend"
)

// Data objects are passed to templates as the interface{}
type Data struct {
	Env    map[string]string
	prefix string
	store  backend.Backend
}

// NewData returns an instantiated Data object
func NewData(be backend.Backend) *Data {
	return &Data{envMap(), "", be}
}

// Get returns the Value associated with Key from the Backend or from ENV
func (d *Data) Get(key string) string {
	l.Debugf("Lookup key %q", key)

	if d.store != nil {
		k := fmt.Sprintf("%s/%s", d.prefix, strings.ToLower(key))
		if val, err := d.store.Get(k); err == nil {
			l.Debugf("Found key %q in backend: %q", key, val)
			return val
		}
	}

	l.Debugf("Did not find %q in backend, will look in ENV", key)
	if v, ok := d.Env[strings.ToUpper(strings.Replace(key, "/", "_", -1))]; ok {
		return v
	}
	return ""
}

func (d *Data) KeyPrefix(p string) { d.prefix = p }

func envMap() map[string]string {
	env := make(map[string]string, len(os.Environ()))
	for _, val := range os.Environ() {
		index := strings.Index(val, "=")
		env[val[:index]] = val[index+1:]
	}
	return env
}
