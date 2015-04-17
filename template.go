package tmplnator

import (
	"bytes"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	l "github.com/Sirupsen/logrus"
)

var extensions = []string{".tpl", ".tmpl", ".templ", ".template"}

// Reader is a convenience interface
type Reader interface {
	Read(p []byte) (n int, err error)
}

// Template type wraps text/template so we can update the info before generatingthe file
type Template struct {
	template *template.Template
	Name     string
	Info     *FileInfo
}

type FileInfo struct {
	skip    bool
	Source  string
	Name    string
	Dir     string
	Owner   int
	Group   int
	Mode    os.FileMode
	Dirmode os.FileMode
}

func NewTemplate(path string) *Template {
	name := filepath.Base(path)
	for _, suffix := range extensions {
		name = strings.TrimSuffix(name, suffix)
	}

	dir := TmpDir
	if dir == "" {
		dir = filepath.Dir(path)
	}
	return &Template{
		Name: filepath.Base(path),
		Info: &FileInfo{
			skip:    false,
			Source:  path,
			Name:    name,
			Dir:     dir,
			Owner:   os.Geteuid(),
			Group:   os.Getegid(),
			Mode:    os.FileMode(0644),
			Dirmode: os.FileMode(0755),
		},
	}
}

func (f *Template) Output() string {
	return filepath.Join(f.Info.Dir, f.Info.Name)
}

func (fi *FileInfo) SetName(format string, items ...interface{}) string {
	fi.Name = fmt.Sprintf(format, items...)
	return ""
}

func (fi *FileInfo) SetDir(format string, items ...interface{}) string {
	if ForceTemp {
		fi.Dir = TmpDir
	} else {
		fi.Dir = fmt.Sprintf(format, items...)
	}
	return ""
}

func (fi *FileInfo) SetOwner(owner string) string {
	u, err := user.Lookup(owner)
	if err != nil {
		errChan <- err
		fi.Owner = os.Geteuid()
		return ""
	}
	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		errChan <- err
		fi.Owner = os.Geteuid()
		return ""
	}
	fi.Owner = uid
	return ""
}

func (fi *FileInfo) SetUID(owner int) string {
	fi.Owner = owner
	return ""
}

func (fi *FileInfo) SetGID(group int) string {
	fi.Group = group
	return ""
}

func (fi *FileInfo) SetMode(mode os.FileMode) string {
	fi.Mode = mode
	return ""
}

func (fi *FileInfo) SetDirmode(mode os.FileMode) string {
	fi.Dirmode = mode
	return ""
}

func (fi *FileInfo) Skip() string {
	fi.skip = true
	return ""
}

func (fi *FileInfo) Src() string {
	return fi.Source
}

func (t *Template) Skip() bool {
	return t.Info.skip
}

func (t *Template) Parse(r Reader) error {
	defer func() {
		if r := recover(); r != nil {
			l.WithFields(l.Fields{
				"package": pkg,
			}).Errorf("Recover from panic: %s", r)
			errChan <- newTemplateError(t.Name, "%v", r)
		}
	}()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(r); err != nil {
		return TemplateError{t.Name, err.Error()}
	}

	tmpl := template.New(t.Name).Funcs(newFuncMap(t.Info))
	tmpl, err := tmpl.Parse(buf.String())
	if err != nil {
		return TemplateError{t.Name, err.Error()}
	}

	t.template = tmpl
	return nil
}
