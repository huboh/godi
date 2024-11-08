package godi

import (
	"fmt"

	"go.uber.org/dig"
)

type (
	// Module is an interface representing a self-contained unit of functionality.
	// It exposes providers that can be imported by other modules within the application.
	Module interface {
		// Config returns the module config containing providers, exports
		// and imports required by the module.
		Config() *ModuleConfig
	}

	// ModuleConfig provides configuration settings for a module's functionality,
	// such as its providers, controllers, and imported modules.
	ModuleConfig struct {
		// IsGlobal indicates whether the module's exported providers are globally
		// available to every other module without needing an explicit import.
		//
		// This is useful for shared utilities or database connections.
		IsGlobal bool

		// Imports specifies other modules that this module depends on.
		// Providers from these modules will be available within this module.
		Imports []Module

		// Exports lists providers from this module that should be accessible to
		// other modules that import this module.
		Exports []Provider

		// ExportsCtors lists constructors for providers that should be accessible
		// in other modules importing this module.
		ExportsCtors []ProviderConstructor

		// Providers lists the providers within the module that are shared across
		// the module's other providers.
		Providers []Provider

		// ProvidersCtors lists constructors for providers that the Godi injector
		// will create and share within this module.
		ProvidersCtors []ProviderConstructor

		// Controllers lists the handlers defined in this module, which handle
		// HTTP requests and define the module's endpoints.
		Controllers []Controller

		// ControllersCtors lists constructors for controllers in this module that
		// will be instantiated by the Godi injector.
		ControllersCtors []ControllerConstructor
	}
)

// module is a wrapper for managing an instance of a Module.
type module struct {
	Module
	scope   scope
	parent  *module
	imports []*module
}

func newModule(m Module, s scope) (*module, error) {
	var (
		err error
		mod = &module{
			scope:  s,
			Module: m,
		}
	)

	err = mod._registerProviders()
	if err != nil {
		return nil, fmt.Errorf("error registering providers: %w", err)
	}

	err = mod._registerControllers()
	if err != nil {
		return nil, fmt.Errorf("error registering controllers: %w", err)
	}

	// recursively create imported modules
	for _, imported := range mod.Config().Imports {
		importedMod, err := newModule(imported, mod.newChildScope(imported))
		if err != nil {
			return nil, fmt.Errorf("error building module (%T): %w", imported, err)
		}

		err = importedMod.assignParent(mod)
		if err != nil {
			return nil, err
		}

		err = importedMod._registerExportedProviders()
		if err != nil {
			return nil, fmt.Errorf("error registering exports: %w", err)
		}
	}

	return mod, nil
}

// assignParent assigns the module's parent and append itself to the parent import list
func (m *module) assignParent(parent *module) error {
	if parent != nil {
		m.parent = parent
		m.parent.imports = append(m.parent.imports, m)
	}
	return nil
}

func (m *module) _registerProviders() error {
	mCfg := m.Config()
	for _, pvdCtor := range mCfg.ProvidersCtors {
		isGlobExport := (mCfg.IsGlobal && m.isExportedProvider(pvdCtor))

		// a global module's exported providers
		// should be made available to all available scopes
		err := m.scope.Provide(pvdCtor, dig.Export(isGlobExport))
		if err != nil {
			return fmt.Errorf("error providing provider (%T): %w", pvdCtor, err)
		}
	}
	return nil
}

// _registerControllers registers controllers in the group named "controllers" in the module scope
func (m *module) _registerControllers() error {
	var (
		mCfg = m.Config()
		opts = []dig.ProvideOption{
			dig.As(new(Controller)),
			dig.Group(groupControllers.String()),
		}
	)

	for _, ctrlCtor := range mCfg.ControllersCtors {
		err := m.scope.Provide(ctrlCtor, opts...)
		if err != nil {
			return fmt.Errorf("error providing controller (%T): %w", ctrlCtor, err)
		}
	}

	return m.scope.Invoke(
		func(input controllerGroupInput) error {
			for _, controller := range input.Controllers {
				_, err := newController(controller, m)
				if err != nil {
					return err
				}
			}
			return nil
		},
	)
}

// _registerExportedProviders registers the current module's exports in it's parent scope
func (m *module) _registerExportedProviders() error {
	mCfg := m.Config()
	// a global module's exports would be
	// available to all available scopes already.
	// so, no need to provide it to parent module's scope
	if (mCfg.IsGlobal) || (m.parent == nil) {
		return nil
	}

	for _, pvdCtor := range mCfg.ExportsCtors {
		err := m.parent.scope.Provide(pvdCtor)
		if err != nil {
			return fmt.Errorf("error providing export (%T): %w", pvdCtor, err)
		}
	}

	return nil
}

func (m *module) newChildScope(mod Module) scope {
	return m.scope.Scope(GetToken(mod))
}

func (m *module) isExportedProvider(provider ProviderConstructor) bool {
	for _, export := range m.Config().ExportsCtors {
		if GetToken(export) == GetToken(provider) {
			return true
		}
	}
	return false
}
