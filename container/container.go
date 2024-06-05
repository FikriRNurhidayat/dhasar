package dhasar_container

import (
	"fmt"
	"reflect"
)

type Container struct {
	dependencies map[reflect.Type]interface{}
}

func New() *Container {
	return &Container{
		dependencies: make(map[reflect.Type]interface{}),
	}
}

func (c *Container) Register(dependency interface{}) {
	t := reflect.TypeOf(dependency)
	c.dependencies[t] = dependency
}

func Get[T any](c *Container) T {
	t := reflect.TypeOf((*T)(nil)).Elem()
	dependency, exists := c.dependencies[t]
	if !exists {
		panic(fmt.Sprintf("no dependency found for type %v", t))
	}
	return dependency.(T)
}
