package enrollmentusecase

import (
	"context"
	common "github/rakadityas/course-management-system/common"
	courseDomain "github/rakadityas/course-management-system/domain/course"
	courseEnrollmentDomain "github/rakadityas/course-management-system/domain/course-enrollment"
	studentDomain "github/rakadityas/course-management-system/domain/student"
	"strconv"
)

// EnrollmentUseCaseInterface defines the interface for the EnrollmentUseCase.
type EnrollmentUseCaseItf interface {
	CourseSignUp(ctx context.Context, req CourseSignUpRequest) (CourseSignUpResp, error)
	ListCourses(ctx context.Context, studentID int64) (ListCoursesResp, error)
	CancelCourse(ctx context.Context, studentID, courseID int64) (CancelCourseResp, error)
	ListClassmates(ctx context.Context, studentID int64) (ListClassmatesResp, error)
}

type EnrollmentUseCase struct {
	studentService          studentDomain.StudentDomainItf
	courseService           courseDomain.CourseDomainItf
	courseEnrollmentService courseEnrollmentDomain.CourseEnrollmentDomainItf
}

func NewEnrollmentUseCase(studentService studentDomain.StudentDomainItf, courseService courseDomain.CourseDomainItf, courseEnrollmentService courseEnrollmentDomain.CourseEnrollmentDomainItf) EnrollmentUseCaseItf {
	return &EnrollmentUseCase{
		studentService:          studentService,
		courseService:           courseService,
		courseEnrollmentService: courseEnrollmentService,
	}
}

// CourseSignUp handles the course sign-up process.
func (enrollmentUC *EnrollmentUseCase) CourseSignUp(ctx context.Context, req CourseSignUpRequest) (CourseSignUpResp, error) {
	// Ensure the student data exists
	studentData, err := enrollmentUC.studentService.GetStudentByID(ctx, req.StudentID)
	if err != nil {
		return CourseSignUpResp{Status: common.StatusFailure, Message: "failed to retrieve student data"}, err
	}
	if studentData == nil {
		return CourseSignUpResp{Status: common.StatusFailure, Message: "student data not found"}, nil
	}

	// Ensure the course data exists
	courseData, err := enrollmentUC.courseService.GetCourseByID(ctx, req.CourseID)
	if err != nil {
		return CourseSignUpResp{Status: common.StatusFailure, Message: "failed to retrieve course data"}, err
	}
	if courseData == nil {
		return CourseSignUpResp{Status: common.StatusFailure, Message: "course data not found"}, nil
	}

	// Ensure the student never made any enrollment at all
	courseEnrollments, err := enrollmentUC.courseEnrollmentService.GetEnrollmentByStudentIDAndCourseID(ctx, req.StudentID, req.CourseID)
	if err != nil {
		return CourseSignUpResp{Status: common.StatusFailure, Message: "failed to retrieve course data"}, err
	}
	if len(courseEnrollments) > 0 {
		return CourseSignUpResp{Status: common.StatusFailure, Message: "student has enrolled before"}, err
	}

	// Create new enrollment
	newEnrollment, err := enrollmentUC.courseEnrollmentService.CreateEnrollment(ctx, req.StudentID, req.CourseID, courseEnrollmentDomain.StatusActive)
	if err != nil {
		return CourseSignUpResp{Status: common.StatusFailure, Message: "failed to sign up course"}, err
	}

	// Return a successful response
	return CourseSignUpResp{
		Status: common.StatusSuccess,
		EnrollmentData: &CourseEnrollment{
			ID:           newEnrollment.ID,
			StudentID:    newEnrollment.StudentID,
			StudentEmail: studentData.Email,
			CourseID:     newEnrollment.CourseID,
			CourseName:   courseData.Name,
			Status:       courseEnrollmentDomain.StatusActive,
			CreateTime:   newEnrollment.CreateTime,
			UpdateTime:   newEnrollment.UpdateTime,
		},
	}, nil
}

// ListCourses retrieves the list of courses a student is enrolled in.
func (enrollmentUC *EnrollmentUseCase) ListCourses(ctx context.Context, studentID int64) (ListCoursesResp, error) {
	// Ensure the student data exists
	studentData, err := enrollmentUC.studentService.GetStudentByID(ctx, studentID)
	if err != nil {
		return ListCoursesResp{Status: common.StatusFailure, Message: "failed to retrieve student data"}, err
	}
	if studentData == nil {
		return ListCoursesResp{Status: common.StatusFailure, Message: "student data not found"}, nil
	}

	// Get course enrollments for the student
	enrollments, err := enrollmentUC.courseEnrollmentService.GetEnrollmentByStudentID(ctx, studentID)
	if err != nil {
		return ListCoursesResp{Status: common.StatusFailure, Message: "failed to retrieve enrollments"}, err
	}

	// Prepare the response
	var courses []CourseDetail
	for _, enrollment := range enrollments {
		course, err := enrollmentUC.courseService.GetCourseByID(ctx, enrollment.CourseID)
		if err != nil {
			return ListCoursesResp{Status: common.StatusFailure, Message: "failed to retrieve course data"}, err
		}
		if course == nil {
			return ListCoursesResp{Status: common.StatusFailure, Message: "course data is not found for courseID: " + strconv.FormatInt(enrollment.CourseID, 10)}, nil
		}

		courses = append(courses, CourseDetail{
			CourseID:   course.ID,
			CourseName: course.Name,
			Status:     enrollment.Status,
			CreateTime: enrollment.CreateTime,
			UpdateTime: enrollment.UpdateTime,
		})
	}

	return ListCoursesResp{
		Status:  common.StatusSuccess,
		Courses: courses,
	}, nil
}

// CancelCourse cancel registered courses on the course enrollment table
func (enrollmentUC *EnrollmentUseCase) CancelCourse(ctx context.Context, studentID, courseID int64) (CancelCourseResp, error) {
	err := enrollmentUC.courseEnrollmentService.UpdateCourseEnrollmentStatus(ctx, studentID, courseID, courseEnrollmentDomain.StatusCancelled)
	if err != nil {
		return CancelCourseResp{Status: common.StatusFailure, Message: "failed to cancel course enrollment"}, err
	}

	return CancelCourseResp{
		Status: common.StatusSuccess,
	}, nil
}

// ListClassmates retrieves the list of classmates for the given student.
func (enrollmentUC *EnrollmentUseCase) ListClassmates(ctx context.Context, studentID int64) (ListClassmatesResp, error) {
	// Ensure the student data exists
	studentData, err := enrollmentUC.studentService.GetStudentByID(ctx, studentID)
	if err != nil {
		return ListClassmatesResp{Status: common.StatusFailure, Message: "failed to retrieve student data"}, err
	}
	if studentData == nil {
		return ListClassmatesResp{Status: common.StatusFailure, Message: "student data not found"}, nil
	}

	// Get course enrollments for the student
	enrollments, err := enrollmentUC.courseEnrollmentService.GetListClassmates(ctx, studentID)
	if err != nil {
		return ListClassmatesResp{Status: common.StatusFailure, Message: "failed to get list of classmates"}, err
	}

	// Create a map to group students by course ID
	mapCourseGroup := make(map[int64][]int64)
	for _, enrollment := range enrollments {
		mapCourseGroup[enrollment.CourseID] = append(mapCourseGroup[enrollment.CourseID], enrollment.StudentID)
	}

	// Prepare the response
	var response ListClassmatesResp
	for courseID, studentIDs := range mapCourseGroup {
		course, err := enrollmentUC.courseService.GetCourseByID(ctx, courseID) // todo: improve this with get bulk
		if err != nil {
			return ListClassmatesResp{Status: common.StatusFailure, Message: "failed to retrieve course data"}, err
		}
		if course == nil {
			return ListClassmatesResp{Status: common.StatusFailure, Message: "course data is not found for courseID: " + strconv.FormatInt(courseID, 10)}, nil
		}

		var classmates []ListClassmatesStudentsResp
		for _, id := range studentIDs {
			if id == studentID {
				continue // Skip the current student
			}

			student, err := enrollmentUC.studentService.GetStudentByID(ctx, id) // todo: improve this with get bulk
			if err != nil {
				return ListClassmatesResp{Status: common.StatusFailure, Message: "failed to retrieve student data: " + strconv.FormatInt(id, 10)}, err
			}
			if student == nil {
				return ListClassmatesResp{Status: common.StatusFailure, Message: "student data is not found for studentID: " + strconv.FormatInt(id, 10)}, nil
			}

			classmates = append(classmates, ListClassmatesStudentsResp{
				StudentID:    strconv.FormatInt(student.ID, 10),
				StudentEmail: student.Email,
			})
		}

		response.Courses = append(response.Courses, ListClassmatesCourseResp{
			CourseID:   course.ID,
			CourseName: course.Name,
			ClassMates: classmates,
		})
	}

	return ListClassmatesResp{
		Status:  common.StatusSuccess,
		Courses: response.Courses,
	}, nil
}
