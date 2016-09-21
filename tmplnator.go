package tmplnator

import (
	"os"
	"path/filepath"
	"runtime"

	l "github.com/Sirupsen/logrus"
	"github.com/albertrdixon/tmplnator/backend"
	"github.com/spf13/afero"
)

var (
	errChan chan error
	data    *Data
	srcFs   afero.Fs
	destFs  afero.Fs

	// Backend is the Key/Value store for Data
	Backend backend.Backend

	// ForceTemp will force all templates to be written to TmpDir if true
	ForceTemp = false

	// TmpDir is the tmpdir we will use
	TmpDir = filepath.Join(os.TempDir(), "T2")

	cpus = runtime.NumCPU()
	pkg  = "tmplnator"
)

func LogLevel(level string) {
	lvl, err := l.ParseLevel(level)
	if err != nil {
		lvl = l.InfoLevel
	}
	l.SetLevel(lvl)
}

func MemGen(root string) []*File {
	initFs(false, true)
	l.Debugf("Initialized Source FS: %v", srcFs)
	l.Debugf("Initialized Destination FS: %v", destFs)
	return generate(root)
}

func RealGen(root string) []*File {
	initFs(false, false)
	return generate(root)
}

func GetFs() (afero.Fs, afero.Fs) {
	return srcFs, destFs
}

func ClearFs() {
	srcFs, destFs = nil, nil
}
