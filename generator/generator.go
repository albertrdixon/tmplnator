package generator

import (
  "github.com/albertrdixon/tmplnator/backend"
  "github.com/albertrdixon/tmplnator/config"
  l "github.com/albertrdixon/tmplnator/logger"
  "github.com/albertrdixon/tmplnator/stack"
  "github.com/oxtoacart/bpool"
  "os"
  "sync"
)

type generator struct {
  files      *stack.Stack
  defaultDir string
  context    *Context
  bpool      *bpool.BufferPool
  wg         *sync.WaitGroup
  threads    int
  del        bool
}

// NewGenerator returns a generator with a parsed file stack
func NewGenerator(c *config.Config) (*generator, error) {
  fs, err := parseFiles(c.TmplDir, c.DefaultDir)
  if err != nil {
    return nil, err
  }

  return &generator{
    files:      fs,
    defaultDir: c.DefaultDir,
    context:    newContext(backend.New(c.Namespace, c.EtcdPeers)),
    bpool:      bpool.NewBufferPool(c.BpoolSize),
    threads:    c.Threads,
    wg:         new(sync.WaitGroup),
    del:        c.Delete,
  }, nil
}

// Generate kicks off file generation. Will spin out generator.threads
// number of goroutines running generator.process()
func (g *generator) Generate() (err error) {
  l.Info("Generating files (threads=%d)", g.threads)
  g.wg.Add(g.threads)

  for i := 0; i < g.threads; i++ {
    go g.process()
  }

  g.wg.Wait()
  l.Quiet(g.defaultDir)
  return nil
}

func (g *generator) process() {
  defer g.wg.Done()

  l.Debug("files.Len(): %d", g.files.Len())
  for g.files.Len() > 0 {
    if f, ok := g.files.Pop().(*file); ok {
      if err := g.write(f); err == nil {
        os.Chown(f.destination(), f.user, f.group)
        if g.del {
          l.Info("Removing %q", f.src)
          if err := os.Remove(f.src); err != nil {
            l.Info("Problem in remove: %v", err)
          }
        }
      } else {
        l.Info("Problem in write: %v", err)
      }
    } else {
      panic("Internal Error: Could not cast stack item as file")
    }
  }
}

func (g *generator) write(f *file) error {
  b := g.bpool.Get()
  defer g.bpool.Put(b)

  err := f.body.Execute(b, g.context)
  if err != nil {
    return err
  }

  l.Info("Generating %q from %q", f.destination(), f.src)
  if err := os.MkdirAll(f.dir, 0755); err != nil {
    return err
  }
  fh, err := os.OpenFile(f.destination(), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, f.mode)
  if err != nil {
    return err
  }
  defer fh.Close()

  _, err = fh.Write(b.Bytes())
  if err != nil {
    return err
  }
  return nil
}
