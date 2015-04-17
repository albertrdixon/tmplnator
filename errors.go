package tmplnator

import (
	"fmt"

	l "github.com/Sirupsen/logrus"
)

// TemplateError are template related errors
type TemplateError struct {
	template string
	msg      string
}

// GeneratorError are generator related errors
type GeneratorError struct {
	thread string
	msg    string
}

func (t TemplateError) Error() string {
	return fmt.Sprintf("%s: %s", t.template, t.msg)
}

func (g GeneratorError) Error() string {
	return fmt.Sprintf("%s: %s", g.thread, g.msg)
}

func newGeneratorError(thread string, format string, items ...interface{}) GeneratorError {
	ge := GeneratorError{thread, fmt.Sprintf(format, items...)}
	l.Error(ge.Error())
	return ge
}

func newTemplateError(template string, format string, items ...interface{}) TemplateError {
	te := TemplateError{template, fmt.Sprintf(format, items...)}
	l.Error(te.Error())
	return te
}
