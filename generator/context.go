package generator

import (
  "github.com/albertrdixon/tmplnator/backend"
  l "github.com/albertrdixon/tmplnator/logger"
  "os"
  "strings"
)

type Context struct {
  Env   map[string]string
  store backend.Backend
}

func newContext(be backend.Backend) *Context {
  return &Context{
    Env:   envMap(),
    store: be,
  }
}

func (c *Context) Var(key string) string {
  if c.store != nil {
    key = strings.ToLower(key)
    if rtn, err := c.store.Get(key); err == nil {
      l.Debug("Got backend[%v]: %v", key, rtn)
      return rtn
    }
  }

  key = strings.ToUpper(strings.Replace(key, "/", "_", -1))
  l.Debug("Lookup c.Env[%v]", key)
  if rtn, ok := c.Env[key]; ok {
    l.Debug("Got: %v", rtn)
    return rtn
  }
  return ""
}

func envMap() map[string]string {
  env := make(map[string]string, len(os.Environ()))
  for _, val := range os.Environ() {
    index := strings.Index(val, "=")
    env[val[:index]] = val[index+1:]
  }
  return env
}
