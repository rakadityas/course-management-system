package courseenrollmentdomain

import "time"

type CourseEnrollment struct {
	ID         int64
	StudentID  int64
	CourseID   int64
	Status     int
	CreateTime time.Time
	UpdateTime time.Time
}

func NewCourseEnrollment(studentID, courseID int64, status int) CourseEnrollment {
	return CourseEnrollment{
		StudentID: studentID,
		CourseID:  courseID,
		Status:    status,
	}
}
