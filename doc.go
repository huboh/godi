// Package godi is a lightweight, modular dependency injection framework designed for building scalable Go applications.
//
// # Overview
//
// Godi simplifies the creation of highly maintainable Go applications through a modular architecture
// powered by dependency injection. It provides a robust foundation for building large-scale applications
// where components are loosely coupled, easily testable, and highly reusable.
//
// Key features include:
//   - Modular architecture with clear separation of concerns
//   - Constructor-based and direct dependency injection
//   - HTTP routing with middleware support via Guards
//   - Flexible module configuration and composition
//   - Type-safe dependency management
//
// # Dependency Injection
//
// Godi supports two main ways of injecting dependencies:
//
//  1. **Constructor-based Injection**:
//
//     Constructors are the building blocks of dependency injection in Godi. They are plain Go functions that:
//     - Accept zero or more dependencies as parameters
//     - Return one or more values of any type
//     - Optionally return an error as the last return value
//
//     Example of a constructor:
//
//     func NewUserService(dep1 *Dependency1, dep2 *Dependency2) (*UserService, error) {
//     return &UserService{
//     dep1: dep1,
//     dep2: dep2,
//     }, nil
//     }
//
//     Any arguments that the constructor has are treated as its dependencies. The dependencies are instantiated
//     in an unspecified order along with any dependencies that they might have, creating a dependency graph at runtime.
//
//  2. **Direct Injection**:
//     If a dependency itself does not require any other dependencies, you can opt to inject it directly without using a constructor.
//     This can be more convenient than defining a constructor, especially for simple dependencies.
//
// # Modules
//
// Modules are the core organizational unit in Godi. Each module encapsulates related functionality
// and can define its dependencies, exports, and HTTP controllers. A module must implement the [godi.Module]
// interface by providing a Config method that returns [*godi.ModuleConfig].
//
// Example of a module:
//
//	type AuthModule struct{}
//
//	func (m *AuthModule) Config() *godi.ModuleConfig {
//		return &godi.ModuleConfig{
//			IsGlobal:         false,                                        // Makes this module's exports available to all other modules
//			Imports:          []godi.Module{&config.Module{}},
//			Exports:          []godi.Provider{},                            // Subset of providers that will be available to other modules
//			ExportsCtor:      []godi.ProviderConstructor{NewAuthService},
//			Providers:        []godi.Provider{},                            // Internal components required by this module
//			ProvidersCtors:   []godi.ProviderConstructor{NewAuthService},
//			Controllers:      []godi.Controller{},                          // HTTP controllers
//			ControllersCtors: []godi.ControllerConstructor{},
//		}
//	}
//
// # Controllers
//
// Controllers handle HTTP routing and request processing. They provide a structured way to define
// endpoints and their associated handlers. A controller must implement the [godi.Controller] interface by
// providing a Config method to return [*godi.ControllerConfig].
//
// Example of a controller:
//
//	type AuthController struct {
//		auth *AuthService
//	}
//
//	func newController(s *AuthService) *AuthController {
//		return &AuthController{
//			auth: s,
//		}
//	}
//
//	func (c *AuthController) Config() *godi.ControllerConfig {
//		return &godi.ControllerConfig{
//			Pattern:     "/auth",                   // Base path for all routes in the controller
//			Metadata:    map[string]string{},       // Controller-wide metadata accessible by guards for decision making
//			Guards:      []godi.Guard{},            // Controller-wide guards
//			GuardsCtors: []godi.GuardConstructor{}, // Controller-wide guards constructors
//			RoutesCfgs:  []*godi.RouteConfig{
//				{
//					Pattern:     "/signin",
//					Method:      http.MethodPost,
//					Handler:     http.HandlerFunc(c.handleSignin),
//					Metadata:    map[string]string{},
//					Guards:      []godi.Guard{},
//					GuardsCtors: []godi.GuardConstructor{},
//				},
//				{
//					Pattern:     "/signup",
//					Method:      http.MethodPost,
//					Handler:     http.HandlerFunc(c.handleSignup),
//					Metadata:    map[string]string{},
//					Guards:      []godi.Guard{},
//					GuardsCtors: []godi.GuardConstructor{},
//				},
//			},
//		}
//	}
//
// # Guards
//
// Guards are used to control access to controllers or individual routes,
// providing an additional layer of security by enforcing runtime or compile-time rules through the incoming request
// or controller/route metadata before handlers are executed.
// They determine whether a given request will be handled by the route handler or not, depending on certain conditions.
//
// A guard must implement the [godi.Guard] interface by providing a Allow method to determine if the request should be allowed.
//
// Example of a guard:
//
//	type AuthGuard struct {
//	    auth *AuthService
//	}
//
//	func newGuard(a *AuthService) *AuthGuard {
//	    return &AuthGuard{
//	        auth: a,
//	    }
//	}
//
//	func (g *AuthGuard) Allow(gCtx godi.GuardContext) (bool, error) {
//	    validated, err := g.auth.Validate(gCtx.Http.R.Header.Get("Authorization"))
//	    if err != nil {
//	        return false, err
//	    }
//
//	    return validated, nil
//	}
//
// Guards can be applied at two different scopes:
//
//  1. **Controller-scoped Guards**:
//
//     Guards defined at the controller level will be applied to all routes within that controller.
//     This is useful for applying authorization to an entire controller.
//
//     func (c *AuthController) Config() *godi.ControllerConfig {
//     return &godi.ControllerConfig{
//     ...
//     Guards:      []godi.Guard{},            // Controller-wide guards
//     GuardsCtors: []godi.GuardConstructor{}, // Controller-wide guards constructors
//     RoutesCfgs:  []*godi.RouteConfig{...},
//     }
//     }
//
//  2. **Route-scoped Guards**:
//
//     Guards can also be defined at the individual route level, allowing you to apply specific guard
//     to only a subset of a controllers routes.
//
//     func (c *AuthController) Config() *godi.ControllerConfig {
//     return &godi.ControllerConfig{
//     ...
//     RoutesCfgs:  []*godi.RouteConfig{
//     {
//     ...
//     Guards:      []godi.Guard{},            // Route guards
//     GuardsCtors: []godi.GuardConstructor{}, // Route guards constructors
//     },
//     },
//     }
//     }
//
// # Structuring Modules
//
// Godi applications are modular, and the main root module is configured to import other modules.
// The root [app.Module] serves as the main entry point, importing any number of sub-modules.
//
//	package app
//
//	import (
//		"github.com/huboh/godi"
//		"github.com/huboh/godi/pkg/modules/config" // Utility module by Godi for reading environment variables
//
//		".../modules/auth"
//		".../modules/database"
//		".../modules/user"
//	)
//
//	type Module struct{}
//
//	func (mod *Module) Config() *godi.ModuleConfig {
//		return &godi.ModuleConfig{
//			Imports: []godi.Module{&config.Module{}, &database.Module{}, &auth.Module{}, &user.Module{}},
//		}
//	}
//
// # Application Setup
//
// To create a new application with Godi, instantiate a new App by providing a root module.
// The following is an example of creating and starting a Godi application.
//
// Creating a Godi application involves defining a root module that composes all other modules:
//
//	package main
//
//	func main() {
//		app, err := godi.New(&app.Module{})
//		if err != nil {
//			log.Fatal("failed to create godi app: ", err)
//		}
//
//	    // Start the application's HTTP server, listening on the specified host and port.
//		err = app.Listen("localhost", "5000")
//		if err != nil {
//			log.Fatal("failed to start app server: ", err)
//		}
//	}
//
// # Best Practices
//
//   - Keep modules focused and cohesive - each module should have a single responsibility
//   - Use constructors to ensure proper initialization and dependency validation
//   - Use metadata to add extra information to routes and controllers for documentation or tooling
//   - Structure your application with clear module boundaries and well-defined interfaces
//   - Consider making commonly used services global by setting `IsGlobal: true` in their module config
package godi // import "github.com/huboh/godi"
