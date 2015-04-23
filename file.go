package tmplnator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	l "github.com/Sirupsen/logrus"
	"github.com/spf13/afero"
)

var extensions = []string{".tpl", ".tmpl", ".templ", ".template"}

// Reader is a convenience interface
type Reader interface {
	Read(p []byte) (n int, err error)
}

// File type wraps text/tsemplate so we can update the info before generatingthe file
type File struct {
	template *template.Template
	info     *FileInfo
	errs     []error
}

type FileInfo struct {
	template string
	fileName string
	fullPath string
	mode     os.FileMode
	dirmode  os.FileMode
	modtime  time.Time
	size     int64
}

func isTemplate(filename string) (answer bool) {
	for _, ext := range extensions {
		if strings.HasSuffix(filename, ext) {
			answer = true
			break
		}
	}
	return
}

func NewFile(path string) *File {
	var filename string
	var fullpath string

	filename = filepath.Base(path)
	for _, suffix := range extensions {
		filename = strings.TrimSuffix(filename, suffix)
	}

	if TmpDir == "" {
		fullpath = filepath.Join(filepath.Dir(path), filename)
	} else {
		fullpath = filepath.Join(TmpDir, filename)
	}

	f := &File{
		template: template.New(path),
		info: &FileInfo{
			fileName: filename,
			fullPath: fullpath,
			mode:     os.FileMode(0644),
			dirmode:  os.FileMode(0755),
			modtime:  time.Time{},
			size:     0,
		},
	}

	return f.RegisterFuncs(newFuncMap(f.Info()))
}

func (f *FileInfo) Name() string       { return f.fileName }
func (f *FileInfo) Mode() os.FileMode  { return f.mode }
func (f *FileInfo) ModTime() time.Time { return f.modtime }
func (f *FileInfo) IsDir() bool        { return false }
func (f *FileInfo) Sys() interface{}   { return nil }
func (f *FileInfo) Size() int64        { return f.size }
func (f *FileInfo) String() string     { return f.Name() }
func (f *FileInfo) Source() string     { return f.template }

func (f *File) TemplateName() string { return f.template.Name() }
func (f *File) Name() string         { return f.info.Name() }
func (f *File) FullPath() string     { return f.info.fullPath }
func (f *File) Dir() string          { return filepath.Dir(f.FullPath()) }
func (f *File) Info() *FileInfo      { return f.info }
func (f *File) String() string {
	return fmt.Sprintf("(template)%s (file)%s", f.TemplateName(), f.FullPath())
}
func (f *File) HasErrs() bool {
	if len(f.errs) > 0 {
		return true
	}
	return false
}

func (f *File) RegisterFuncs(funcMap map[string]interface{}) *File {
	f.template = f.template.Funcs(funcMap)
	return f
}

func (f *FileInfo) SetFilename(format string, items ...interface{}) string {
	path := fmt.Sprintf(format, items...)
	l.Debugf("%v: Changing name to %s", f, filepath.Base(path))
	f.fileName = filepath.Base(path)
	return f.SetFullpath(filepath.Dir(path))
}

func (f *FileInfo) SetFullpath(format string, items ...interface{}) string {
	path := fmt.Sprintf(format, items...)
	l.Debugf("%v: Changing path to %s", f, path)
	fullpath := filepath.Join(path, f.Name())
	if ForceTemp {
		fullpath = filepath.Join(TmpDir, f.Name())
	}
	f.fullPath = filepath.Clean(fullpath)
	return ""
}

func (f *FileInfo) SetMode(mode os.FileMode) string {
	l.Debugf("%v: Changing filemode to %v", f, mode)
	f.mode = mode
	return ""
}

func (f *FileInfo) SetDirmode(mode os.FileMode) string {
	l.Debugf("%v: Changing dirmode to %v", f, mode)
	f.dirmode = mode
	return ""
}

func (f *File) Open() (afero.File, error) {
	return destFs.Open(f.FullPath())
}

func (f *File) Read(p []byte) (n int, err error) {
	fh, err := f.Open()
	if err != nil {
		return
	}

	n, err = fh.Read(p)
	return
}

func ParseTemplate(f *File, fs afero.Fs) error {
	defer func() {
		if r := recover(); r != nil {
			f.errs = append(f.errs, NewFileError(f, "panic: %v", r))
		}
	}()

	fh, err := fs.Open(f.template.Name())
	defer fh.Close()

	if err != nil {
		f.errs = append(f.errs, err)
		return err
	}

	var b bytes.Buffer
	if _, err := b.ReadFrom(fh); err != nil {
		e := NewFileError(f, err.Error())
		f.errs = append(f.errs, e)
		return e
	}

	t, err := f.template.Parse(b.String())
	if err != nil {
		e := NewFileError(f, err.Error())
		f.errs = append(f.errs, e)
		return e
	}

	f.template = t
	return nil
}

func WriteFile(f *File, fs afero.Fs) error {
	b := new(bytes.Buffer)
	if err := f.template.Execute(b, data); err != nil {
		e := NewFileError(f, err.Error())
		f.errs = append(f.errs, e)
		return e
	}

	if _, err := fs.Stat(f.Dir()); err != nil {
		if err = fs.MkdirAll(f.Dir(), f.info.dirmode); err != nil {
			e := NewFileError(f, err.Error())
			f.errs = append(f.errs, e)
			return e
		}
	}

	gf, err := fs.Create(f.FullPath())
	defer gf.Close()

	if err != nil {
		e := NewFileError(f, err.Error())
		f.errs = append(f.errs, e)
		return e
	}

	if l.GetLevel() == l.DebugLevel {
		l.Debugf("%q CONTENT: %s", f.Name(), b.String())
	}
	n, err := b.WriteTo(gf)
	if err != nil {
		e := NewFileError(f, err.Error())
		f.errs = append(f.errs, e)
		return e
	}
	f.info.size = n
	fs.Chmod(f.FullPath(), f.info.Mode())
	return nil
}
