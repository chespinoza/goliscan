package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/ryanuber/go-license"

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

	c.checker.OkStateHandler = func(pkgName string, license *license.License) {
		printer.AddData("OK", "Found accepted license", pkgName, license.Type)
	}

	c.checker.ExceptionedStateHandler = func(pkgName string, license *license.License) {
		printer.AddData("WARNING", "Found exceptioned package", pkgName, license.Type)
	}

	c.checker.CriteriaUnknownStateHandler = func(pkgName string, license *license.License) {
		printer.AddData("WARNING", "Criteria for license unknown", pkgName, license.Type)
	}

	c.checker.CriticalStateHandler = func(pkgName string, license *license.License) {
		printer.AddData("CRITICAL", "Found unaccepted license", pkgName, license.Type)
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
