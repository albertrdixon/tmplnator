// Package logger provides super simple logging
// functionality for command line based programs.
package logger

import (
  "fmt"
)

var (
  // Level is the output level
  Level int
  // Fmt is the log format
  Fmt string
)

// Quiet prints msg if Level == 0
func Quiet(msg string, args ...interface{}) {
  if Level == 0 {
    m := fmt.Sprintf(msg, args...)
    fmt.Printf("%s\n", m)
  }
}

// Info prints msg if Level >= 1
func Info(msg string, args ...interface{}) {
  if Level >= 1 {
    m := fmt.Sprintf(msg, args...)
    fmt.Printf(Fmt, m)
  }
}

// Error always prints its message as an error
func Error(msg string, args ...interface{}) {
  m := fmt.Sprintf(msg, args...)
  fmt.Print(fmt.Errorf(Fmt, m))
}

// Debug prints msg if Level > 1
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
