package config

import (
  "testing"
)

func TestRuntimeVersion(t *testing.T) {
  var vers string

  // dev build and sha
  vers = RuntimeVersion("0.0-dev", "sha")
  if vers != "0.0-dev-sha" {
    t.Errorf("Dev version with sha incorrect: %s\n", vers)
  }

  // dev build and no sha
  vers = RuntimeVersion("0.0-dev", "")
  if vers != "0.0-dev" {
    t.Errorf("Dev version with no sha incorrect: %s\n", vers)
  }

  // release build and sha
  vers = RuntimeVersion("0.0", "sha")
  if vers != "0.0" {
    t.Errorf("Release version with sha incorrect: %s\n", vers)
  }

  // release build and no sha
  vers = RuntimeVersion("0.0", "")
  if vers != "0.0" {
    t.Errorf("Release version with no sha incorrect: %s\n", vers)
  }
}
