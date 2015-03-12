package backend

import (
  "fmt"
)

type Error struct {
  Code int
  Msg  string
}

func (e Error) Error() string {
  return fmt.Sprintf("[%d] %s", e.Code, e.Msg)
}

func KeyNotFound(key string) Error {
  return Error{
    Code: 1,
    Msg:  fmt.Sprintf("Key not found: %s", key),
  }
}
