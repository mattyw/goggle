package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
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
	params := []*ast.Field{}
	results := []*ast.Field{}
	pos := fd.Type.Func
	if fd.Type.Params.List != nil {

		params = fd.Type.Params.List
	}
	if fd.Type.Results != nil {
		results = fd.Type.Results.List
	}

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
	p := os.Getenv("GOPATH")
	if p == "" {
		return
	}
	gopath := strings.Split(p, ":")
	for i, d := range gopath {
		gopath[i] = filepath.Join(d, "src")
	}
	r := runtime.GOROOT()
	if r != "" {
		gopath = append(gopath, r+"/src/pkg")
	}
	for _, path := range gopath {
		funcs, err := walk(path)
		if err != nil {
			fmt.Printf("failed to walk %v", err)
			return
		}
		fmt.Println(len(funcs))
	}
}

func walk(path string) ([]*F, error) {
	funcs := []*F{}
	walker := func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".go" {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return nil // Ignore error for now
			}
			f, err := inspectFile(path, string(data))
			funcs = append(funcs, f...)
			if err != nil {
				return nil
			}
		}
		return nil
	}
	err := filepath.Walk(path, walker)
	if err != nil {
		return nil, err
	}
	return funcs, nil
}

// TODO we probably need to return []F from here
func inspectFile(filename, contents string) ([]*F, error) {
	funcs := []*ast.FuncDecl{}
	result := []*F{}
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, filename, contents, 0)
	if err != nil {
		return nil, err
	}

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			funcs = append(funcs, x)
		}
		return true
	})
	for _, f := range funcs {
		fun := NewF(f, fset)
		result = append(result, fun)
	}
	return result, nil
}
