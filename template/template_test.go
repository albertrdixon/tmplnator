package template

import (
  "github.com/albertrdixon/tmplnator/stack"
  "io/ioutil"
  "os"
  "path/filepath"
  "syscall"
  "testing"
)

var tmpl = `
{{ destination . "/out/%s" .Env.BAR }}
{{ mode . 0777 }}
{{ user . 0 }}
Foo = {{ .Env.BAR }}
`
var expected = `



Foo = Baz
`

func TestTemplate(t *testing.T) {
  os.Setenv("BAR", "Baz")
  env = envMap()

  prefix, err := ioutil.TempDir("", "tmpltest")
  defer os.RemoveAll(prefix)
  if err != nil {
    t.Errorf("ParseTemplate(): Could not create tmp dir: %v", err)
  }

  file, targetDir := filepath.Join(prefix, "test"), filepath.Join(prefix, "out", "Baz")
  os.Mkdir(targetDir, 0750)
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

  f, err = os.Open(filepath.Join(targetDir, "test"))
  if err != nil {
    t.Errorf("template.Write(): problems opening written file: %v", err)
  }
  info, err := f.Stat()
  if err != nil {
    t.Errorf("template.Write(): problems gettings stats for written file: %v", err)
  }

  var body = make([]byte, info.Size())
  _, err = f.Read(body)
  if err != nil {
    t.Errorf("template.Write(): problems reading written file: %v", err)
  }

  if string(body) != expected {
    t.Errorf("expected %q, got %q", expected, body)
  }
  if info.Mode() != os.FileMode(0777) {
    t.Logf("expected file mode %v, got %v, umask likely interfered", os.FileMode(0777), info.Mode())
  }
  if info.Sys().(*syscall.Stat_t).Uid != 0 {
    t.Logf("expected uid %d, got %d, I likely do not have permissions", 0, info.Sys().(*syscall.Stat_t).Uid)
  }
}
