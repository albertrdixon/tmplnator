package template

import (
  "bytes"
  "crypto/sha1"
  "encoding/json"
  "fmt"
  "io"
  "net/url"
  "os"
  "reflect"
  "strings"
)

func newFuncMap() map[string]interface{} {
  return map[string]interface{}{
    "destination":     destination,
    "mode":            mode,
    "user":            user,
    "group":           group,
    "to_json":         marshalJSON,
    "from_json":       UnmarshalJSON,
    "from_json_array": UnmarshalJSONArray,
    "first":           arrayFirst,
    "last":            arrayLast,
    "file_exists":     fileExists,
    "has_key":         hasKey,
    "default":         defaultValue,
    "parseURL":        parseURL,
    "split":           strings.Split,
    "join":            strings.Join,
    "has_suffix":      strings.HasSuffix,
    "contains":        strings.Contains,
    "fields":          strings.Fields,
    "downcase":        strings.ToLower,
    "upcase":          strings.ToUpper,
    "trim_suffix":     strings.TrimSuffix,
    "trim_space":      strings.TrimSpace,
  }
}

func destination(t Template, d string, p ...interface{}) string {
  files[t.Sha1].Path = fmt.Sprintf(d, p...)
  return ""
}

func mode(t Template, m os.FileMode) string {
  files[t.Sha1].Mode = m
  return ""
}

func user(t Template, u int) string {
  files[t.Sha1].User = u
  return ""
}

func group(t Template, g int) string {
  files[t.Sha1].Group = g
  return ""
}

func UnmarshalJSON(data string) (map[string]interface{}, error) {
  var ret map[string]interface{}
  err := json.Unmarshal([]byte(data), &ret)
  return ret, err
}

func UnmarshalJSONArray(data string) ([]interface{}, error) {
  var ret []interface{}
  err := json.Unmarshal([]byte(data), &ret)
  return ret, err
}

func fileExists(path string) (bool, error) {
  _, err := os.Stat(path)
  if err == nil {
    return true, nil
  }
  if os.IsNotExist(err) {
    return false, nil
  }
  return false, err
}

func hasKey(item map[string]string, key string) bool {
  if _, ok := item[key]; ok {
    return true
  }
  return false
}

func defaultValue(arg interface{}, def interface{}) interface{} {
  if arg == nil {
    return def
  }
  if as, ok := arg.(string); ok {
    if as == "" {
      return def
    }
  }
  return arg
}

// arrayLast returns last item in the array
func arrayLast(input interface{}) interface{} {
  arr := reflect.ValueOf(input)
  return arr.Index(arr.Len() - 1).Interface()
}

// arrayFirst returns first item in the array or nil if the
// input is nil or empty
func arrayFirst(input interface{}) interface{} {
  if input == nil {
    return nil
  }

  arr := reflect.ValueOf(input)

  if arr.Len() == 0 {
    return nil
  }

  return arr.Index(0).Interface()
}

func marshalJSON(input interface{}) (string, error) {
  var buf bytes.Buffer
  enc := json.NewEncoder(&buf)
  if err := enc.Encode(input); err != nil {
    return "", err
  }
  return strings.TrimSuffix(buf.String(), "\n"), nil
}

func hashSha1(input string) string {
  h := sha1.New()
  io.WriteString(h, input)
  return fmt.Sprintf("%x", h.Sum(nil))
}

func parseURL(rawurl string) (*url.URL, error) {
  u, err := url.Parse(rawurl)
  if err != nil {
    return nil, err
  }
  return u, nil
}
