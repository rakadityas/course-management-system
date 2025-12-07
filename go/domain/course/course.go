package coursedomain

import "context"

type CourseDomainItf interface {
	GetCourseByID(ctx context.Context, id int64) (*Course, error)
}

type CourseService struct {
	repo CourseRepository
}

func NewCourseService(repo CourseRepository) CourseDomainItf {
	return &CourseService{repo: repo}
}

func (s *CourseService) GetCourseByID(ctx context.Context, id int64) (*Course, error) {
	return s.repo.GetCourseByID(ctx, id)
}
