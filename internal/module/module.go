package module

import (
	"io/ioutil"

	"golang.org/x/mod/modfile"
)

type Module struct {
	Name         string
	Path         string
	Dependencies []*Module
}

func (m *Module) AddDependency(module *Module) {
	m.Dependencies = append(m.Dependencies, module)
}

func New(name string) *Module {
	return &Module{
		Name: name,
	}
}

func FromFile(path string, skipIndirectDeps bool) (*Module, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	file, err := modfile.ParseLax(path, data, nil)
	if err != nil {
		return nil, err
	}

	module := New(file.Module.Mod.Path)
	module.Path = path

	for _, r := range file.Require {
		if skipIndirectDeps && r.Indirect {
			continue
		}
		module.AddDependency(New(r.Mod.Path))
	}

	return module, nil
}
