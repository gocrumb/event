package main

import (
	"errors"
	"go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"
	"strings"
)

type Package struct {
	name  string
	dir   string
	files []File
}

func parsePackage(dir string, names []string) (*Package, error) {
	files := []*ast.File{}
	pkg := &Package{}
	pkg.dir = dir
	fs := token.NewFileSet()
	for _, name := range names {
		if !strings.HasSuffix(name, ".go") {
			continue
		}
		file, err := parser.ParseFile(fs, name, nil, 0)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
		pkg.files = append(pkg.files, File{
			File: file,
			pkg:  pkg,
		})
	}
	if len(files) == 0 {
		return nil, errors.New("no buildable Go files")
	}
	pkg.name = files[0].Name.Name
	return pkg, nil
}

func parsePackageFiles(names []string) (*Package, error) {
	return parsePackage(".", names)
}

func parsePackageDir(dir string) (*Package, error) {
	pkg, err := build.Default.ImportDir(dir, 0)
	if err != nil {
		return nil, err
	}
	names := []string{}
	names = append(names, pkg.GoFiles...)
	names = append(names, pkg.CgoFiles...)
	names = append(names, pkg.SFiles...)
	return parsePackage(dir, prefixNames(names, dir))
}

func (p *Package) check(fs *token.FileSet, files []*ast.File) error {
	defs := map[*ast.Ident]types.Object{}
	config := types.Config{
		Importer:    importer.Default(),
		FakeImportC: true,
	}
	_, err := config.Check(p.dir, fs, files, &types.Info{
		Defs: defs,
	})
	return err
}

type File struct {
	*ast.File
	pkg *Package
}

func prefixNames(names []string, dir string) []string {
	if dir == "." {
		return names
	}
	r := make([]string, len(names))
	for i, name := range names {
		r[i] = filepath.Join(dir, name)
	}
	return r
}
