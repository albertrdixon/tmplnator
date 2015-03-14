package generator

import (
  "fmt"
  l "github.com/Sirupsen/logrus"
  "github.com/albertrdixon/tmplnator/backend"
  "github.com/albertrdixon/tmplnator/config"
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
  l.WithField("threads", g.threads).Info("Generating files")
  g.wg.Add(g.threads)

  for i := 0; i < g.threads; i++ {
    go g.process(i)
  }

  g.wg.Wait()
  if l.GetLevel() == l.ErrorLevel {
    fmt.Println(g.defaultDir)
  }
  return nil
}

func (g *generator) process(id int) {
  l.WithFields(l.Fields{
    "id":              id,
    "file_stack_size": g.files.Len(),
  }).Debug("Starting processing thread")
  defer g.wg.Done()

  for g.files.Len() > 0 {
    if f, ok := g.files.Pop().(*file); ok {
      l.WithFields(l.Fields{
        "id":       id,
        "template": f.src,
      }).Debug("Processing a template")
      if err := g.write(f); err == nil {
        os.Chown(f.destination(), f.user, f.group)
        if g.del {
          l.WithField("path", f.src).Info("Removing file")
          if err := os.Remove(f.src); err != nil {
            l.WithField("error", err).Info("Failed to remove file")
          }
        }
      } else {
        l.WithField("error", err).Info("Failed to write file")
      }
    } else {
      l.WithField("item", f).Panic("Internal Error: Could not cast stack item as *file")
    }
  }
}

func (g *generator) write(f *file) error {
  b := g.bpool.Get()
  defer g.bpool.Put(b)

  l.WithFields(l.Fields{
    "template": f.src,
    "context":  g.context,
  }).Debug("Executing template")
  err := f.body.Execute(b, g.context)
  if err != nil {
    return err
  }

  if f.skip {
    return nil
  }

  l.WithField("path", f.dir).Debug("Creating directory")
  if _, err := os.Stat(f.dir); err != nil {
    if err = os.MkdirAll(f.dir, f.dirmode); err != nil {
      return err
    }
  }

  l.WithField("path", f.destination()).Debug("Creating file")
  fh, err := os.OpenFile(f.destination(), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, f.mode)
  if err != nil {
    return err
  }
  defer fh.Close()

  l.WithFields(l.Fields{
    "template": f.src,
    "file":     f.destination(),
  }).Info("Generating file")
  _, err = fh.Write(b.Bytes())
  if err != nil {
    return err
  }
  return nil
}
