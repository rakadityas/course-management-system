package studentdomain

import "context"

type StudentDomainItf interface {
	GetStudentByID(ctx context.Context, studentID int64) (*Student, error)
}

type StudentService struct {
	repo StudentRepository
}

func NewStudentService(repo StudentRepository) StudentDomainItf {
	return &StudentService{repo: repo}
}

// GetStudentByID retrieves a student by their ID.
func (s *StudentService) GetStudentByID(ctx context.Context, id int64) (*Student, error) {
	return s.repo.GetStudentByID(ctx, id)
}
