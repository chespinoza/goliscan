package commands

import (
	"github.com/urfave/cli"

	"github.com/chespinoza/goliscan/config"
	"github.com/chespinoza/goliscan/scanner"
)

type ListCommand struct {
	config.OutputSettings

	printer *scanner.OutputPrinter
}

func (l *ListCommand) Execute(context *cli.Context) {
	outputHandler := l.getOutputHandler()

	err := outputHandler.HandleLicensesOutput()
	if err != nil {
		ThrowError(err)
	}
}

func (l *ListCommand) getOutputHandler() (outputHandler *scanner.LicensesOutputHandler) {
	outputHandler = scanner.NewLicensesOutputHandler(l.OutputSettings, l.listHandler)

	printer, err := outputHandler.GetPrinter()
	if err != nil {
		ThrowError(err)
	}

	l.printer = printer

	return
}

func (l *ListCommand) listHandler(pkgName string, licenseSearchResult scanner.LicenseSearchResult) {
	license := licenseSearchResult.License
	error := licenseSearchResult.Error

	if license != nil {
		l.printer.AddData("INFO", "Found license", pkgName, license.Type, licenseSearchResult.Direct)
	} else if error != nil {
		l.printer.AddData("INFO", error.Error(), pkgName, "unknown", licenseSearchResult.Direct)
	}
}

func init() {
	RegisterCommand("list", "Scan vendored packages and list all found licenses", &ListCommand{})
}
