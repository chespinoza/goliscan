package scanner

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type PackageInfo struct {
	FromVendor bool
}

type PackagesScanner struct {
	pkgNames map[string]*PackageInfo
}

func (s *PackagesScanner) GetPackages(root string) (map[string]*PackageInfo, error) {
	root, err := s.resolveRootSymlink(root)
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(root, s.scanImports)
	if err != nil {
		return nil, err
	}

	return s.pkgNames, nil
}

func (s *PackagesScanner) resolveRootSymlink(root string) (string, error) {
	info, err := os.Lstat(root)
	if err != nil {
		return "", err
	}

	if info.Mode()&os.ModeSymlink == os.ModeSymlink {
		root, err = os.Readlink(root)
		if err != nil {
			return "", err
		}
	}

	return root, nil
}

func (s *PackagesScanner) scanImports(path string, info os.FileInfo, err error) (walkErr error) {
	next, walkErr := s.checkIfShouldBeSkipped(path, info)
	if next || walkErr != nil {
		return
	}

	fileSet := token.NewFileSet()
	astFile, walkErr := parser.ParseFile(fileSet, path, nil, parser.ImportsOnly)
	if walkErr != nil {
		return walkErr
	}

	isVendor := strings.Contains(path, "/vendor/")

	for _, importSpec := range astFile.Imports {
		pkgName := strings.Replace(importSpec.Path.Value, "\"", "", -1)
		s.markPackage(pkgName, isVendor)
	}

	return nil
}

func (s *PackagesScanner) markPackage(pkgName string, isVendor bool) {
	if s.pkgNames[pkgName] == nil {
		s.pkgNames[pkgName] = &PackageInfo{
			FromVendor: false,
		}
	}

	s.pkgNames[pkgName].FromVendor = s.pkgNames[pkgName].FromVendor || isVendor
}

func (s *PackagesScanner) checkIfShouldBeSkipped(path string, info os.FileInfo) (bool, error) {
	if info.IsDir() {
		name := info.Name()
		if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_") || name == "testdata" {
			return true, filepath.SkipDir
		}
		return true, nil
	}

	if filepath.Ext(path) != ".go" {
		return true, nil
	}

	return false, nil
}

func NewPackagesScanner() *PackagesScanner {
	return &PackagesScanner{
		pkgNames: make(map[string]*PackageInfo),
	}
}
