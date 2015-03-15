package file

import (
  l "github.com/Sirupsen/logrus"
  "os"
  "path/filepath"
  "text/template"
)

func ParseFiles(dir string, def string) (fq *FileQueue, err error) {
  l.WithField("directory", dir).Info("Parsing files")
  fq = newFileQueue()
  err = filepath.Walk(dir, walkfunc(def, fq))
  fq.populateQueue()
  return
}

func walkfunc(def string, fq *FileQueue) filepath.WalkFunc {
  return func(path string, info os.FileInfo, err error) error {
    ext := filepath.Ext(path)
    if info.Mode().IsRegular() && ext != ".skip" && ext != ".ignore" {
      return parseFile(path, def, fq)
    }
    l.WithField("path", path).Debug("Skipping")
    return nil
  }
}

func parseFile(path string, def string, fq *FileQueue) (err error) {
  l.WithField("path", path).Debug("Parsing file")

  f := newFile(path, def, filepath.Base(path))
  contents, err := f.Read()
  if err != nil {
    return
  }

  t, err := newTemplate(path, f).Parse(string(contents))
  if err != nil {
    return
  }
  f.Template(t)
  fq.add(f)
  return
}

func newTemplate(path string, f File) *template.Template {
  return template.New(path).Funcs(newFuncMap(f))
}
