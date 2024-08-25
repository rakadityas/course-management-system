package studentdomain

import (
	"context"
	"database/sql"
	"fmt"
)

// StudentRepository defines the interface for student-related database operations.
type StudentRepository interface {
	GetStudentByID(ctx context.Context, id int64) (*Student, error)
}

// StudentDB implements the StudentRepository interface using a SQL database.
type StudentDB struct {
	DB *sql.DB
}

// NewSQLStudentRepository creates a new StudentDB instance with the given database connection.
func NewSQLStudentRepository(db *sql.DB) *StudentDB {
	return &StudentDB{DB: db}
}

// GetStudentByID retrieves a student from the database by their ID.
func (repo *StudentDB) GetStudentByID(ctx context.Context, id int64) (*Student, error) {
	query := `
		SELECT id, email, create_time, update_time
		FROM students
		WHERE id = ?
	`
	row := repo.DB.QueryRowContext(ctx, query, id)

	student := &Student{}
	err := row.Scan(&student.ID, &student.Email, &student.CreateTime, &student.UpdateTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No student found
		}
		return nil, fmt.Errorf("failed to retrieve student: %v", err)
	}

	return student, nil
}
