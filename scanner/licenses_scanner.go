package scanner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ryanuber/go-license"
)

type LicensesScanner struct{}

func (l *LicensesScanner) GetLicenses(root string) (licenses map[string]*license.License, err error) {
	pkgScanner := NewPackagesScanner()

	pkgNames, err := pkgScanner.GetPackages(root)
	if err != nil {
		return
	}

	licenses = make(map[string]*license.License)

	for pkgName := range pkgNames {
		finalPkgName, lic, err := l.getPkgLicense(root, pkgName)
		if err != nil || lic == nil {
			continue
		}
		licenses[finalPkgName] = lic
	}

	return
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
