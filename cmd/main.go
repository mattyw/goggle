package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type F struct {
	decl   *ast.FuncDecl
	fset   *token.FileSet
	Input  []string
	Output []string
	Pos    token.Pos
}

func (f *F) String() string {
	pos := f.fset.Position(f.Pos)
	return fmt.Sprintf("%s (%s) %s", pos, strings.Join(f.Input, ","), strings.Join(f.Output, ","))
}

func NewF(fd *ast.FuncDecl, fs *token.FileSet) *F {
	pos := fd.Type.Func
	params := fd.Type.Params.List
	results := fd.Type.Results.List

	input := []string{}
	for _, i := range params {
		for _ = range i.Names {
			input = append(input, fmt.Sprintf("%s", (i.Type)))
		}
	}
	output := []string{}
	for _, o := range results {
		output = append(output, fmt.Sprintf("%s", (o.Type)))
	}
	return &F{fset: fs, decl: fd, Input: input, Output: output, Pos: pos}
}

func main() {
	funcs := []*ast.FuncDecl{}
	// src is the input for which we want to inspect the AST.
	src := `
package p
const c = 1.0
var X = f(3.14)*2 + c
func foo(x,y int) int {
    return x + y
}
func foo2(x int, y float64) int {
    return x + y
}
func bar(x int, y int) (int, int) {
    return x,y
}
`

	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "src.go", src, 0)
	if err != nil {
		panic(err)
	}

	// Inspect the AST and print all identifiers and literals.
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			funcs = append(funcs, x)
		}
		return true
	})
	for _, f := range funcs {
		fun := NewF(f, fset)
		fmt.Println(fun.String())
	}

}
