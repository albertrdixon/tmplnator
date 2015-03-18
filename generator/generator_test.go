package generator

import (
  "github.com/albertrdixon/tmplnator/backend"
  "github.com/albertrdixon/tmplnator/file"
  "github.com/oxtoacart/bpool"
  "os"
  "sync"
  "testing"
)

var filetest = []struct {
  names          []string
  templates      []string
  expectedOutput map[string]string
  expectError    bool
  stackSize      int
}{
  {
    names:          []string{"one"},
    templates:      []string{`{{ dir "/some/path" }}{{ mode 0777 }}Body Text {{ .Env.TEST_VAR }}`},
    expectedOutput: map[string]string{"one": "Body Text VALUE"},
    expectError:    false,
    stackSize:      1,
  },
  {
    names: []string{"first", "second"},
    templates: []string{
      `{{ dir "/some/other/path" }}{{ mode 0755 }}Body Text One {{ .Env.TEST_VAR }}`,
      `{{ dir "some/other/path" }}{{ name "2nd" }}{{ mode 0644 }}Body Text Two {{ .Get "foo/bar" }}`,
    },
    expectedOutput: map[string]string{
      "first": "Body Text One VALUE",
      "2nd":   "Body Text Two baz",
    },
    expectError: true,
    stackSize:   0,
  },
}

func newTestGenerator(fq *file.Queue, be backend.Backend) *generator {
  return &generator{
    files:      fq,
    defaultDir: "/var/tmp/testing",
    context:    newContext(be),
    bpool:      bpool.NewBufferPool(2),
    threads:    2,
    wg:         new(sync.WaitGroup),
    del:        false,
  }
}

func TestProcess(t *testing.T) {
  mb := backend.NewMock(
    map[string]string{
      "foo":     "bar",
      "foo/bar": "baz",
      "one":     "two",
    },
    map[string][]string{
      "foo":     []string{"bar", "baz"},
      "foo/baz": []string{"bim", "biff"},
    },
  )

  for _, ft := range filetest {
    file.Testing = true
    fq := file.NewFileQueue()
    for idx, tm := range ft.templates {
      mf := file.NewFile(tm, ft.names[idx])
      err := file.ParseFile(mf, fq)
      if err != nil {
        t.Errorf("Parsing failed, please fix it! %v", err)
        t.FailNow()
      }
    }

    if !t.Failed() {
      g := newTestGenerator(fq, mb)
      fq.PopulateQueue()
      err := g.Generate()
      if err != nil {
        t.Errorf("Generate(%q): Should not have produced an error: %v", ft.names, err)
      }

      for _, f := range fq.Files() {
        fi := f.Info()
        if f.Output() != ft.expectedOutput[fi.Name] {
          t.Errorf("Generate(%q): Output not expected, Got file=%q out=%q", ft.names, fi.Name, f.Output())
        }
      }
    }
  }
}

func init() {
  os.Setenv("TEST_VAR", "VALUE")
}
