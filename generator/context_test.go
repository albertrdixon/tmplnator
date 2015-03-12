package generator

import (
  "github.com/albertrdixon/tmplnator/backend"
  "os"
  "testing"
)

var mb backend.Backend

func TestVar(t *testing.T) {
  var vartests = []struct {
    key      string
    expected string
  }{
    {"foo", "bar"},
    {"one", "two"},
    {"foo/bar", "baz"},
    {"john/jay", "signer"},
    {"not/there", ""},
  }

  c := newContext(mb)
  for _, vt := range vartests {
    out := c.Var(vt.key)
    if out != vt.expected {
      t.Errorf("Var(key=%s): Expected %q, got %q", vt.key, vt.expected, out)
    }
  }
}

func init() {
  mb = backend.NewMock(
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
  os.Setenv("ONE", "THREE")
  os.Setenv("JOHN_JAY", "signer")
}
