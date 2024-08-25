package courseenrollmentdomain

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// Custom error for when no rows are updated.
var ErrNoRowsAffected = errors.New("no rows were updated")

type CourseEnrollmentRepository interface {
	CreateEnrollment(ctx context.Context, courseEnrollment CourseEnrollment) (CourseEnrollment, error)
	GetEnrollmentByStudentID(ctx context.Context, studentID int64) ([]CourseEnrollment, error)
	GetEnrollmentByStudentIDAndCourseID(ctx context.Context, studentID, courseID int64) ([]CourseEnrollment, error)
	UpdateCourseEnrollmentStatus(ctx context.Context, studentID, courseID int64, newStatus int) error
	GetListClassmates(ctx context.Context, studentID int64) ([]CourseEnrollment, error)
}

type CourseEnrollmentDB struct {
	DB *sql.DB
}

// NewSQLCourseRepository creates a new StudentDB instance with the given database connection.
func NewSQLCourseEnrollmentRepository(db *sql.DB) *CourseEnrollmentDB {
	return &CourseEnrollmentDB{DB: db}
}

// CreateEnrollment inserts a new course enrollment record into the database.
func (repo *CourseEnrollmentDB) CreateEnrollment(ctx context.Context, courseEnrollment CourseEnrollment) (CourseEnrollment, error) {
	query := `
		INSERT INTO course_enrollments (student_id, course_id, status, create_time, update_time)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := repo.DB.ExecContext(ctx, query, courseEnrollment.StudentID, courseEnrollment.CourseID, courseEnrollment.Status, courseEnrollment.CreateTime, courseEnrollment.UpdateTime)
	if err != nil {
		return CourseEnrollment{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return CourseEnrollment{}, err
	}

	courseEnrollment.ID = id
	return courseEnrollment, nil
}

// GetEnrollmentByStudentID retrieves all course enrollments for a given student.
func (repo *CourseEnrollmentDB) GetEnrollmentByStudentID(ctx context.Context, studentID int64) ([]CourseEnrollment, error) {
	query := `
		SELECT id, student_id, course_id, status, create_time, update_time
		FROM course_enrollments
		WHERE student_id = ? and status = 1
	`
	rows, err := repo.DB.QueryContext(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enrollments []CourseEnrollment
	for rows.Next() {
		var enrollment CourseEnrollment
		if err := rows.Scan(&enrollment.ID, &enrollment.StudentID, &enrollment.CourseID, &enrollment.Status, &enrollment.CreateTime, &enrollment.UpdateTime); err != nil {
			return nil, err
		}
		enrollments = append(enrollments, enrollment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return enrollments, nil
}

// UpdateCourseEnrollmentStatus updates the status of course enrollments for a specific student and course.
// Returns an error if no rows are affected by the update.
func (repo *CourseEnrollmentDB) UpdateCourseEnrollmentStatus(ctx context.Context, studentID, courseID int64, newStatus int) error {
	query := `
		UPDATE course_enrollments
		SET status = ?, update_time = ?
		WHERE student_id = ? AND course_id = ?
	`
	result, err := repo.DB.ExecContext(ctx, query, newStatus, time.Now(), studentID, courseID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNoRowsAffected
	}

	return nil
}

// GetListClassmates retrieves all students who have signed up for the same course as the specified student.
func (repo *CourseEnrollmentDB) GetListClassmates(ctx context.Context, studentID int64) ([]CourseEnrollment, error) {
	query := `
		SELECT ce.id, ce.student_id, ce.course_id, ce.status, ce.create_time, ce.update_time
		FROM course_enrollments ce
		JOIN course_enrollments ce2 ON ce.course_id = ce2.course_id
		WHERE ce2.student_id = ? AND ce.student_id != ? and ce2.status = 1 and ce.status = 1
	`
	rows, err := repo.DB.QueryContext(ctx, query, studentID, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classmates []CourseEnrollment
	for rows.Next() {
		var enrollment CourseEnrollment
		if err := rows.Scan(&enrollment.ID, &enrollment.StudentID, &enrollment.CourseID, &enrollment.Status, &enrollment.CreateTime, &enrollment.UpdateTime); err != nil {
			return nil, err
		}
		classmates = append(classmates, enrollment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return classmates, nil
}

// GetEnrollmentByStudentIDAndCourseID retrieves enrollments for a student and course.
func (repo *CourseEnrollmentDB) GetEnrollmentByStudentIDAndCourseID(ctx context.Context, studentID, courseID int64) ([]CourseEnrollment, error) {
	query := `
		SELECT id, student_id, course_id, status, create_time, update_time
		FROM course_enrollments
		WHERE student_id = ? AND course_id = ?
	`

	rows, err := repo.DB.QueryContext(ctx, query, studentID, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enrollments []CourseEnrollment
	for rows.Next() {
		var enrollment CourseEnrollment
		if err := rows.Scan(&enrollment.ID, &enrollment.StudentID, &enrollment.CourseID, &enrollment.Status, &enrollment.CreateTime, &enrollment.UpdateTime); err != nil {
			return nil, err
		}
		enrollments = append(enrollments, enrollment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return enrollments, nil
}
