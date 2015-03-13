package main

import (
  "fmt"
  "github.com/albertrdixon/tmplnator/config"
  "github.com/albertrdixon/tmplnator/generator"
  l "github.com/albertrdixon/tmplnator/logger"
  "github.com/ian-kent/gofigure"
  "io/ioutil"
  "os"
)

func main() {
  var cfg = config.Defaults

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
  if d, err := ioutil.TempDir("", "T2"); err == nil {
    cfg.DefaultDir = d
  }

  g, err := generator.NewGenerator(&cfg)
  if err != nil {
    fmt.Printf("ERROR: %v", err)
  }

  if err := g.Generate(); err != nil {
    fmt.Printf("ERROR: %v", err)
  }
  os.Exit(0)
}
