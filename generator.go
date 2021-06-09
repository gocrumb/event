package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"strings"
)

type Generator struct {
	buf bytes.Buffer
	pkg *Package

	types []string
}

func (g *Generator) generate() ([]byte, error) {
	b := bytes.Buffer{}

	types := []Type{}
	for _, f := range g.pkg.files {
		ast.Inspect(f.File, newGenDeclFn(g, &types))
	}

	fmt.Fprintf(&b, "// Code generated by \"gocrumb event %s\"; DO NOT EDIT\n", strings.Join(os.Args[1:], " "))
	fmt.Fprintln(&b)
	fmt.Fprintf(&b, "package %s\n", g.pkg.name)
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "type EventType int")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "const (")
	for i, t := range types {
		fmt.Fprintf(&b, "\tEvent%s", t.name)
		if i == 0 {
			fmt.Fprintln(&b, " EventType = iota")
		} else {
			fmt.Fprintln(&b)
		}
	}
	fmt.Fprintln(&b, ")")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "type Event interface {")
	fmt.Fprintln(&b, "\tType() EventType")
	fmt.Fprintln(&b, "\tTrigger()")
	fmt.Fprintln(&b, "}")
	fmt.Fprintln(&b)
	for _, t := range types {
		fmt.Fprintf(&b, "func (e %s) Type() EventType {\n", t.name)
		fmt.Fprintf(&b, "\treturn Event%s\n", t.name)
		fmt.Fprintln(&b, "}")
		fmt.Fprintln(&b)
		fmt.Fprintf(&b, "func (e %s) Trigger() {\n", t.name)
		fmt.Fprintf(&b, "\tEmitter%s.Trigger(e)\n", t.name)
		fmt.Fprintln(&b, "}")
		fmt.Fprintln(&b)
		fmt.Fprintf(&b, "type %sHandler interface {\n", t.name)
		fmt.Fprintf(&b, "\tHandle(%s)\n", t.name)
		fmt.Fprintln(&b, "}")
		fmt.Fprintln(&b)
		fmt.Fprintf(&b, "type %sHandlerFunc func(%s)\n", t.name, t.name)
		fmt.Fprintln(&b)
		fmt.Fprintf(&b, "func (f %sHandlerFunc) Handle(e %s) {\n", t.name, t.name)
		fmt.Fprintln(&b, "\tf(e)")
		fmt.Fprintln(&b, "}")
		fmt.Fprintln(&b)
		fmt.Fprintf(&b, "type %sEmitter struct {\n", t.name)
		fmt.Fprintf(&b, "\thandlers []%sHandler", t.name)
		fmt.Fprintln(&b, "}")
		fmt.Fprintln(&b)
		fmt.Fprintf(&b, "func (m *%sEmitter) Trigger(e %s) {\n", t.name, t.name)
		fmt.Fprintln(&b, "\tfor _, h := range m.handlers {")
		fmt.Fprintln(&b, "\t\th.Handle(e)")
		fmt.Fprintln(&b, "\t}")
		fmt.Fprintln(&b, "}")
		fmt.Fprintln(&b)
		fmt.Fprintf(&b, "func (m *%sEmitter) Handle(h %sHandler) {\n", t.name, t.name)
		fmt.Fprintln(&b, "\tm.handlers = append(m.handlers, h)")
		fmt.Fprintln(&b, "}")
		fmt.Fprintln(&b)
		fmt.Fprintf(&b, "func (m *%sEmitter) HandleFunc(f func(%s)) {\n", t.name, t.name)
		fmt.Fprintf(&b, "\tm.Handle(%sHandlerFunc(f))\n", t.name)
		fmt.Fprintln(&b, "}")
		fmt.Fprintln(&b)
		fmt.Fprintf(&b, "func On%s(f func(%s)) {\n", t.name, t.name)
		fmt.Fprintf(&b, "\tEmitter%s.Handle(%sHandlerFunc(f))\n", t.name, t.name)
		fmt.Fprintln(&b, "}")
		fmt.Fprintln(&b)
		fmt.Fprintf(&b, "func Emit%s(e %s) {\n", t.name, t.name)
		fmt.Fprintf(&b, "\tEmitter%s.Trigger(e)\n", t.name)
		fmt.Fprintln(&b, "}")
		fmt.Fprintln(&b)
	}
	fmt.Fprintln(&b, "var (")
	for _, t := range types {
		fmt.Fprintf(&b, "\tEmitter%s = %sEmitter{}\n", t.name, t.name)
	}
	fmt.Fprintln(&b, ")")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "func Trigger(e Event) {")
	fmt.Fprintln(&b, "\te.Trigger()")
	fmt.Fprintln(&b, "}")

	return format.Source(b.Bytes())
}

func newGenDeclFn(g *Generator, types *[]Type) func(node ast.Node) bool {
	return func(node ast.Node) bool {
		decl, ok := node.(*ast.GenDecl)
		if !ok || decl.Tok != token.TYPE {
			return true
		}

		for _, spec := range decl.Specs {
			tspec := spec.(*ast.TypeSpec)
			for _, t := range g.types {
				if t == tspec.Name.Name {
					*types = append(*types, Type{
						name: t,
					})
					break
				}
			}
		}

		return false
	}
}

type Type struct {
	name string
}
