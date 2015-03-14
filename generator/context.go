package generator

import (
  l "github.com/Sirupsen/logrus"
  "github.com/albertrdixon/tmplnator/backend"
  "os"
  "strings"
)

// Context type objects are passed into the template during template.Execute().
type Context struct {
  store backend.Backend
}

func newContext(be backend.Backend) *Context {
  return &Context{be}
}

// Get performs a lookup of the given key in the backend. Failing that,
// it attempts to find the key in ENV.
func (c *Context) Get(key string) string {
  l.WithField("key", key).Debug("Lookup key")

  if c.store != nil {
    key = strings.ToLower(key)
    if val, err := c.store.Get(key); err == nil {
      l.WithFields(l.Fields{
        "key":   key,
        "value": val,
      }).Debug("Found in backend")
      return val
    }
  }

  l.WithField("key", key).Debug("Not in backend, looking in ENV")
  return os.Getenv(strings.ToUpper(strings.Replace(key, "/", "_", -1)))
}
