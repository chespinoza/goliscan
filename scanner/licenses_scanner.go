package scanner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ryanuber/go-license"
)

type LicensesScanner struct{}

func (l *LicensesScanner) GetLicenses(root string) (licenses map[string]LicenseSearchResult, err error) {
	pkgScanner := NewPackagesScanner()

	pkgNames, err := pkgScanner.GetPackages(root)
	if err != nil {
		return
	}

	licenses = make(map[string]LicenseSearchResult)

	for pkgName, pkgInfo := range pkgNames {
		finalPkgName, lic, err := l.getPkgLicense(root, pkgName)
		if l.isKnownLicenseError(err) {
			continue
		}

		licenses[finalPkgName] = newLicenseSearchResult(lic, pkgInfo, err)
	}

	return
}

func (l *LicensesScanner) isKnownLicenseError(err error) bool {
	return err != nil &&
		err.Error() != license.ErrMultipleLicenses &&
		err.Error() != license.ErrNoLicenseFile &&
		err.Error() != license.ErrUnrecognizedLicense
}

func (l *LicensesScanner) getPkgLicense(root, pkgName string) (finalPkgName string, lic *license.License, err error) {
	finalPkgName = pkgName

	path := filepath.Join(root, "vendor", pkgName)
	pkgPath, err := os.Stat(path)
	if err != nil {
		return
	}

	if pkgPath.IsDir() {
		lic, err = license.NewFromDir(path)
		if err != nil {
			if err.Error() == license.ErrNoLicenseFile {
				return l.lookupPkgParentDir(root, pkgName)
			}
			return
		}
	}

	return
}

func (l *LicensesScanner) lookupPkgParentDir(root, pkgName string) (FinalPkgName string, lic *license.License, err error) {
	FinalPkgName = pkgName
	parts := strings.Split(pkgName, "/")
	partsCount := len(parts)

	if partsCount < 2 {
		return
	}

	FinalPkgName = strings.Join(parts[:partsCount-1], "/")
	return l.getPkgLicense(root, FinalPkgName)
}

func NewLicenseScanner() *LicensesScanner {
	return &LicensesScanner{}
}
