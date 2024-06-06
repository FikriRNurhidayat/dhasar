package dhasar

import (
	"fmt"
	"sync"
)

type Container struct {
	mu           sync.RWMutex
	dependencies map[string]interface{}
}

func NewContainer() *Container {
	return &Container{
		dependencies: make(map[string]interface{}),
	}
}

func (c *Container) Register(name string, dependency interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.dependencies[name] = dependency
}

func (c *Container) Resolve(name string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	dependency, ok := c.dependencies[name]
	if !ok {
		return nil, fmt.Errorf("dependency %s not found", name)
	}
	return dependency, nil
}

func Get[T any](c *Container, name string) T {
	dependency, err := c.Resolve(name)
	if err != nil {
		panic(fmt.Sprintf("%s is not resolved.", name))
	}

	dep, ok := dependency.(T)
	if !ok {
		panic(fmt.Sprintf("%s is not resolved.", name))
	}

	return dep
}
