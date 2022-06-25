package module

type Graph []*Module

func (g Graph) Prune(filter func(*Module) bool) Graph {
	res := excludeItems(g, filter)
	for i, mod := range res {
		res[i] = &Module{
			Name:         mod.Name,
			Path:         mod.Path,
			Dependencies: Graph(mod.Dependencies).Prune(filter),
		}
	}
	return res
}

func excludeItems(list []*Module, filter func(*Module) bool) []*Module {
	result := make([]*Module, 0, len(list))
	for _, mod := range list {
		if !filter(mod) {
			result = append(result, mod)
		}
	}
	return result
}
