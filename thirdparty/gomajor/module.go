package gomajor

// Module contains the module path and versions
type Module struct {
	Path     string
	Versions []string
}

// MultiModule contains multiple modules with different major versions
type MultiModule struct {
	// Modules is a list of modules with different major versions, in ascending order.
	Modules []*Module
}

func (mm *MultiModule) Versions() []string {
	versions := make([]string, 0)
	for _, m := range mm.Modules {
		versions = append(versions, m.Versions...)
	}
	return versions
}
