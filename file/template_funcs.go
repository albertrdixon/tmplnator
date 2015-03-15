package file

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
  "time"
)

func newFuncMap(f File) map[string]interface{} {
  return map[string]interface{}{
    "dir":             f.setDir,
    "name":            f.setName,
    "mode":            f.setMode,
    "dir_mode":        f.setDirMode,
    "user":            f.setUser,
    "group":           f.setGroup,
    "skip":            f.setSkip,
    "env":             os.Getenv,
    "source":          f.Src,
    "timestamp":       timestamp,
    "to_json":         marshalJSON,
    "from_json":       UnmarshalJSON,
    "from_json_array": UnmarshalJSONArray,
    "fmt":             fmt.Sprintf,
    "first":           arrayFirst,
    "last":            arrayLast,
    "file_exists":     fileExists,
    "has_key":         hasKey,
    "in_env":          inEnv,
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

func UnmarshalJSON(data string) (j map[string]interface{}, e error) {
  e = json.Unmarshal([]byte(data), &j)
  return
}

func UnmarshalJSONArray(data string) (ja []interface{}, e error) {
  e = json.Unmarshal([]byte(data), &ja)
  return
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

func hasKey(item map[string]string, key string) (ok bool) {
  _, ok = item[key]
  return
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

func inEnv(key string) bool {
  if ok := os.Getenv(key); ok == "" {
    return false
  }
  return true
}

func timestamp() string {
  return time.Now().String()
}
