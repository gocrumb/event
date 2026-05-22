package main

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"testing"
)

const fixtureSrc = `package fixture

type OrderPlaced struct {
	OrderID    int
	CustomerID int
}

type OrderShipped struct {
	OrderID int
}
`

func generateFixture(t *testing.T) (string, []byte) {
	t.Helper()
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "fixture.go")
	if err := os.WriteFile(srcPath, []byte(fixtureSrc), 0644); err != nil {
		t.Fatal(err)
	}
	pkg, err := parsePackage(dir, []string{srcPath})
	if err != nil {
		t.Fatalf("parsePackage: %v", err)
	}
	g := Generator{
		pkg:   pkg,
		types: []string{"OrderPlaced", "OrderShipped"},
	}
	out, err := g.generate()
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	return dir, out
}

func TestGenerateParses(t *testing.T) {
	_, out := generateFixture(t)
	fs := token.NewFileSet()
	if _, err := parser.ParseFile(fs, "event.go", out, 0); err != nil {
		t.Fatalf("generated code did not parse: %v\n%s", err, out)
	}
}

func TestGenerateTypeChecks(t *testing.T) {
	dir, out := generateFixture(t)
	fs := token.NewFileSet()
	srcFile, err := parser.ParseFile(fs, filepath.Join(dir, "fixture.go"), nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	genFile, err := parser.ParseFile(fs, filepath.Join(dir, "event.go"), out, 0)
	if err != nil {
		t.Fatal(err)
	}
	conf := types.Config{Importer: importer.Default()}
	pkg, err := conf.Check("fixture", fs, []*ast.File{srcFile, genFile}, &types.Info{
		Defs: map[*ast.Ident]types.Object{},
	})
	if err != nil {
		t.Fatalf("type check failed: %v\n%s", err, out)
	}

	scope := pkg.Scope()
	for _, name := range []string{
		"Event",
		"EventType",
		"EventOrderPlaced",
		"EventOrderShipped",
		"OrderPlacedHandler",
		"OrderPlacedHandlerFunc",
		"OrderPlacedEmitter",
		"OrderShippedHandler",
		"OrderShippedHandlerFunc",
		"OrderShippedEmitter",
		"EmitterOrderPlaced",
		"EmitterOrderShipped",
		"OnOrderPlaced",
		"OnOrderShipped",
		"EmitOrderPlaced",
		"EmitOrderShipped",
		"Trigger",
	} {
		if scope.Lookup(name) == nil {
			t.Errorf("expected symbol %q in generated package", name)
		}
	}
}

func TestGeneratedSignaturesUseContext(t *testing.T) {
	dir, out := generateFixture(t)
	fs := token.NewFileSet()
	srcFile, err := parser.ParseFile(fs, filepath.Join(dir, "fixture.go"), nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	genFile, err := parser.ParseFile(fs, filepath.Join(dir, "event.go"), out, 0)
	if err != nil {
		t.Fatal(err)
	}
	conf := types.Config{Importer: importer.Default()}
	pkg, err := conf.Check("fixture", fs, []*ast.File{srcFile, genFile}, nil)
	if err != nil {
		t.Fatalf("type check failed: %v", err)
	}

	cases := []struct {
		name   string
		params []string
	}{
		{"EmitOrderPlaced", []string{"context.Context", "fixture.OrderPlaced"}},
		{"EmitOrderShipped", []string{"context.Context", "fixture.OrderShipped"}},
		{"OnOrderPlaced", []string{"func(context.Context, fixture.OrderPlaced)"}},
		{"OnOrderShipped", []string{"func(context.Context, fixture.OrderShipped)"}},
		{"Trigger", []string{"context.Context", "fixture.Event"}},
	}
	for _, c := range cases {
		obj := pkg.Scope().Lookup(c.name)
		if obj == nil {
			t.Errorf("missing symbol %q", c.name)
			continue
		}
		sig, ok := obj.Type().(*types.Signature)
		if !ok {
			t.Errorf("%s: not a function (%T)", c.name, obj.Type())
			continue
		}
		if sig.Params().Len() != len(c.params) {
			t.Errorf("%s: want %d params, got %d", c.name, len(c.params), sig.Params().Len())
			continue
		}
		for i, want := range c.params {
			if got := sig.Params().At(i).Type().String(); got != want {
				t.Errorf("%s param %d: want %s, got %s", c.name, i, want, got)
			}
		}
	}

	// Event interface must require Trigger(context.Context).
	iface, ok := pkg.Scope().Lookup("Event").Type().Underlying().(*types.Interface)
	if !ok {
		t.Fatalf("Event is not an interface")
	}
	var trigger *types.Func
	for i := 0; i < iface.NumMethods(); i++ {
		if iface.Method(i).Name() == "Trigger" {
			trigger = iface.Method(i)
			break
		}
	}
	if trigger == nil {
		t.Fatal("Event.Trigger missing")
	}
	sig := trigger.Type().(*types.Signature)
	if sig.Params().Len() != 1 || sig.Params().At(0).Type().String() != "context.Context" {
		t.Errorf("Event.Trigger: want (context.Context), got %s", sig)
	}
}
