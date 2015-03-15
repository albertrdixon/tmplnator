package file

import (
  "github.com/albertrdixon/tmplnator/stack"
  "io/ioutil"
  "os"
  "path/filepath"
  "testing"
)

func TestParseFile(t *testing.T) {
  var filetest = []struct {
    name        string
    body        string
    expectError bool
    stackSize   int
  }{
    {
      name:        "good",
      body:        `{{ dir "/some/path" }}{{ mode 0777 }} Body Text {{ .Env.VAR }}`,
      expectError: false,
      stackSize:   1,
    },
    {
      name:        "bad",
      body:        `{{ dir "/some/other/path" }{{ mode 0755 "one too many" }}Body Text {{ .Env.BAD Something }}`,
      expectError: true,
      stackSize:   0,
    },
  }

  dir, err := ioutil.TempDir("", "tmpltest")
  defer os.RemoveAll(dir)
  if err != nil {
    t.Errorf("ParseFile(): Could not create tmp dir: %v", err)
  }

  for _, ft := range filetest {
    fp := filepath.Join(dir, ft.name)
    fh, err := os.OpenFile(fp, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
    if err != nil {
      t.Errorf("ParseFile(): Could not create testFile: %v", err)
    }
    fh.WriteString(ft.body)
    fh.Close()

    st := stack.NewStack()
    err = parseFile(fp, "", st)
    if !ft.expectError && err != nil {
      t.Errorf("ParseFile(%q): Expected no error while parsing, got: %v", ft.name, err)
    }
    if ft.expectError && err == nil {
      t.Errorf("ParseFile(%q): Expected an error while parsing", ft.name)
    }
    if st.Len() != ft.stackSize {
      t.Errorf("ParseFile(%q): Expected stack size to be %d, got %d", ft.name, ft.stackSize, st.Len())
    }
  }
}
