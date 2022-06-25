package module

import (
	"testing"
	"testing/fstest"

	"golang.org/x/exp/slices"
)

func moduleEqual(a, b *Module) bool { return a.Name == b.Name }

func FuzzNew(f *testing.F) {
	f.Add("mod1")
	f.Add("github.com/org1/mod1")
	f.Add("gitlab.com/org1/mod2")

	f.Fuzz(func(t *testing.T, name string) {
		m := New(name)
		if m.Name != name {
			t.Errorf("Expected %s, got %s", name, m.Name)
		}
	})
}

func TestModule_AddDependency(t *testing.T) {
	tests := []struct {
		name string
		deps []*Module
	}{
		{
			name: "case1",
			deps: []*Module{
				New("a"),
				New("b"),
			},
		},
		{
			name: "case2",
			deps: []*Module{
				New("mod1"),
			},
		},
		{
			name: "case3",
			deps: []*Module{
				New("a"),
				New("c"),
				New("b"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := New(test.name)
			m.AddDependency(test.deps...)

			if !slices.EqualFunc(m.Dependencies, test.deps, moduleEqual) {
				t.Errorf("Expected %v, got %v", test.deps, m.Dependencies)
			}
		})
	}
}

func TestFromFile(t *testing.T) {
	mFs := fstest.MapFS{
		"go.mod": {
			Data: []byte(`
				module github.com/jcyamacho/godeps

				go 1.18
				
				require (
					github.com/hashicorp/go-multierror v1.1.1
					github.com/spf13/cobra v1.5.0
					golang.org/x/mod v0.5.1
				)
				
				require (
					github.com/hashicorp/errwrap v1.1.0 // indirect
					github.com/inconshreveable/mousetrap v1.0.0 // indirect
					github.com/spf13/pflag v1.0.5 // indirect
					golang.org/x/xerrors v0.0.0-20220609144429-65e65417b02f // indirect
				)
			`),
		},
	}

	mod, err := FromFile(mFs, "go.mod", false)
	if err != nil {
		t.Fatal(err)
	}

	if mod.Name != "github.com/jcyamacho/godeps" {
		t.Errorf("Expected %s, got %s", "github.com/jcyamacho/godeps", mod.Name)
	}

	deps := Graph(mod.Dependencies)
	if len(deps) != 7 {
		t.Errorf("Expected %d dependencies, got %d", 7, len(deps))
	}

	names := []string{
		"github.com/hashicorp/go-multierror",
		"github.com/spf13/cobra",
		"golang.org/x/mod",
		"github.com/hashicorp/errwrap",
		"github.com/inconshreveable/mousetrap",
		"github.com/spf13/pflag",
		"golang.org/x/xerrors",
	}

	for _, depName := range names {
		if dep := deps.FindByName(depName); dep == nil {
			t.Errorf("Expected dependency %s, got nil", depName)
		}
	}
}
