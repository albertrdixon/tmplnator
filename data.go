package tmplnator

import (
	"os"
	"strings"

	l "github.com/Sirupsen/logrus"
	"github.com/albertrdixon/tmplnator/backend"
)

// Data objects are passed to templates as the interface{}
type Data struct {
	Env   map[string]string
	store backend.Backend
}

// NewData returns an instantiated Data object
func NewData(be backend.Backend) *Data {
	return &Data{envMap(), be}
}

// Get returns the Value associated with Key from the Backend or from ENV
func (d *Data) Get(key string) string {
	l.WithField("key", key).Debug("Lookup key")

	if d.store != nil {
		key = strings.ToLower(key)
		if val, err := d.store.Get(key); err == nil {
			l.WithFields(l.Fields{
				"key":   key,
				"value": val,
			}).Debug("Found in backend")
			return val
		}
	}

	l.Debugf("Did not find %q in backend, will look in ENV", key)
	return d.Env[strings.ToUpper(strings.Replace(key, "/", "_", -1))]
}

func envMap() map[string]string {
	env := make(map[string]string, len(os.Environ()))
	for _, val := range os.Environ() {
		index := strings.Index(val, "=")
		env[val[:index]] = val[index+1:]
	}
	return env
}
