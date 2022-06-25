package mermaid

import (
	"testing"

	"github.com/jcyamacho/godeps/internal/module"
)

func TestFlowchart(t *testing.T) {
	modA := module.New("a")
	modB := module.New("b")
	modC := module.New("c")
	modD := module.New("d")

	modA.AddDependency(modB, modC)
	modB.AddDependency(modD)
	modC.AddDependency(modD)

	data := Flowchart([]*module.Module{modA, modB, modC, modD})
	expected := "flowchart\n\ta --> b\n\ta --> c\n\tb --> d\n\tc --> d\n"

	if data != expected {
		t.Errorf("Expected %s, got %s", expected, data)
	}
}
