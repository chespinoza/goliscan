package scanner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
)

var defaultOutputTemplate = "[{{.Type | printf \"%8s\"}}]  {{.Message | printf \"%-28s\"}}  license = {{.License | printf \"%-14s\"}}  direct dependency = {{.Direct | printf \"%-3s\"}}  package = {{.Package}}"

type outputLine struct {
	Type    string
	Message string
	License string
	Direct  string
	Package string
}

type OutputPrinter struct {
	Lines []*outputLine

	tmpl *template.Template
}

func (l *OutputPrinter) AddData(level, message, pkgName, license string, direct bool) {
	directString := "no"
	if direct {
		directString = "yes"
	}

	l.Lines = append(l.Lines, &outputLine{
		Type:    level,
		Message: message,
		License: license,
		Direct:  directString,
		Package: pkgName,
	})
}

func (l *OutputPrinter) Print() (err error) {
	for _, line := range l.Lines {
		result := &bytes.Buffer{}
		err = l.tmpl.Execute(result, line)
		if err != nil {
			return
		}

		fmt.Println(strings.TrimSpace(result.String()))
	}

	return
}

func (l *OutputPrinter) PrintJSON() (err error) {
	jsonBytes, err := json.Marshal(l.Lines)
	if err != nil {
		return err
	}

	fmt.Println(string(jsonBytes))

	return
}

func NewOutputPrinter(outputTemplate string) (*OutputPrinter, error) {
	if outputTemplate == "" {
		outputTemplate = defaultOutputTemplate
	}

	tmpl, err := template.New("output-template").Parse(outputTemplate)
	if err != nil {
		return nil, err
	}

	outputPrinter := &OutputPrinter{
		tmpl: tmpl,
	}

	return outputPrinter, nil
}
