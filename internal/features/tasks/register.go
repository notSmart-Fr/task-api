package tasks

import (
	"io"
	"net/http"
)

type Module struct {
	store Store
}

func NewModule(store Store) *Module {
	return &Module{store: store}
}

func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /tasks", NewListTasksEndpoint(m.store).Handle)
	mux.HandleFunc("GET /tasks/{id}", NewGetTaskEndpoint(m.store).Handle)
	mux.HandleFunc("POST /tasks", NewCreateTaskEndpoint(m.store).Handle)
	mux.HandleFunc("PUT /tasks/{id}", NewUpdateTaskEndpoint(m.store).Handle)
	mux.HandleFunc("DELETE /tasks/{id}", NewDeleteTaskEndpoint(m.store).Handle)
}

func (m *Module) Close() error {
	if closer, ok := m.store.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
