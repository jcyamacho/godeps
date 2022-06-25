package module

import (
	"io/fs"
	"sort"
	"sync"

	"github.com/hashicorp/go-multierror"
)

const goModFileName = "go.mod"

func ScanDir(fsys fs.FS, skipIndirectDeps bool) (Graph, error) {
	fileMods, err := scanGoModFiles(fsys)
	if err != nil {
		return nil, err
	}

	modules := &sync.Map{}
	g := &multierror.Group{}
	for _, file := range fileMods {
		file := file
		g.Go(func() error {
			mod, err := FromFile(fsys, file, skipIndirectDeps)
			if err != nil {
				return err
			}
			modules.Store(mod.Name, mod)
			for _, dep := range mod.Dependencies {
				modules.LoadOrStore(dep.Name, dep)
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return toModulesList(modules), nil
}

func scanGoModFiles(fsys fs.FS) ([]string, error) {
	var files []string
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == goModFileName {

			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func toModulesList(m *sync.Map) []*Module {
	var list []*Module
	m.Range(func(_, value any) bool {
		mod := value.(*Module)
		for i, dep := range mod.Dependencies {
			if d, loaded := m.Load(dep.Name); loaded {
				mod.Dependencies[i] = d.(*Module)
			}
		}
		sort.Slice(mod.Dependencies, func(i, j int) bool {
			return mod.Dependencies[i].Name < mod.Dependencies[j].Name
		})
		list = append(list, mod)
		return true
	})
	sort.Slice(list, func(i, j int) bool {
		return len(list[i].Dependencies) > len(list[j].Dependencies)
	})
	return list
}
