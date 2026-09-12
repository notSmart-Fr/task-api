package tasks

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	db *sql.DB
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
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

// Get/tasks with pagination
func (s *SQLiteStore) GetAll(limit, offset int) ([]Task, int, error) {
	// 1. Get total count for pagination metadata
	var total int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&total); err != nil {
		return nil, 0, err
	}

	// 2. Fetch paginated records
	query := "SELECT id, title, done FROM tasks ORDER BY id ASC LIMIT ? OFFSET ?"
	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tasksList []Task
	for rows.Next() {
		var t Task
		var doneInt int
		if err := rows.Scan(&t.ID, &t.Title, &doneInt); err != nil {
			return nil, 0, err
		}
		t.Done = doneInt == 1
		tasksList = append(tasksList, t)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	if tasksList == nil {
		tasksList = []Task{}
	}

	return tasksList, total, nil
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
func (s *SQLiteStore) Create(title string) (Task, error) {
	res, err := s.db.Exec("INSERT INTO tasks (title, done) VALUES (?, 0)", title)
	if err != nil {
		return Task{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return Task{}, err
	}

	return Task{
		ID:    int(id),
		Title: title,
		Done:  false,
	}, nil
}

// Stage 3: PUT /tasks/{id}
func (s *SQLiteStore) Update(id int, title *string, done *bool) (Task, error) {
	existing, err := s.GetByID(id)
	if err != nil {
		return Task{}, err
	}

	if title != nil {
		existing.Title = *title
	}
	if done != nil {
		existing.Done = *done
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
