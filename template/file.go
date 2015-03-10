package template

import (
  "os"
  "path/filepath"
)

type File struct {
  Prefix string
  Path   string
  Name   string
  Mode   os.FileMode
  User   int
  Group  int
}

func (f *File) Destination() string {
  return filepath.Join(f.Prefix, f.Path, f.Name)
}

func (f *File) DestinationDir() string {
  return filepath.Join(f.Prefix, f.Path)
}

func GetDestination(t *Template) string {
  return files[t.Sha1].Destination()
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
