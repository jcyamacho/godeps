package module

import (
	"testing"
)

// a -> b, a -> c, b -> d, c -> d
func createGraph() Graph {
	modA := New("a")
	modB := New("b")
	modC := New("c")
	modD := New("d")

	modA.AddDependency(modB, modC)
	modB.AddDependency(modD)
	modC.AddDependency(modD)

	return Graph{modA, modB, modC, modD}
}

func TestPrune(t *testing.T) {
	g := createGraph()

	g = g.Prune(func(m *Module) bool {
		return m.Name == "b"
	})

	if len(g) != 3 {
		t.Errorf("Expected 3 modules, got %d", len(g))
	}

	if g[0].Name != "a" || g[1].Name != "c" || g[2].Name != "d" {
		t.Errorf("Expected a, c, d, got %v", g)
	}

	if len(g[0].Dependencies) != 1 || g[0].Dependencies[0].Name != "c" {
		t.Errorf("Expected a -> c, got %v", g[0].Dependencies)
	}
}

func TestFindByName(t *testing.T) {
	g := createGraph()

	expecteds := []string{"a", "c", "d"}
	unexpecteds := []string{"x", "f", "e"}

	for _, expected := range expecteds {
		if m := g.FindByName(expected); m == nil {
			t.Errorf("Expected %s, got nil", expected)
		}
	}

	for _, unexpected := range unexpecteds {
		if m := g.FindByName(unexpected); m != nil {
			t.Errorf("Expected nil, got %s", unexpected)
		}
	}
}
