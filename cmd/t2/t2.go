package main

import (
  "fmt"
  "github.com/albertrdixon/tmplnator/config"
  "github.com/albertrdixon/tmplnator/generator"
  l "github.com/albertrdixon/tmplnator/logger"
  "github.com/ian-kent/gofigure"
  "io/ioutil"
  "os"
  "path/filepath"
)

func main() {
  var cfg = config.Config{
    TmplDir:     "/templates",
    DefaultDir:  "",
    Delete:      false,
    Threads:     4,
    BpoolSize:   4,
    Verbosity:   1,
    ShowVersion: false,
  }

  // gofigure.Debug = true
  err := gofigure.Gofigure(&cfg)
  if err != nil {
    fmt.Printf("Error parsing config: %v\n", err)
    os.Exit(1)
  }

  if cfg.ShowVersion {
    fmt.Println(config.RuntimeVersion(config.CodeVersion, config.Build))
    os.Exit(0)
  }
  switch {
  case cfg.Verbosity <= 0:
    l.Level = 0
  case cfg.Verbosity == 1:
    l.Level = 1
  case cfg.Verbosity >= 2:
    l.Level = 2
  }

  if _, err := os.Stat(cfg.TmplDir); err != nil {
    fmt.Printf("Problems reading dir %q: %v\n", cfg.TmplDir, err)
    os.Exit(2)
  }
  if cfg.DefaultDir == "" {
    d, err := ioutil.TempDir("", "T2")
    if err != nil {
      d = filepath.Join(os.TempDir(), "T2")
    }
    cfg.DefaultDir = d
  }

  // be := backend.New(cfg.Namespace, cfg.EtcdPeers)
  g, err := generator.NewGenerator(&cfg)
  if err != nil {
    fmt.Printf("ERROR: %v", err)
  }

  if err := g.Generate(); err != nil {
    fmt.Printf("ERROR: %v", err)
  }
  os.Exit(0)
}
