package generator

import (
  "fmt"
  l "github.com/Sirupsen/logrus"
  "github.com/albertrdixon/tmplnator/backend"
  "github.com/albertrdixon/tmplnator/config"
  "github.com/albertrdixon/tmplnator/file"
  "github.com/oxtoacart/bpool"
  "sync"
)

type generator struct {
  // files      *stack.Stack
  files      *file.FileQueue
  defaultDir string
  context    *Context
  bpool      *bpool.BufferPool
  wg         *sync.WaitGroup
  threads    int
  del        bool
  errors     chan error
}

// NewGenerator returns a generator with a parsed file stack
func NewGenerator(c *config.Config) (*generator, error) {
  fq, err := file.ParseFiles(c.TmplDir, c.DefaultDir)
  if err != nil {
    return nil, err
  }

  return &generator{
    files:      fq,
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
  if l.GetLevel() <= l.ErrorLevel {
    fmt.Println(g.defaultDir)
  }
  return nil
}

func (g *generator) process(id int) {
  l.WithFields(l.Fields{
    "thread_id":       id,
    "file_stack_size": g.files.Len(),
  }).Debug("Starting processing thread")
  defer g.wg.Done()
  defer g.catch(id)

  for f := range g.files.Queue() {
    l.WithFields(l.Fields{
      "thread_id": id,
      "template":  f.Src(),
    }).Debug("Processing template")
    if err := g.write(f); err == nil {
      if g.del {
        l.WithField("template", f.Src()).Info("Removing template")
        if err := f.DeleteTemplate(); err != nil {
          l.WithField("error", err).Error("Failed to remove file")
        }
      }
    } else {
      l.WithField("error", err).Fatal("Failed to write file")
    }
  }
}

func (g *generator) write(f file.File) (err error) {
  b := g.bpool.Get()
  defer g.bpool.Put(b)

  err = f.Write(b, g.context)
  return
}

func (g *generator) catch(tid int) {
  if r := recover(); r != nil {
    l.WithFields(l.Fields{
      "thread_id": tid,
      "message":   r,
    }).Fatal("Recovered from panic!")
  }
}
