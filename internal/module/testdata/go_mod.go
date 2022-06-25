package testdata

import "bytes"

type Mod struct {
	Name    string
	Version string
}

type GoMod struct {
	*Mod
	Deps         []*Mod
	IndirectDeps []*Mod
}

func NewGoMod(name string) *GoMod {
	return &GoMod{
		Mod: &Mod{
			Name:    name,
			Version: "v1.0.0",
		},
	}
}

func (g *GoMod) AddDep(d *Mod) *GoMod {
	g.Deps = append(g.Deps, d)
	return g
}

func (g *GoMod) AddIndirectDep(d *Mod) *GoMod {
	g.IndirectDeps = append(g.IndirectDeps, d)
	return g
}

func (g *GoMod) Build() []byte {
	var buf bytes.Buffer
	buf.WriteString("module " + g.Name + "\n")
	buf.WriteString("go 1.18\n")
	for _, dep := range g.Deps {
		buf.WriteString("require " + dep.Name + " " + dep.Version + "\n")
	}
	for _, dep := range g.IndirectDeps {
		buf.WriteString("require " + dep.Name + " " + dep.Version + " // indirect\n")
	}
	return buf.Bytes()
}
