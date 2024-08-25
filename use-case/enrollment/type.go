package enrollmentusecase

import "time"

// CourseSignUp related
type (
	CourseSignUpRequest struct {
		StudentID int64 `json:"student_id"`
		CourseID  int64 `json:"course_id"`
	}

	// CourseSignUpResp represents the response structure for course sign-up.
	CourseSignUpResp struct {
		Status         string            `json:"status"`
		Message        string            `json:"message,omitempty"`
		EnrollmentData *CourseEnrollment `json:"enrollment_data,omitempty"`
	}
)

// CourseEnrollment related
type (
	// CourseEnrollment represents the course enrollment details.
	CourseEnrollment struct {
		ID           int64     `json:"id"`
		StudentID    int64     `json:"student_id"`
		StudentEmail string    `json:"student_email"`
		CourseID     int64     `json:"course_id"`
		CourseName   string    `json:"course_name"`
		Status       int       `json:"status"`
		CreateTime   time.Time `json:"create_time"`
		UpdateTime   time.Time `json:"update_time"`
	}

	// ListCoursesResp represents the response structure for listing courses.
	ListCoursesResp struct {
		Status  string         `json:"status"`
		Message string         `json:"message,omitempty"`
		Courses []CourseDetail `json:"courses,omitempty"`
	}

	// CourseDetail provides detailed information about a course.
	CourseDetail struct {
		CourseID   int64     `json:"course_id"`
		CourseName string    `json:"course_name"`
		Status     int       `json:"status"`
		CreateTime time.Time `json:"create_time"`
		UpdateTime time.Time `json:"update_time"`
	}
)

// CancelCourseResp related
type (
	// CancelCourseRequest represents the request payload for canceling a course enrollment.
	CancelCourseRequest struct {
		StudentID int64 `json:"student_id"`
		CourseID  int64 `json:"course_id"`
	}

	// CancelCourseResp represents the response structure for course cancel
	CancelCourseResp struct {
		Status  string `json:"status"`
		Message string `json:"message,omitempty"`
	}
)

// ListClassmatesResp related
type (
	// ListClassmatesResp represents lists of students within the same course as the student
	ListClassmatesResp struct {
		Status  string                     `json:"status"`
		Message string                     `json:"message,omitempty"`
		Courses []ListClassmatesCourseResp `json:"courses"`
	}

	ListClassmatesCourseResp struct {
		CourseID   int64                        `json:"course_id"`
		CourseName string                       `json:"course_name"`
		ClassMates []ListClassmatesStudentsResp `json:"class_mates"`
	}

	ListClassmatesStudentsResp struct {
		StudentID    string `json:"student_id"`
		StudentEmail string `json:"student_email"`
	}
)
