package scanner

import (
	"os"
	"sort"

	"github.com/ryanuber/go-license"

	"gitlab.com/tmaczukin/goliscan/config"
)

type licensesOutputHandlerFn func(pkgName string, license *license.License)

type LicensesOutputHandler struct {
	Printer *OutputPrinter

	handler  licensesOutputHandlerFn
	licenses map[string]*license.License
	settings config.OutputSettings
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
	h.licenses, err = licenseScanner.GetLicenses(root)

	return
}

func (h *LicensesOutputHandler) handleSortedList() (err error) {
	var keys []string
	for pkgName := range h.licenses {
		keys = append(keys, pkgName)
	}
	sort.Strings(keys)

	for _, pkgName := range keys {
		license := h.licenses[pkgName]
		h.handler(pkgName, license)
	}

	return
}

func NewLicensesOutputHandler(settings config.OutputSettings, handler licensesOutputHandlerFn) *LicensesOutputHandler {
	return &LicensesOutputHandler{
		handler:  handler,
		settings: settings,
	}
}
