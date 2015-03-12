package config

import (
  "fmt"
)

// Config is the config struct
type Config struct {
  gofigure    interface{} `envPrefix:"T2" order:"flag,env"`
  TmplDir     string      `env:"TMPL_DIR" flag:"template-dir" flagDesc:"Template directory"`
  DefaultDir  string      `env:"DEFAULT_DIR" flag:"default-dir" flagDesc:"Default output directory"`
  Delete      bool        `env:"DELETE" flag:"delete" flagDesc:"Remove templates after processing"`
  Threads     int         `env:"THREADS" flag:"threads" flagDesc:"Number of processing threads"`
  BpoolSize   int         `env:"BPOOL_SIZE" flag:"bpool-size" flagDesc:"Size of write buffer pool"`
  EtcdPeers   []string    `env:"ETCD_PEERS" flag:"etcd-peers" flagDesc:"etcd peers in host:port (can be provided multiple times)"`
  Namespace   string      `env:"NAMESPACE" flag:"namespace" flagDesc:"etcd key namespace"`
  Verbosity   int         `env:"VERBOSITY" flag:"v" flagDesc:"Verbosity (0:quiet output, 1:default, 2:debug output)"`
  ShowVersion boolflag    `flag:"version" flagDesc:"show version"`
}

type boolflag bool

func (b *boolflag) String() string {
  return fmt.Sprintf("%t", *b)
}

func (b *boolflag) Set(v string) error {
  if v == "true" {
    *b = true
  } else {
    *b = false
  }
  return nil
}

func (b *boolflag) IsBoolFlag() bool {
  return true
}

var (
  // Build is passed in via ldflags
  Build string
)

func RuntimeVersion(version string, build string) string {
  var vers string
  if build != "" && len(version) > 4 && version[len(version)-4:] == "-dev" {
    vers = fmt.Sprintf("%s-%s", version, build)
  } else {
    vers = version
  }
  return vers
}
