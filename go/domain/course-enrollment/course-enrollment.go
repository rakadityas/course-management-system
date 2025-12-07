package courseenrollmentdomain

import (
	"context"
	"time"
)

type CourseEnrollmentDomainItf interface {
	CreateEnrollment(ctx context.Context, studentID, courseID int64, status int) (CourseEnrollment, error)
	GetEnrollmentByStudentID(ctx context.Context, studentID int64) ([]CourseEnrollment, error)
	GetEnrollmentByStudentIDAndCourseID(ctx context.Context, studentID, courseID int64) ([]CourseEnrollment, error)
	UpdateCourseEnrollmentStatus(ctx context.Context, studentID, courseID int64, newStatus int) error
	GetListClassmates(ctx context.Context, studentID int64) ([]CourseEnrollment, error)
}

type CourseEnrollmentService struct {
	repo CourseEnrollmentRepository
}

func NewCourseEnrollmentService(repo CourseEnrollmentRepository) CourseEnrollmentDomainItf {
	return &CourseEnrollmentService{repo: repo}
}

func (s *CourseEnrollmentService) CreateEnrollment(ctx context.Context, studentID, courseID int64, status int) (CourseEnrollment, error) {
	enrollment := CourseEnrollment{
		StudentID:  studentID,
		CourseID:   courseID,
		Status:     status,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	return s.repo.CreateEnrollment(ctx, enrollment)
}

func (s *CourseEnrollmentService) GetEnrollmentByStudentID(ctx context.Context, studentID int64) ([]CourseEnrollment, error) {
	return s.repo.GetEnrollmentByStudentID(ctx, studentID)
}

func (s *CourseEnrollmentService) UpdateCourseEnrollmentStatus(ctx context.Context, studentID, courseID int64, newStatus int) error {
	return s.repo.UpdateCourseEnrollmentStatus(ctx, studentID, courseID, newStatus)
}

func (s *CourseEnrollmentService) GetListClassmates(ctx context.Context, studentID int64) ([]CourseEnrollment, error) {
	return s.repo.GetListClassmates(ctx, studentID)
}

// GetEnrollmentByStudentIDAndCourseID retrieves course enrollments for a student and course.
func (service *CourseEnrollmentService) GetEnrollmentByStudentIDAndCourseID(ctx context.Context, studentID, courseID int64) ([]CourseEnrollment, error) {
	return service.repo.GetEnrollmentByStudentIDAndCourseID(ctx, studentID, courseID)
}
