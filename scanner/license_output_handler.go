package scanner

import (
	"errors"
	"os"
	"sort"

	"github.com/ryanuber/go-license"

	"gitlab.com/tmaczukin/goliscan/config"
)

type licensesOutputHandlerFn func(pkgName string, licenseSearchResult LicenseSearchResult)

type LicenseSearchResult struct {
	License *license.License
	Error   error
	Direct  bool
}

func newLicenseSearchResult(lic *license.License, pkgInfo *PackageInfo, err error) LicenseSearchResult {
	if err != nil {
		var error string
		switch {
		case err.Error() == license.ErrUnrecognizedLicense:
			error = "Could not guess license type"
		case err.Error() == license.ErrNoLicenseFile:
			error = "Unable to find any license file"
		case err.Error() == license.ErrMultipleLicenses:
			error = "Multiple license files found"
		}

		return LicenseSearchResult{
			License: nil,
			Error:   errors.New(error),
			Direct:  !pkgInfo.FromVendor,
		}
	}

	return LicenseSearchResult{
		License: lic,
		Error:   nil,
		Direct:  !pkgInfo.FromVendor,
	}
}

type LicensesOutputHandler struct {
	Printer *OutputPrinter

	handler  licensesOutputHandlerFn
	settings config.OutputSettings

	licenseSearchResults map[string]LicenseSearchResult
}

func (h *LicensesOutputHandler) GetPrinter() (*OutputPrinter, error) {
	if h.Printer == nil {
		printer, err := NewOutputPrinter(h.settings.OutputTemplate)
		if err != nil {
			return nil, err
		}

		h.Printer = printer
	}

	return h.Printer, nil

}

func (h *LicensesOutputHandler) HandleLicensesOutput() (err error) {
	err = h.searchLicenses()
	if err != nil {
		return
	}

	err = h.handleSortedList()
	if err != nil {
		return
	}

	printer, err := h.GetPrinter()
	if err != nil {
		return
	}

	if h.settings.UseJSON {
		err = printer.PrintJSON()
	} else {
		err = printer.Print()
	}

	return
}

func (h *LicensesOutputHandler) searchLicenses() (err error) {
	root, err := os.Getwd()
	if err != nil {
		return
	}

	licenseScanner := NewLicenseScanner()
	h.licenseSearchResults, err = licenseScanner.GetLicenses(root)

	return
}

func (h *LicensesOutputHandler) handleSortedList() (err error) {
	var keys []string
	for pkgName := range h.licenseSearchResults {
		keys = append(keys, pkgName)
	}
	sort.Strings(keys)

	for _, pkgName := range keys {
		licenseSearchResult := h.licenseSearchResults[pkgName]

		if h.shouldSkip(licenseSearchResult.Direct) {
			continue
		}

		h.handler(pkgName, licenseSearchResult)
	}

	return
}

func (h *LicensesOutputHandler) shouldSkip(isDependencyDirect bool) bool {
	return (h.settings.DirectOnly && !isDependencyDirect) ||
		(h.settings.IndirectOnly && isDependencyDirect)
}

func NewLicensesOutputHandler(settings config.OutputSettings, handler licensesOutputHandlerFn) *LicensesOutputHandler {
	return &LicensesOutputHandler{
		handler:  handler,
		settings: settings,
	}
}
