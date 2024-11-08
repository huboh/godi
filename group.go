package godi

import "go.uber.org/dig"

// group is a type that represents a specific group
// of values of the same type in the dependency injection container/scope
type group string

func (g group) String() string {
	return string(g)
}

const (
	groupPipes        group = "pipes"
	groupGuards       group = "guards"
	groupFilters      group = "filters"
	groupProviders    group = "providers"
	groupControllers  group = "controllers"
	groupInterceptors group = "interceptors"
)

// guardGroupInput is used for injecting the collection of Guard instances
// grouped under `groupGuards` in a particular.
type guardGroupInput struct {
	dig.In
	Guards []Guard `group:"guards"`
}

// controllerGroupInput is used for injecting the collection of Controller instances
// grouped under `groupControllers` in a particular.
type controllerGroupInput struct {
	dig.In
	Controllers []Controller `group:"controllers"`
}
