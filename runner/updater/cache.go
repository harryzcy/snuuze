package updater

import (
	"os"
	"strings"
)

var (
	readFile = os.ReadFile
)

type Cache struct {
	Files map[string][]byte
}

func NewCache() *Cache {
	return &Cache{
		Files: map[string][]byte{},
	}
}

func (c *Cache) Get(path string) ([]byte, error) {
	if data, ok := c.Files[path]; ok {
		return data, nil
	}
	return c.Read(path)
}

func (c *Cache) Set(path string, data []byte) {
	c.Files[path] = data
}

func (c *Cache) Read(path string) ([]byte, error) {
	data, err := readFile(path)
	if err != nil {
		return nil, err
	}
	c.Set(path, data)
	return data, nil
}

func (c *Cache) Commit() error {
	for path := range c.Files {
		err := os.WriteFile(path, c.Files[path], 0644)
		if err != nil {
			return err
		}
	}
	c.Files = nil
	return nil
}

func (c *Cache) ListGoMod() ([]string, error) {
	var paths []string
	for path := range c.Files {
		if strings.HasSuffix(path, "go.mod") {
			paths = append(paths, path)
		}
	}
	return paths, nil
}
