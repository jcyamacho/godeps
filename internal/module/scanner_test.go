package module

import (
	"fmt"
	"testing"
	"testing/fstest"

	"github.com/jcyamacho/godeps/internal/module/testdata"

	"golang.org/x/exp/slices"
)

func TestScanGoModFiles(t *testing.T) {
	paths := []string{
		"a/b/c/d/go.mod",
		"a/b/c/e/go.mod",
		"a/b/c/go.mod",
	}
	allPaths := []string{
		"a/b/c/f/go1.mod",
		"a/b/c/g/go.1mod",
		"a/b/c/h/go.mod1",
		"a/b/c/i/2go.mod",
	}
	allPaths = append(allPaths, paths...)

	fsys := make(fstest.MapFS)
	for i, path := range allPaths {
		fsys[path] = &fstest.MapFile{
			Data: testdata.NewGoMod(fmt.Sprintf("mod%d", i)).Build(),
		}
	}
	files, err := scanGoModFiles(fsys)
	if err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(files, paths) {
		t.Errorf("Expected %v, got %v", paths, files)
	}
}

func TestScanDir(t *testing.T) {
	modD := testdata.NewGoMod("mod/d")
	modB := testdata.NewGoMod("mod/b").AddDep(modD.Mod)
	modC := testdata.NewGoMod("mod/c").AddDep(modD.Mod)
	modA := testdata.NewGoMod("mod/a").AddDep(modB.Mod).AddDep(modC.Mod).AddIndirectDep(modD.Mod)

	fsys := fstest.MapFS{
		"mod/a/go.mod": {
			Data: modA.Build(),
		},
		"mod/b/go.mod": {
			Data: modB.Build(),
		},
		"mod/c/go.mod": {
			Data: modC.Build(),
		},
		"mod/d/go.mod": {
			Data: modD.Build(),
		},
	}

	graph, err := ScanDir(fsys, false)
	if err != nil {
		t.Fatal(err)
	}

	if len(graph) != 4 {
		t.Errorf("Expected 4 modules, got %d", len(graph))
	}

	a := graph.FindByName(modA.Name)
	b := graph.FindByName(modB.Name)
	c := graph.FindByName(modC.Name)
	d := graph.FindByName(modD.Name)

	if a == nil || b == nil || c == nil || d == nil {
		t.Fatalf("Expected to find modules %s, %s, %s, %s", modA.Name, modB.Name, modC.Name, modD.Name)
	}

	if !a.HasDep(b.Name) || !a.HasDep(c.Name) || !a.HasDep(d.Name) {
		t.Errorf("Expected %s to have deps %s, %s, %s", a.Name, b.Name, c.Name, d.Name)
	}

	if !b.HasDep(d.Name) {
		t.Errorf("Expected %s to have dep %s", b.Name, d.Name)
	}

	if !c.HasDep(d.Name) {
		t.Errorf("Expected %s to have dep %s", c.Name, d.Name)
	}
}
