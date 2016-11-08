package commands

import (
	"fmt"

	"github.com/codegangsta/cli"

	"gitlab.com/tmaczukin/goliscan/config"
	"gitlab.com/tmaczukin/goliscan/scanner"
)

type CheckCommand struct {
	config.OutputSettings

	File string `short:"c" long:"config" description:"Configuration file"`

	foundNonApproved bool
	checker          *scanner.LicenseChecker
}

func (c *CheckCommand) Execute(context *cli.Context) {
	checker, err := scanner.NewLicenseChecker(config.NewConfig(c.File))
	if err != nil {
		ThrowError(err)
	}

	c.checker = checker

	outputHandler := c.prepareHandlers()

	err = outputHandler.HandleLicensesOutput()
	if err != nil {
		ThrowError(err)
	}

	if c.foundNonApproved {
		err = fmt.Errorf("At least one unaccepted license was found!")
	}

	if err != nil {
		ThrowError(err)
	}
}

func (c *CheckCommand) prepareHandlers() (outputHandler *scanner.LicensesOutputHandler) {
	outputHandler = scanner.NewLicensesOutputHandler(c.OutputSettings, c.checker.Check)

	printer, err := outputHandler.GetPrinter()
	if err != nil {
		ThrowError(err)
	}

	c.checker.OkStateHandler = func(pkgName string, licenseSearchResult scanner.LicenseSearchResult) {
		printer.AddData("OK", "Found accepted license", pkgName, licenseSearchResult.License.Type, licenseSearchResult.Direct)
	}

	c.checker.ExceptionedStateHandler = func(pkgName string, licenseSearchResult scanner.LicenseSearchResult) {
		printer.AddData("WARNING", "Found exceptioned package", pkgName, licenseSearchResult.License.Type, licenseSearchResult.Direct)
	}

	c.checker.CriteriaUnknownStateHandler = func(pkgName string, licenseSearchResult scanner.LicenseSearchResult) {
		license := licenseSearchResult.License
		error := licenseSearchResult.Error

		if license != nil {
			printer.AddData("WARNING", "Criteria for license unknown", pkgName, license.Type, licenseSearchResult.Direct)
		} else if error != nil {
			printer.AddData("WARNING", error.Error(), pkgName, "unknown", licenseSearchResult.Direct)
		}
	}

	c.checker.CriticalStateHandler = func(pkgName string, licenseSearchResult scanner.LicenseSearchResult) {
		printer.AddData("CRITICAL", "Found unaccepted license", pkgName, licenseSearchResult.License.Type, licenseSearchResult.Direct)
		c.foundNonApproved = true
	}

	return
}

func init() {
	command := &CheckCommand{
		foundNonApproved: false,
	}

	RegisterCommand("check", "Scan vendored packages and check all found licenses", command)
}
