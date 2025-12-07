package coursedomain

import "time"

type Course struct {
	ID         int64
	Name       string
	CreateTime time.Time
	UpdateTime time.Time
}

func NewCourse(name string) Course {
	return Course{
		Name: name,
	}
}
