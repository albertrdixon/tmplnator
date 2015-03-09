package template

import (
	"os"
	"path/filepath"
)

type File struct {
	Prefix string
	Path   string
	Name   string
	Mode   uint32
	User   int
	Group  int
}

func (f *File) Destination() string {
	return filepath.Join(f.Prefix, f.Path, f.Name)
}

func newFile(pre string, name string) *File {
	if pre == "" {
		pre = "/"
	}
	return &File{
		Prefix: pre,
		Path:   "",
		Name:   name,
		Mode:   0644,
		User:   os.Getuid(),
		Group:  os.Getegid(),
	}
}
