package dhasar_container

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

func Get[T any](c *Container, name string) (T, error) {
	dependency, err := c.Resolve(name)
	if err != nil {
		var zero T
		return zero, err
	}
	dep, ok := dependency.(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("dependency %s is of incorrect type", name)
	}
	return dep, nil
}
