package logger

import (
  "fmt"
)

var (
  Level int
  Fmt   string
)

func Quiet(msg string, args ...interface{}) {
  if Level == 0 {
    m := fmt.Sprintf(msg, args...)
    fmt.Printf("%s\n", m)
  }
}

func Info(msg string, args ...interface{}) {
  if Level == 1 {
    m := fmt.Sprintf(msg, args...)
    fmt.Printf(Fmt, m)
  }
}

func Error(msg string, args ...interface{}) {
  m := fmt.Sprintf(msg, args...)
  fmt.Print(fmt.Errorf(Fmt, m))
}

func Debug(msg string, args ...interface{}) {
  if Level > 1 {
    m := fmt.Sprintf(msg, args...)
    fmt.Printf("%s\n", m)
  }
}

func init() {
  Level = 1
  Fmt = "==> %s\n"
}
