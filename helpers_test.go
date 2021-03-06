package tmplnator

import (
	"os"
	"testing"

	"github.com/albertrdixon/tmplnator/backend"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderTemplate(te *testing.T) {
	is := assert.New(te)
	must := require.New(te)
	if testing.Verbose() {
		LogLevel("debug")
	}
	initFs(true, true)

	path := "template"
	os.Setenv("foo", "bar")
	os.Setenv("slice", "0,1,2")
	os.Setenv("json", `{"key": "value"}`)
	Backend = backend.NewMock(map[string]string{"/test/vars/bif": "Generated"}, nil)
	data = NewData(nil)
	var tests = []struct {
		name    string
		content string
		out     string
		pass    bool
	}{
		{"from_json", `{{ get (from_json .Env.json) "key" }}`, "value", true},
		{"has_key (yes)", `{{ printf "%t" (has_key .Env "foo") }}`, "true", true},
		{"has_key (no)", `{{ printf "%t" (has_key .Env .Env.foo) }}`, "false", true},
		{"has_key (fail)", `{{ printf "%t" (has_key .Env 2) }}`, "false", true},
		{"titleize", `{{ titleize .Env.foo }}`, "Bar", true},
		{"titleize (empty)", `{{ titleize .Env.bar }}`, "", true},
		{"downcase", `{{ downcase .Env.foo }}`, "bar", true},
		{"upcase", `{{ upcase .Env.foo }}`, "BAR", true},
		{"trim", `{{ trim .Env.foo "b" }}`, "ar", true},
		{"titleize + downcase", `{{ titleize (downcase .Env.foo) }}`, "Bar", true},
		{"eq", `{{ printf "%t" (eq .Env.empty "fail") }}`, "false", false},
		{"eql", `{{ printf "%t" (eql .Env.bar "fail") }}`, "false", true},
		{"eq (fix)", `{{ printf "%t" (eq (def .Env.bar "baz") "baz") }}`, "true", true},
		{"get & split", `{{ get (split "0,1,2" ",") 1 }}`, "1", true},
		{"get & split (2)", `{{ get (split .Env.slice ",") 2 }}`, "2", true},
		{"contains", `{{ printf "%t" (contains .Env.foo .Env.slice) }}`, "false", true},
		{"has_suffix", `{{ printf "%t" (has_suffix .Env.slice "2") }}`, "true", true},
		{"split & join", `{{ join (split .Env.slice ",") "-" }}`, "0-1-2", true},
		{"join (fail)", `{{ join (split .Env.slice ",") 2 }}`, "", false},
		{"truncate", `{{ truncate 3 "foobar" }}`, "foo", true},
		{"truncate 2", `{{ truncate 5 "foo" }}`, "foo", true},
	}

	for i, t := range tests {
		fh, err := srcFs.Create(path)
		must.NoError(err)
		_, err = fh.WriteString(t.content)
		must.NoError(err)

		f := NewFile(path)
		f.Info().SetFullpath(t.name)
		must.NoError(ParseTemplate(f, srcFs))

		err = WriteFile(f, destFs)
		if t.pass {
			is.NoError(err, "[%d %s]", i, t.name)
			b, err := afero.ReadFile(destFs, f.FullPath())
			is.NoError(err, "[%d %s]", i, t.name)
			is.Equal(t.out, string(b), "[%d %s]", i, t.name)
		} else {
			is.Error(err, "[%d %s]", i, t.name)
		}
		must.NoError(srcFs.Remove(path))
		if t.pass {
			must.NoError(destFs.Remove(t.name))
		}
	}
}
