package mermaid

import (
	"strings"

	"github.com/jcyamacho/godeps/internal/module"
)

func Flowchart(modules []*module.Module) string {
	var sb strings.Builder
	sb.WriteString("flowchart\n")
	for _, mod := range modules {
		for _, dep := range mod.Dependencies {
			sb.WriteString("\t")
			sb.WriteString(mod.Name)
			sb.WriteString(" --> ")
			sb.WriteString(dep.Name)
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
