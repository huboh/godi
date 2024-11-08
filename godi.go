package godi

import (
	"net/http"

	"go.uber.org/dig"
)

// App represents the main application
type App struct {
	*HttpServer
	module    *module
	container *dig.Container
}

// New initializes a new instance of App, configuring the root module and dependencies.
func New(module Module) (*App, error) {
	c := dig.New()
	s := newHttpServer(http.NewServeMux())

	err := c.Provide(func() *HttpServer { return s })
	if err != nil {
		return nil, err
	}

	m, err := newModule(module, c.Scope(GetToken(module)))
	if err != nil {
		return nil, err
	}

	return &App{
		module:     m,
		container:  c,
		HttpServer: s,
	}, nil
}
