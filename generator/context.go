package generator

import (
  "github.com/albertrdixon/tmplnator/backend"
  l "github.com/albertrdixon/tmplnator/logger"
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
  if c.store != nil {
    key = strings.ToLower(key)
    if rtn, err := c.store.Get(key); err == nil {
      l.Debug("Got backend[%v]: %v", key, rtn)
      return rtn
    }
  }
  l.Debug("No backend configured, looking up %q in ENV", key)
  return os.Getenv(strings.ToUpper(strings.Replace(key, "/", "_", -1)))
}
