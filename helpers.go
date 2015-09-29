package tmplnator

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	"code.google.com/p/go-uuid/uuid"
)

func newFuncMap(f *FileInfo) map[string]interface{} {
	return map[string]interface{}{
		"file":        f.SetFilename,
		"path":        f.SetFullpath,
		"mode":        f.SetMode,
		"dir_mode":    f.SetDirmode,
		"source":      f.Source,
		"timestamp":   timestamp,
		"to_json":     toJSON,
		"from_json":   fromJSON,
		"get":         get,
		"eql":         reflect.DeepEqual,
		"exists":      exists,
		"-e":          exists,
		"has_key":     hasKey,
		"def":         def,
		"default":     def,
		"url":         parseURL,
		"sha":         hash,
		"downcase":    downcase,
		"upcase":      upcase,
		"titleize":    titleize,
		"strip":       trimSpace,
		"split":       split,
		"replace":     replace,
		"trim":        trim,
		"trim_suffix": trimSuffix,
		"fields":      fields,
		"has_suffix":  hasSuffix,
		"contains":    contains,
		"join":        join,
		"uuid":        uuid.New,
	}
}

func parseURL(rawurl string) (u *url.URL, e error) { u, e = url.Parse(rawurl); return }
func timestamp() string                            { return time.Now().String() }

func hasSuffix(a, b interface{}) bool         { return ssb(strings.HasSuffix, a, b) }
func contains(a, b interface{}) bool          { return ssb(strings.Contains, a, b) }
func trim(a, b interface{}) interface{}       { return ssm(strings.Trim, a, b) }
func trimSuffix(a, b interface{}) interface{} { return ssm(strings.TrimSuffix, a, b) }
func downcase(s interface{}) interface{}      { return sm(strings.ToLower, s) }
func upcase(s interface{}) interface{}        { return sm(strings.ToUpper, s) }
func titleize(s interface{}) interface{}      { return sm(strings.Title, s) }
func trimSpace(s interface{}) interface{}     { return sm(strings.TrimSpace, s) }

func fromJSON(d interface{}) (j map[string]interface{}, e error) {
	switch d := d.(type) {
	case string:
		e = json.Unmarshal([]byte(d), &j)
	case []byte:
		e = json.Unmarshal(d, &j)
	default:
		return nil, errors.New("from_json only accepts a string or []byte")
	}
	return
}

func toJSON(i interface{}) (string, error) {
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(i); err != nil {
		return "", err
	}
	return strings.TrimSpace(b.String()), nil
}

func concat(a interface{}, bs ...interface{}) string {
	switch a := a.(type) {
	default:
		return ""
	case string:
		b := new(bytes.Buffer)
		b.WriteString(a)
		for _, s := range bs {
			if s, ok := s.(string); ok {
				b.WriteString(s)
			}
		}
		return b.String()
	}
}

func exists(p interface{}) (bool, error) {
	if p, ok := p.(string); ok {
		if _, e := os.Stat(p); e != nil {
			if os.IsNotExist(e) {
				return false, nil
			}
			return false, e
		}
		return true, nil
	}
	return false, nil
}

func hasKey(m, k interface{}) bool {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	ma, ke := reflect.ValueOf(m), reflect.ValueOf(k)
	if !ma.IsValid() || ma.IsNil() || ma.Kind() != reflect.Map {
		return false
	}
	if v := ma.MapIndex(ke); !v.IsValid() {
		return false
	}
	return true
}

func def(a, b interface{}) interface{} {
	switch a := a.(type) {
	case string:
		if a == "" {
			return b
		}
	default:
		if a == nil {
			return b
		}
	}
	return a
}

func get(s, i interface{}) interface{} {
	a, b := reflect.ValueOf(s), reflect.ValueOf(i)
	switch a.Kind() {
	default:
		return ""
	case reflect.Array:
		return a.Index(int(b.Int())).Interface()
	case reflect.String:
		return a.Index(int(b.Int())).Interface()
	case reflect.Slice:
		return a.Index(int(b.Int())).Interface()
	case reflect.Map:
		return a.MapIndex(b).Interface()
	}
}

func hash(d interface{}) (string, error) {
	if d, ok := d.(string); ok {
		h := sha1.New()
		io.WriteString(h, d)
		return fmt.Sprintf("%x", h.Sum(nil)), nil
	}
	return "", errors.New("sha only accepts a string")
}

func fields(s interface{}) interface{} {
	if s == nil {
		return ""
	}
	switch s := s.(type) {
	default:
		return s
	case string:
		return strings.Fields(s)
	}
}

func split(a, b interface{}) interface{} {
	if a == nil {
		return ""
	}
	switch a := a.(type) {
	default:
		return a
	case string:
		if b, ok := b.(string); ok {
			return strings.Split(a, b)
		}
		return a
	}
}

func join(a, b interface{}) (interface{}, error) {
	switch a := a.(type) {
	default:
		return "", errors.New("strings.Join requires a string slice")
	case []string:
		if b, ok := b.(string); ok {
			return strings.Join(a, b), nil
		}
		return "", errors.New("The separator for strings.Join must be a string")
	}
}

func replace(s interface{}, a, b string, n int) interface{} {
	switch s := s.(type) {
	default:
		return s
	case string:
		return strings.Replace(s, a, b, n)
	}
}

func sm(fn func(string) string, s interface{}) interface{} {
	if s == nil {
		return ""
	}
	switch s := s.(type) {
	default:
		return s
	case string:
		return fn(s)
	}
}

func ssm(fn func(string, string) string, a, b interface{}) interface{} {
	if a == nil {
		return ""
	}
	switch a := a.(type) {
	default:
		return a
	case string:
		if b, ok := b.(string); ok {
			return fn(a, b)
		}
		return a
	}
}

func ssb(fn func(string, string) bool, a, b interface{}) bool {
	if a == nil || b == nil {
		return false
	}
	switch a := a.(type) {
	default:
		return false
	case string:
		if b, ok := b.(string); ok {
			return fn(a, b)
		}
		return false
	}
}
