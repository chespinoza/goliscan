package commands

import (
	"github.com/codegangsta/cli"
	"github.com/ryanuber/go-license"

	"gitlab.com/tmaczukin/goliscan/config"
	"gitlab.com/tmaczukin/goliscan/scanner"
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

func (l *ListCommand) listHandler(pkgName string, license *license.License) {
	l.printer.AddData("INFO", "Found license", pkgName, license.Type)
}

func init() {
	RegisterCommand("list", "Scan vendored packages and list all found licenses", &ListCommand{})
}
