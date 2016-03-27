package main

import (
	"go/ast"

	_ "golang.org/x/tools/go/gcimporter"
)

type FuncDefinition struct {
	Name string
}

func (fd *FuncDefinition) Parse(list []*ast.Field) bool {
	return true
}
