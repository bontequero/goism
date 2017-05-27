package tu

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
)

func parse(fSet *token.FileSet, pkgPath string) (*ast.Package, error) {
	const parseFlags = parser.ParseComments

	pkgs, err := parser.ParseDir(fSet, pkgPath, nil, parseFlags)
	if err != nil {
		return nil, err
	}

	if len(pkgs) > 1 {
		pkgNames := reflect.ValueOf(pkgs).MapKeys()
		return nil, fmt.Errorf("more than one package found: %v", pkgNames)
	}
	for _, pkg := range pkgs {
		return pkg, nil
	}
	return nil, fmt.Errorf("can not find Go package in %s", pkgPath)
}
