package tmplnator

import (
	"bytes"
	"os"
	"testing"

	"github.com/albertrdixon/tmplnator/backend"
	"github.com/spf13/afero"
)

var testDir = "fixtures/test/"
var expected = []byte("Generated")

func TestGenerate(t *testing.T) {
	if testing.Verbose() {
		LogLevel("debug")
	} else {
		LogLevel("fatal")
	}

	os.Setenv("FOO", "Generated")
	b := backend.NewMock(map[string]string{"/test/vars/bif": "Generated"}, nil)
	Backend = b
	t.Logf("Running MemGen(%s)", testDir)
	files := MemGen(testDir)
	if len(files) < 5 {
		t.Errorf("Expected 5 files generated, got %d", len(files))
	} else {
		t.Logf("Checking Files: %v", files)
		for _, file := range files {
			t.Logf("Checking File: %v", file)
			if file.HasErrs() {
				t.Errorf("%q: Got errors!", file.TemplateName())
				for _, err := range file.errs {
					t.Errorf("%v", err)
				}
			} else {
				if _, err := destFs.Stat(file.FullPath()); err != nil {
					t.Errorf("%q: File not generated! %v", file.Name(), err)
				}

				content, err := afero.ReadFile(destFs, file.FullPath())
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
