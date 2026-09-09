package tasks

import "net/http"

func RegisterHandlers(mux *http.ServeMux, store Store) {
	mux.HandleFunc("GET /tasks", NewListTasksEndpoint(store).Handle)
	mux.HandleFunc("GET /tasks/{id}", NewGetTaskEndpoint(store).Handle)
	mux.HandleFunc("POST /tasks", NewCreateTaskEndpoint(store).Handle)
	mux.HandleFunc("PUT /tasks/{id}", NewUpdateTaskEndpoint(store).Handle)
	mux.HandleFunc("DELETE /tasks/{id}", NewDeleteTaskEndpoint(store).Handle)
}
