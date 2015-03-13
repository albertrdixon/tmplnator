package generator

import (
  "github.com/albertrdixon/tmplnator/backend"
  l "github.com/albertrdixon/tmplnator/logger"
  "os"
  "strings"
)

type Context struct {
  store backend.Backend
}

func newContext(be backend.Backend) *Context {
  return &Context{be}
}

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
