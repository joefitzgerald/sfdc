package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"go/importer"
	"go/types"
)

// GetPackageName finds the package name for the given directory
func GetPackageName(directory, skipPrefix, skipSuffix string) (string, error) {
	pkgDir, err := build.Default.ImportDir(directory, 0)
	if err != nil {
		return "", fmt.Errorf("cannot process directory %s: %s", directory, err)
	}

	var files []*ast.File
	fs := token.NewFileSet()
	for _, name := range pkgDir.GoFiles {
		if !strings.HasSuffix(name, ".go") ||
			(skipSuffix != "" && strings.HasPrefix(name, skipPrefix) &&
				strings.HasSuffix(name, skipSuffix)) {
			continue
		}
		if directory != "." {
			name = filepath.Join(directory, name)
		}
		f, err := parser.ParseFile(fs, name, nil, 0)
		if err != nil {
			return "", fmt.Errorf("parsing file %v: %v", name, err)
		}
		files = append(files, f)
	}
	if len(files) == 0 {
		return "", fmt.Errorf("%s: no buildable Go files", directory)
	}

	// type-check the package
	defs := make(map[*ast.Ident]types.Object)
	config := types.Config{FakeImportC: true, Importer: importer.Default()}
	info := &types.Info{Defs: defs}
	if _, err := config.Check(directory, fs, files, info); err != nil {
		return "", fmt.Errorf("type-checking package: %v", err)
	}

	return files[0].Name.Name, nil
}
