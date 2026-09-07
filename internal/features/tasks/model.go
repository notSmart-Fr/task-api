package tasks

// Task represents a to-do item in the system.
type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// CreateTaskInput holds payload data for creating a task.
type CreateTaskInput struct {
	Title string `json:"title"`
}

// UpdateTaskInput holds payload data for updating a task.
type UpdateTaskInput struct {
	Title *string `json:"title"`
	Done  *bool   `json:"done"`
}
