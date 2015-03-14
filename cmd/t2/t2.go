package main

import (
  "fmt"
  l "github.com/Sirupsen/logrus"
  "github.com/albertrdixon/tmplnator/config"
  "github.com/albertrdixon/tmplnator/generator"
  "github.com/ian-kent/gofigure"
  "io/ioutil"
  "os"
)

func setup() *config.Config {
  var cfg = config.Defaults

  // gofigure.Debug = true
  err := gofigure.Gofigure(&cfg)
  if err != nil {
    fmt.Printf("Error parsing config: %v\n", err)
    os.Exit(1)
  }
  registerLogger(cfg.Verbosity)

  if cfg.ShowVersion {
    fmt.Println(config.RuntimeVersion(config.CodeVersion, config.Build))
    os.Exit(0)
  }
  if _, err := os.Stat(cfg.TmplDir); err != nil {
    l.WithFields(l.Fields{
      "directory": cfg.TmplDir,
      "error":     err,
    }).Fatal("Problems reading template dir")
  }

  if d, err := ioutil.TempDir("", "T2"); err == nil {
    cfg.DefaultDir = d
  }

  return &cfg
}

func registerLogger(lvl int) {
  l.SetOutput(os.Stdout)
  switch {
  case lvl <= 0:
    l.SetLevel(l.ErrorLevel)
  case lvl == 1:
    l.SetLevel(l.InfoLevel)
  case lvl >= 2:
    l.SetLevel(l.DebugLevel)
  }
}

func generateFiles(cfg *config.Config) (err error) {
  g, err := generator.NewGenerator(cfg)
  if err != nil {
    return
  }

  err = g.Generate()
  return
}

func main() {
  cfg := setup()

  if err := generateFiles(cfg); err != nil {
    l.WithField("error", err).Fatal("ERROR!")
  }
  os.Exit(0)
}
