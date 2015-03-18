package file

import (
  "os"
  "path/filepath"
  "testing"
)

func TestFileExists(t *testing.T) {
  tmpFile := filepath.Join(os.TempDir(), "sg-test")
  os.Create(tmpFile)
  defer os.Remove(tmpFile)
  fakeFile := "path/to/fake_file"

  // Test fileExists on a real file
  shouldBeTrue, err := fileExists(tmpFile)
  if err != nil {
    t.Errorf("fileExists returned an error with %q: %v", tmpFile, err)
  }
  if !shouldBeTrue {
    t.Errorf("Could not find %q even though it was created", tmpFile)
  }

  // Test fileExists on a fake file
  shouldBeFalse, err := fileExists(fakeFile)
  if err != nil {
    t.Errorf("fileExists returned an error with %q: %v", fakeFile, err)
  }
  if shouldBeFalse {
    t.Error("Fake file returned true, but hsould be false")
  }
}

func TestParseURL(t *testing.T) {
  goodURL := "http://my.domain.net:8888"
  expectedHost := "my.domain.net:8888"

  good, err := parseURL(goodURL)
  if err != nil {
    t.Errorf("Should not have gotten an error parsing %q: %v", goodURL, err)
  }
  if good == nil {
    t.Errorf("Should have parsed %q correctly, but got nil", good)
  }
  if good.Host != expectedHost {
    t.Errorf("Expected Host to be %q but got %q", expectedHost, good.Host)
  }
}

func TestDefaultValue(t *testing.T) {
  var defaulttests = []struct {
    in       interface{}
    def      interface{}
    expected interface{}
  }{
    {nil, "default", "default"},
    {nil, 2, 2},
    {"", "default", "default"},
    {"", 3, 3},
    {"realoutput", "shouldnotshow", "realoutput"},
  }

  for _, dt := range defaulttests {
    out := defaultValue(dt.in, dt.def)
    if out != dt.expected {
      t.Errorf("defaultValue(%v, %v): %q, want %q", dt.in, dt.def, out, dt.expected)
    }
  }
}

func TestJSON(t *testing.T) {
  goodJSON := "{\"Key1\":\"Value1\",\"Key2\":\"Value2\"}"
  badJSON := "{\"Key3\":\"Blah\""

  out, err := UnmarshalJSON(goodJSON)
  if err != nil {
    t.Errorf("UnmarshalJSON(good): expected no error, got %v", err)
  }
  if out["Key2"] != "Value2" {
    t.Error("UnmarshalJSON(good): didn't get expected structure back")
  }

  out, err = UnmarshalJSON(badJSON)
  if err == nil {
    t.Errorf("UnmarshalJSON(bad): expected error, got %v", err)
  }
  if out != nil {
    t.Error("UnmarshalJSON(bad): expected out to be nil")
  }
}
