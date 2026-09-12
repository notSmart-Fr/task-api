package router

import (
	"errors"
	"io"
	"net/http"
	"slices"
)

// Module defines a feature that owns its routes and resources.
type Module interface {
	RegisterRoutes(mux *http.ServeMux)
	io.Closer
}

func RegisterModules(mux *http.ServeMux, modules ...Module) {
	for _, module := range modules {
		module.RegisterRoutes(mux)
	}
}

func CloseModules(modules ...Module) error {
	var closeErr error
	for _, module := range slices.Backward(modules) {
		if err := module.Close(); err != nil {
			closeErr = errors.Join(closeErr, err)
		}
	}
	return closeErr
}
