package main

import (
	"os"
	"path/filepath"

	l "github.com/Sirupsen/logrus"
	"github.com/albertrdixon/tmplnator"
	"gopkg.in/alecthomas/kingpin.v1"
)

var (
	debug       = kingpin.Flag("debug", "Enable debug mode").Short('d').Bool()
	quiet       = kingpin.Flag("quiet", "Enable quiet mode").Short('q').Bool()
	printTmpDir = kingpin.Flag("print-tmp", "Print out TmpDir").Short('p').Bool()
	forceTmp    = kingpin.Flag("force-tmp", "Force all generated files to be written to TmpDir").Short('F').Bool()
	tmpDir      = kingpin.Flag("tmpdir", "Set TmpDir").Default(filepath.Join(os.TempDir(), "T2")).OverrideDefaultFromEnvar("T2_TMPDIR").Short('T').String()
	rootDir     = kingpin.Arg("template-directory", "Directory under which there are templates").Required().ExistingDir()
)

func main() {
	kingpin.Version(tmplnator.Version)
	kingpin.Parse()

	l.SetOutput(os.Stdout)
	l.SetLevel(l.InfoLevel)
	if *quiet {
		l.SetLevel(l.FatalLevel)
	} else if *debug {
		l.SetLevel(l.DebugLevel)
	}

	if *forceTmp {
		tmplnator.ForceTemp = true
	}
	if *tmpDir != "" {
		tmplnator.TmpDir = *tmpDir
	}

	tmplnator.InitFs(false, false)
	tmplnator.Generate(*rootDir)
	if *printTmpDir {
		println(tmplnator.TmpDir)
	}
}
