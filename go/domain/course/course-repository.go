package coursedomain

import (
	"context"
	"database/sql"
	"fmt"
)

type CourseRepository interface {
	GetCourseByID(ctx context.Context, id int64) (*Course, error)
}

type CourseDB struct {
	DB *sql.DB
}

// NewSQLCourseRepository creates a new StudentDB instance with the given database connection.
func NewSQLCourseRepository(db *sql.DB) *CourseDB {
	return &CourseDB{DB: db}
}

// GetCourseByID retrieves a course by its ID from the database.
func (repo *CourseDB) GetCourseByID(ctx context.Context, id int64) (*Course, error) {
	query := `
		SELECT id, name, create_time, update_time
		FROM courses
		WHERE id = ?
	`
	row := repo.DB.QueryRowContext(ctx, query, id)

	var course Course
	err := row.Scan(&course.ID, &course.Name, &course.CreateTime, &course.UpdateTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No course found
		}
		return nil, fmt.Errorf("failed to retrieve course: %w", err)
	}

	return &course, nil
}
