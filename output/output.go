package output

import (
	"fmt"
	"github.com/fatih/color"
)

type Output struct {
	indentation int
	testOnly    bool
}

func New() *Output {
	output := Output{
		indentation: 0,
		testOnly:    false,
	}
	return &output
}

func NewForTesting() *Output {
	output := Output{
		indentation: 0,
		testOnly:    true,
	}
	return &output
}

func (o *Output) Indent() *Output {
	o.indentation = o.indentation + 1
	return o
}

func (o *Output) Dedent() *Output {
	o.indentation = o.indentation - 1
	return o
}

func (o *Output) Error(err error) {
	o.Failure("%v", err)
}

func (o *Output) Failure(format string, a ...interface{}) *Output {
	if !o.testOnly {
		color.Red(o.indentFmt(format, a...))
	}
	return o
}

func (o *Output) Warning(format string, a ...interface{}) *Output {
	if !o.testOnly {
		color.Yellow(o.indentFmt(format, a...))
	}
	return o
}

func (o *Output) Info(format string, a ...interface{}) *Output {
	if !o.testOnly {
		color.Cyan(o.indentFmt(format, a...))
	}
	return o
}

func (o *Output) Done() *Output {
	return o.Success("Done.")
}

func (o *Output) Success(format string, a ...interface{}) *Output {
	if !o.testOnly {
		color.Green(o.indentFmt(format, a...))
	}
	return o
}

func (o *Output) indentFmt(format string, a ...interface{}) string {
	indent := ""
	for i := 0; i < o.indentation; i++ {
		indent = indent + "  "
	}
	if len(a) > 0 {
		return indent + fmt.Sprintf(format, a...)
	}
	return indent + format
}
