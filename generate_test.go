package tmplnator

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/albertrdixon/tmplnator/backend"
)

var testDir = "fixtures/test/"
var expected = []byte("Generated")

func TestGenerate(t *testing.T) {
	LogLevel("fatal")
	os.Setenv("FOO", "Generated")
	b := backend.NewMock(map[string]string{"/test/vars/bif": "Generated"}, nil)
	Backend = b
	files := MemGen(testDir)
	if len(files) < 5 {
		t.Errorf("Expected 5 files generated, got %d", len(files))
	} else {
		for _, file := range MemGen(testDir) {
			if file.HasErrs() {
				t.Errorf("%q: Got errors!", file.TemplateName())
				for _, err := range file.errs {
					t.Errorf("%v", err)
				}
			} else {
				if _, err := destFs.Stat(file.FullPath()); err != nil {
					t.Errorf("%q: File not generated! %v", file.Name(), err)
				}

				content, err := ioutil.ReadAll(file)
				if err != nil {
					t.Errorf("%q: Error reading file: %v", file.Name(), err)
				}
				content = bytes.TrimSpace(content)
				if !bytes.Equal(content, expected) {
					t.Errorf("%q: Expected content to be %q, got %q", file.Name(), string(expected), string(content))
				}
			}
		}
	}
}
