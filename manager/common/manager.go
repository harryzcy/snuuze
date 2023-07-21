package common

type Manager interface {
	Name() string
	Match(path string) bool
}
