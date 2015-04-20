package tmplnator

import "testing"

var testDir = "fixtures/test/"

func TestGenerate(t *testing.T) {
	InitFs(false, true)
	for _, file := range Generate(testDir) {
		if file.HasErrs() {
			t.Errorf("%q: Got errors!", file.TemplateName())
			for _, err := range file.errs {
				t.Errorf("%v", err)
			}
		}
	}
}
