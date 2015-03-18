package file

import (
  // "io/ioutil"
  "bytes"
  "os"
  "testing"
)

var filetest = []struct {
  name           string
  template       string
  expectedOutput string
  expectedInfo   Info
  expectError    bool
  stackSize      int
}{
  {
    name:           "bad",
    template:       `{{ dir "/some/other/path" }{{ mode 0755 "one too many" }}Body Text {{ env "BAD" Something }}`,
    expectedOutput: "",
    expectedInfo:   Info{},
    expectError:    true,
    stackSize:      0,
  },
  {
    name:           "change_everything",
    template:       `{{ dir "/some/path" }}{{ name "name_changed" }}{{ mode 0777 }}{{ user 10000 }}Body Text`,
    expectedOutput: "Body Text",
    expectedInfo: Info{
      Name: "name_changed",
      Dir:  "/some/path",
      Mode: os.FileMode(0777),
      User: 10000,
    },
    expectError: false,
    stackSize:   1,
  },
}

func TestParseFile(t *testing.T) {
  Testing = true
  for _, ft := range filetest {
    fq := NewFileQueue()
    mf := NewFile(ft.template, ft.name)
    err := ParseFile(mf, fq)
    fq.PopulateQueue()

    if !ft.expectError && err != nil {
      t.Errorf("ParseFile(%q): Expected no error while parsing, got: %v", ft.name, err)
    }
    if ft.expectError && err == nil {
      t.Errorf("ParseFile(%q): Expected an error while parsing", ft.name)
    }
    if fq.Len() != ft.stackSize {
      t.Errorf("ParseFile(%q): Expected stack size to be %d, got %d", ft.name, ft.stackSize, fq.Len())
    }
  }
}

func TestWriteFile(t *testing.T) {
  Testing = true
  for _, ft := range filetest {
    fq := NewFileQueue()
    mf := NewFile(ft.template, ft.name)
    err := ParseFile(mf, fq)
    if err != nil {
      if !ft.expectError {
        t.Errorf("WriteFile(%q): Parsing failed, please fix it.", ft.name)
      }
    } else {
      err = mf.Write(new(bytes.Buffer), nil)
      if err != nil {
        t.Errorf("WriteFile(%q): Did not expect error in write: %v", ft.name, err)
      }

      out, info := mf.Output(), mf.Info()
      if out != ft.expectedOutput {
        t.Errorf("WriteFile(%q): Expected output=%q, got output=%q", ft.name, ft.expectedOutput, out)
      }

      if info.Name != ft.expectedInfo.Name {
        t.Errorf("WriteFile(%q): Expected filename=%q, got filename=%q", ft.name, ft.expectedInfo.Name, info.Name)
      }
      if info.Dir != ft.expectedInfo.Dir {
        t.Errorf("WriteFile(%q): Expected dir=%q, got dir=%q", ft.name, ft.expectedInfo.Dir, info.Dir)
      }
      if info.Mode != ft.expectedInfo.Mode {
        t.Errorf("WriteFile(%q): Expected mode=%q, got mode=%q", ft.name, ft.expectedInfo.Mode, info.Mode)
      }
      if info.User != ft.expectedInfo.User {
        t.Errorf("WriteFile(%q): Expected user=%d, got user=%d", ft.name, ft.expectedInfo.User, info.User)
      }
    }
  }
}
