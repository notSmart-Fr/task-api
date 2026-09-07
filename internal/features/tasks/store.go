package tasks

import (
	"errors"
	"sync"
)

var ErrNotFound = errors.New("task not found")

// Store defines the data storage boundary for tasks.
type Store interface {
	GetAll() []Task
	GetByID(id int) (Task, error)
	Create(title string) Task
	Update(id int, input UpdateTaskInput) (Task, error)
	Delete(id int) error
}

// MemoryStore implements Store using an in-memory slice.
type MemoryStore struct {
	sync.RWMutex
	tasks  []Task
	nextID int
}

// NewMemoryStore initializes the store pre-filled with 3 example tasks.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		tasks: []Task{
			{ID: 1, Title: "Learn Go HTTP routing", Done: true},
			{ID: 2, Title: "Build in-memory CRUD API", Done: false},
			{ID: 3, Title: "Write unit tests", Done: false},
		},
		nextID: 4,
	}
}

func (s *MemoryStore) GetAll() []Task {
	s.RLock()
	defer s.RUnlock()
	return s.tasks
}

func (s *MemoryStore) GetByID(id int) (Task, error) {
	s.RLock()
	defer s.RUnlock()

	for _, t := range s.tasks {
		if t.ID == id {
			return t, nil
		}
	}
	return Task{}, ErrNotFound
}

func (s *MemoryStore) Create(title string) Task {
	s.Lock()
	defer s.Unlock()

	t := Task{ID: s.nextID, Title: title, Done: false}
	s.nextID++
	s.tasks = append(s.tasks, t)
	return t
}

func (s *MemoryStore) Update(id int, input UpdateTaskInput) (Task, error) {
	s.Lock()
	defer s.Unlock()

	for i, t := range s.tasks {
		if t.ID == id {
			if input.Title != nil {
				s.tasks[i].Title = *input.Title
			}
			if input.Done != nil {
				s.tasks[i].Done = *input.Done
			}
			return s.tasks[i], nil
		}
	}
	return Task{}, ErrNotFound
}

func (s *MemoryStore) Delete(id int) error {
	s.Lock()
	defer s.Unlock()

	for i, t := range s.tasks {
		if t.ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}
