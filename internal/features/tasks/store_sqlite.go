package tasks

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore opens (or creates) tasks.db, creates the table if missing, and seeds initial tasks.
func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Stage 0: Create table if not exists
	schema := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		done INTEGER NOT NULL DEFAULT 0
	);`

	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("failed to create tasks table: %w", err)
	}

	store := &SQLiteStore{db: db}

	// Stage 0: Seed example tasks only if the table is empty
	if err := store.seedIfEmpty(); err != nil {
		return nil, fmt.Errorf("failed to seed database: %w", err)
	}

	return store, nil
}

func (s *SQLiteStore) seedIfEmpty() error {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		seeds := []string{
			"Learn Go HTTP routing",
			"Build in-memory CRUD API",
			"Connect CRUD to SQLite database",
		}
		for i, title := range seeds {
			done := 0
			if i == 0 {
				done = 1
			}
			_, err := s.db.Exec("INSERT INTO tasks (title, done) VALUES (?, ?)", title, done)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Stage 1: GET /tasks
func (s *SQLiteStore) GetAll() []Task {
	rows, err := s.db.Query("SELECT id, title, done FROM tasks")
	if err != nil {
		return []Task{}
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		var doneInt int
		if err := rows.Scan(&t.ID, &t.Title, &doneInt); err == nil {
			t.Done = doneInt == 1
			tasks = append(tasks, t)
		}
	}
	if err := rows.Err(); err != nil {
		return []Task{}
	}
	return tasks
}

// Stage 1: GET /tasks/{id}
func (s *SQLiteStore) GetByID(id int) (Task, error) {
	var t Task
	var doneInt int
	err := s.db.QueryRow("SELECT id, title, done FROM tasks WHERE id = ?", id).Scan(&t.ID, &t.Title, &doneInt)
	if err == sql.ErrNoRows {
		return Task{}, ErrNotFound
	} else if err != nil {
		return Task{}, err
	}

	t.Done = doneInt == 1
	return t, nil
}

// Stage 2: POST /tasks
func (s *SQLiteStore) Create(title string) Task {
	res, err := s.db.Exec("INSERT INTO tasks (title, done) VALUES (?, 0)", title)
	if err != nil {
		return Task{}
	}

	id, err := res.LastInsertId()
	if err != nil {
		return Task{}
	}

	return Task{
		ID:    int(id),
		Title: title,
		Done:  false,
	}
}

// Stage 3: PUT /tasks/{id}
func (s *SQLiteStore) Update(id int, input UpdateTaskInput) (Task, error) {
	// Fetch existing first to handle partial updates
	existing, err := s.GetByID(id)
	if err != nil {
		return Task{}, err
	}

	if input.Title != nil {
		existing.Title = *input.Title
	}
	if input.Done != nil {
		existing.Done = *input.Done
	}

	doneInt := 0
	if existing.Done {
		doneInt = 1
	}

	_, err = s.db.Exec("UPDATE tasks SET title = ?, done = ? WHERE id = ?", existing.Title, doneInt, id)
	if err != nil {
		return Task{}, err
	}

	return existing, nil
}

// Stage 3: DELETE /tasks/{id}
func (s *SQLiteStore) Delete(id int) error {
	res, err := s.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
