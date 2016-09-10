package scanner

import (
	"github.com/ryanuber/go-license"

	"gitlab.com/tmaczukin/goliscan/config"
)

type StateHandlerFn func(pkgName string, license *license.License)

type LicenseChecker struct {
	OkStateHandler              StateHandlerFn
	ExceptionedStateHandler     StateHandlerFn
	CriteriaUnknownStateHandler StateHandlerFn
	CriticalStateHandler        StateHandlerFn

	configuration *config.Config
}

func (c *LicenseChecker) Check(pkgName string, license *license.License) {
	var handler StateHandlerFn

	switch {
	case c.isExceptioned(license, pkgName):
		handler = c.ExceptionedStateHandler
	case c.isCritical(license, pkgName):
		handler = c.CriticalStateHandler
	case c.isOK(license):
		handler = c.OkStateHandler
	default:
		handler = c.CriteriaUnknownStateHandler
	}

	handler(pkgName, license)
}

func (c *LicenseChecker) isExceptioned(license *license.License, pkgName string) bool {
	return c.configuration.Unaccepted[license.Type] && c.configuration.Exceptions[pkgName]
}

func (c *LicenseChecker) isCritical(license *license.License, pkgName string) bool {
	return (c.configuration.Unaccepted[license.Type] && !c.configuration.Exceptions[pkgName]) ||
		(!c.configuration.Accepted[license.Type] && !c.configuration.AllowUnknown)
}

func (c *LicenseChecker) isOK(license *license.License) bool {
	return c.configuration.Accepted[license.Type]
}

func NewLicenseChecker(configuration *config.Config) (*LicenseChecker, error) {
	err := configuration.ReadConfig()
	if err != nil {
		return nil, err
	}

	checker := &LicenseChecker{
		configuration: configuration,
	}

	return checker, nil
}
