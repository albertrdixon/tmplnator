package file

import (
  "bytes"
  "fmt"
  l "github.com/Sirupsen/logrus"
  "io/ioutil"
  "os"
  "path/filepath"
  tmpl "text/template"
)

type templateFile struct {
  template *tmpl.Template
  src      string
  name     string
  dir      string
  user     int
  group    int
  mode     os.FileMode
  dirmode  os.FileMode
  skip     bool
}

func (tf *templateFile) Write(b *bytes.Buffer, data interface{}) (err error) {
  l.WithFields(l.Fields{
    "template": tf.src,
    "data":     data,
  }).Debug("Executing template")
  err = tf.template.Execute(b, data)
  if err != nil {
    return err
  }

  if tf.skip {
    return nil
  }

  l.WithField("path", tf.dir).Debug("Creating directory")
  if _, err := os.Stat(tf.dir); err != nil {
    if err = os.MkdirAll(tf.dir, tf.dirmode); err != nil {
      return err
    }
  }

  l.WithField("path", tf.Destination()).Debug("Creating file")
  fh, err := os.OpenFile(tf.Destination(), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, tf.mode)
  if err != nil {
    return err
  }
  defer fh.Close()
  defer os.Chown(tf.Destination(), tf.user, tf.group)

  l.WithFields(l.Fields{
    "template": tf.src,
    "file":     tf.Destination(),
  }).Info("Generating file")
  _, err = fh.Write(b.Bytes())
  if err != nil {
    return err
  }
  return nil
}

func (tf *templateFile) Read() (b []byte, err error) {
  b, err = ioutil.ReadFile(tf.src)
  return
}

func (tf *templateFile) Src() string {
  return tf.src
}

func (tf *templateFile) Destination() string {
  return filepath.Join(tf.dir, tf.name)
}

func (tf *templateFile) DeleteTemplate() (err error) {
  err = os.Remove(tf.src)
  return
}

func (tf *templateFile) Template(t *tmpl.Template) {
  tf.template = t
}

func (tf *templateFile) setDir(dir string, args ...interface{}) string {
  for i, a := range args {
    if a == nil {
      args[i] = ""
    }
  }
  tf.dir = fmt.Sprintf(dir, args...)
  return ""
}

func (tf *templateFile) setName(name string, args ...interface{}) string {
  for i, a := range args {
    if a == nil {
      args[i] = ""
    }
  }
  tf.name = fmt.Sprintf(name, args...)
  return ""
}

func (tf *templateFile) setUser(uid int) string {
  tf.user = uid
  return ""
}

func (tf *templateFile) setGroup(gid int) string {
  tf.group = gid
  return ""
}

func (tf *templateFile) setMode(m os.FileMode) string {
  tf.mode = m
  return ""
}

func (tf *templateFile) setDirMode(dm os.FileMode) string {
  tf.dirmode = dm
  return ""
}

func (tf *templateFile) setSkip() string {
  tf.skip = true
  return ""
}

func newTemplateFile(path string, def string, name string) File {
  return &templateFile{
    src:     path,
    name:    name,
    dir:     def,
    mode:    os.FileMode(0644),
    dirmode: os.FileMode(0755),
    user:    os.Geteuid(),
    group:   os.Getegid(),
    skip:    false,
  }
}
