package main

import (
	"os"
	"path/filepath"
	"runtime"

	l "github.com/Sirupsen/logrus"
	"github.com/albertrdixon/tmplnator"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	debug       = kingpin.Flag("debug", "Enable debug mode").Short('d').Bool()
	quiet       = kingpin.Flag("quiet", "Enable quiet mode").Short('q').Bool()
	printTmpDir = kingpin.Flag("print-tmp", "Print out TmpDir").Short('p').Bool()
	forceTmp    = kingpin.Flag("force-tmp", "Force all generated files to be written to TmpDir").Short('F').Bool()
	tmpDir      = kingpin.Flag("tmpdir", "Set TmpDir").Default(filepath.Join(os.TempDir(), "T2")).PlaceHolder("\"$TMPDIR/T2\"").OverrideDefaultFromEnvar("T2_TMPDIR").Short('T').String()
	rootDir     = kingpin.Arg("template-directory", "Directory under which there are templates").Required().ExistingDir()
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	kingpin.Version(tmplnator.Version)
	kingpin.Parse()

	l.SetOutput(os.Stdout)
	tmplnator.LogLevel("info")
	if *quiet {
		tmplnator.LogLevel("fatal")
	} else if *debug {
		tmplnator.LogLevel("debug")
	}

	if *forceTmp {
		tmplnator.ForceTemp = true
	}
	if *tmpDir != "" {
		tmplnator.TmpDir = *tmpDir
	}

	tmplnator.RealGen(*rootDir)
	if *printTmpDir {
		println(tmplnator.TmpDir)
	}
}
