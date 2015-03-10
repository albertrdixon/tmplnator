package template

import (
  "github.com/albertrdixon/tmplnator/stack"
  "io/ioutil"
  "os"
  "path/filepath"
  "testing"
)

var tmpl = `
{{ destination . "/out"}}
Foo = {{ .Vars.BAR }}
`
var expected = `

Foo = Baz
`

func TestTemplate(t *testing.T) {
  os.Setenv("BAR", "Baz")
  varMap = envMap()

  prefix, err := ioutil.TempDir("", "tmpltest")
  defer os.RemoveAll(prefix)
  if err != nil {
    t.Errorf("ParseTemplate(): Could not create tmp dir: %v", err)
  }

  file, targetDir := filepath.Join(prefix, "test"), filepath.Join(prefix, "out")
  os.Mkdir(targetDir, 0755)
  f, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
  if err != nil {
    t.Errorf("ParseTemplate(): Could not create testFile: %v", err)
  }
  f.WriteString(tmpl)
  f.Close()

  st := stack.NewStack()
  err = ParseTemplate(file, prefix, st)
  if err != nil {
    t.Errorf("ParseTemplate(): did not expect an error: %v", err)
  }

  tm, ok := st.Pop().(*Template)
  if !ok {
    t.Error("ParseTemplate(): Could not cast stack item as Template")
  }

  err = tm.Write()
  if err != nil {
    t.Errorf("template.Write(): did not expect error: %v", err)
  }

  out, err := ioutil.ReadFile(filepath.Join(targetDir, "test"))
  if err != nil {
    t.Errorf("template.Write(): problems reading written file: %v", err)
  }

  if string(out) != expected {
    t.Errorf("expected %q, got %q", expected, out)
  }
}
