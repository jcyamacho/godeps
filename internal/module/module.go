package module

import (
	"io"
	"io/fs"

	"golang.org/x/mod/modfile"
)

type Module struct {
	Name         string
	Path         string
	Dependencies []*Module
}

func (m *Module) AddDependency(modules ...*Module) {
	m.Dependencies = append(m.Dependencies, modules...)
}

func (m *Module) HasDep(name string) bool {
	for _, dep := range m.Dependencies {
		if dep.Name == name {
			return true
		}
	}
	return false
}

func New(name string) *Module {
	return &Module{
		Name: name,
	}
}

func FromFile(fsys fs.FS, path string, skipIndirectDeps bool) (*Module, error) {
	data, err := readFile(fsys, path)
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

func readFile(fsys fs.FS, path string) ([]byte, error) {
	file, err := fsys.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return io.ReadAll(file)
}
