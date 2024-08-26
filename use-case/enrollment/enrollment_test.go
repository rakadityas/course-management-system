package enrollmentusecase

import (
	"context"
	"errors"
	common "github/rakadityas/course-management-system/common"
	courseDomain "github/rakadityas/course-management-system/domain/course"
	courseEnrollmentDomain "github/rakadityas/course-management-system/domain/course-enrollment"
	courseEnrollmentDomainMock "github/rakadityas/course-management-system/domain/course-enrollment/mocks"
	courseDomainMock "github/rakadityas/course-management-system/domain/course/mocks"
	studentDomain "github/rakadityas/course-management-system/domain/student"
	studentDomainMock "github/rakadityas/course-management-system/domain/student/mocks"
	"reflect"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
)

func TestEnrollmentUseCase_CourseSignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		studentID int64 = 1
		courseID  int64 = 101
		status          = courseEnrollmentDomain.StatusActive
	)
	constCreateTime := time.Date(2023, 8, 25, 0, 0, 0, 0, time.UTC)
	constUpdateTime := time.Date(2023, 8, 25, 1, 0, 0, 0, time.UTC)

	type fields struct {
		studentService          studentDomain.StudentDomainItf
		courseService           courseDomain.CourseDomainItf
		courseEnrollmentService courseEnrollmentDomain.CourseEnrollmentDomainItf
	}
	type args struct {
		ctx context.Context
		req CourseSignUpRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    CourseSignUpResp
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				studentService: func() studentDomain.StudentDomainItf {
					mock := studentDomainMock.NewMockStudentDomainItf(ctrl)
					mock.EXPECT().GetStudentByID(gomock.Any(), studentID).Return(&studentDomain.Student{ID: studentID, Email: "student@example.com"}, nil)
					return mock
				}(),
				courseService: func() courseDomain.CourseDomainItf {
					mock := courseDomainMock.NewMockCourseDomainItf(ctrl)
					mock.EXPECT().GetCourseByID(gomock.Any(), courseID).Return(&courseDomain.Course{ID: courseID, Name: "Course Name"}, nil)
					return mock
				}(),
				courseEnrollmentService: func() courseEnrollmentDomain.CourseEnrollmentDomainItf {
					mock := courseEnrollmentDomainMock.NewMockCourseEnrollmentDomainItf(ctrl)
					mock.EXPECT().GetEnrollmentByStudentIDAndCourseID(gomock.Any(), studentID, courseID).Return([]courseEnrollmentDomain.CourseEnrollment{}, nil)
					mock.EXPECT().CreateEnrollment(gomock.Any(), studentID, courseID, status).Return(courseEnrollmentDomain.CourseEnrollment{
						ID:         1,
						StudentID:  studentID,
						CourseID:   courseID,
						Status:     status,
						CreateTime: constCreateTime,
						UpdateTime: constUpdateTime,
					}, nil)
					return mock
				}(),
			},
			args: args{
				ctx: context.Background(),
				req: CourseSignUpRequest{
					StudentID: studentID,
					CourseID:  courseID,
				},
			},
			want: CourseSignUpResp{
				Status: common.StatusSuccess,
				EnrollmentData: &CourseEnrollment{
					ID:           1,
					StudentID:    studentID,
					StudentEmail: "student@example.com",
					CourseID:     courseID,
					CourseName:   "Course Name",
					Status:       status,
					CreateTime:   constCreateTime,
					UpdateTime:   constUpdateTime,
				},
			},
			wantErr: false,
		},
		{
			name: "Student Not Found",
			fields: fields{
				studentService: func() studentDomain.StudentDomainItf {
					mock := studentDomainMock.NewMockStudentDomainItf(ctrl)
					mock.EXPECT().GetStudentByID(gomock.Any(), studentID).Return(nil, nil)
					return mock
				}(),
				courseService: func() courseDomain.CourseDomainItf {
					return courseDomainMock.NewMockCourseDomainItf(ctrl)
				}(),
				courseEnrollmentService: func() courseEnrollmentDomain.CourseEnrollmentDomainItf {
					return courseEnrollmentDomainMock.NewMockCourseEnrollmentDomainItf(ctrl)
				}(),
			},
			args: args{
				ctx: context.Background(),
				req: CourseSignUpRequest{
					StudentID: studentID,
					CourseID:  courseID,
				},
			},
			want: CourseSignUpResp{
				Status:  common.StatusFailure,
				Message: "student data not found",
			},
			wantErr: false,
		},
		{
			name: "Course Not Found",
			fields: fields{
				studentService: func() studentDomain.StudentDomainItf {
					mock := studentDomainMock.NewMockStudentDomainItf(ctrl)
					mock.EXPECT().GetStudentByID(gomock.Any(), studentID).Return(&studentDomain.Student{ID: studentID, Email: "student@example.com"}, nil)
					return mock
				}(),
				courseService: func() courseDomain.CourseDomainItf {
					mock := courseDomainMock.NewMockCourseDomainItf(ctrl)
					mock.EXPECT().GetCourseByID(gomock.Any(), courseID).Return(nil, nil)
					return mock
				}(),
				courseEnrollmentService: func() courseEnrollmentDomain.CourseEnrollmentDomainItf {
					return courseEnrollmentDomainMock.NewMockCourseEnrollmentDomainItf(ctrl)
				}(),
			},
			args: args{
				ctx: context.Background(),
				req: CourseSignUpRequest{
					StudentID: studentID,
					CourseID:  courseID,
				},
			},
			want: CourseSignUpResp{
				Status:  common.StatusFailure,
				Message: "course data not found",
			},
			wantErr: false,
		},
		{
			name: "Enrollment Already Exists",
			fields: fields{
				studentService: func() studentDomain.StudentDomainItf {
					mock := studentDomainMock.NewMockStudentDomainItf(ctrl)
					mock.EXPECT().GetStudentByID(gomock.Any(), studentID).Return(&studentDomain.Student{ID: studentID, Email: "student@example.com"}, nil)
					return mock
				}(),
				courseService: func() courseDomain.CourseDomainItf {
					mock := courseDomainMock.NewMockCourseDomainItf(ctrl)
					mock.EXPECT().GetCourseByID(gomock.Any(), courseID).Return(&courseDomain.Course{ID: courseID, Name: "Course Name"}, nil)
					return mock
				}(),
				courseEnrollmentService: func() courseEnrollmentDomain.CourseEnrollmentDomainItf {
					mock := courseEnrollmentDomainMock.NewMockCourseEnrollmentDomainItf(ctrl)

					mock.EXPECT().GetEnrollmentByStudentIDAndCourseID(gomock.Any(), studentID, courseID).Return([]courseEnrollmentDomain.CourseEnrollment{
						{ID: 1, StudentID: studentID, CourseID: courseID, Status: status}}, nil)
					return mock
				}(),
			},
			args: args{
				ctx: context.Background(),
				req: CourseSignUpRequest{
					StudentID: studentID,
					CourseID:  courseID,
				},
			},
			want: CourseSignUpResp{
				Status:  common.StatusFailure,
				Message: "student has enrolled before",
			},
			wantErr: false,
		},
		{
			name: "Create Enrollment Error",
			fields: fields{
				studentService: func() studentDomain.StudentDomainItf {
					mock := studentDomainMock.NewMockStudentDomainItf(ctrl)
					mock.EXPECT().GetStudentByID(gomock.Any(), studentID).Return(&studentDomain.Student{ID: studentID, Email: "student@example.com"}, nil)
					return mock
				}(),
				courseService: func() courseDomain.CourseDomainItf {
					mock := courseDomainMock.NewMockCourseDomainItf(ctrl)
					mock.EXPECT().GetCourseByID(gomock.Any(), courseID).Return(&courseDomain.Course{ID: courseID, Name: "Course Name"}, nil)
					return mock
				}(),
				courseEnrollmentService: func() courseEnrollmentDomain.CourseEnrollmentDomainItf {
					mock := courseEnrollmentDomainMock.NewMockCourseEnrollmentDomainItf(ctrl)
					mock.EXPECT().GetEnrollmentByStudentIDAndCourseID(gomock.Any(), studentID, courseID).Return([]courseEnrollmentDomain.CourseEnrollment{}, nil)
					mock.EXPECT().CreateEnrollment(gomock.Any(), studentID, courseID, status).Return(courseEnrollmentDomain.CourseEnrollment{}, errors.New("create enrollment error"))
					return mock
				}(),
			},
			args: args{
				ctx: context.Background(),
				req: CourseSignUpRequest{
					StudentID: studentID,
					CourseID:  courseID,
				},
			},
			want: CourseSignUpResp{
				Status:  common.StatusFailure,
				Message: "failed to sign up course",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enrollmentUC := &EnrollmentUseCase{
				studentService:          tt.fields.studentService,
				courseService:           tt.fields.courseService,
				courseEnrollmentService: tt.fields.courseEnrollmentService,
			}
			got, err := enrollmentUC.CourseSignUp(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnrollmentUseCase.CourseSignUp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EnrollmentUseCase.CourseSignUp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnrollmentUseCase_ListCourses(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const studentID int64 = 1
	timestamp := time.Date(2023, 8, 25, 0, 0, 0, 0, time.UTC)

	type fields struct {
		studentService          studentDomain.StudentDomainItf
		courseService           courseDomain.CourseDomainItf
		courseEnrollmentService courseEnrollmentDomain.CourseEnrollmentDomainItf
	}
	type args struct {
		ctx       context.Context
		studentID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ListCoursesResp
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				studentService: func() studentDomain.StudentDomainItf {
					mock := studentDomainMock.NewMockStudentDomainItf(ctrl)
					mock.EXPECT().GetStudentByID(gomock.Any(), studentID).Return(&studentDomain.Student{ID: studentID, Email: "student@example.com"}, nil)
					return mock
				}(),
				courseService: func() courseDomain.CourseDomainItf {
					mock := courseDomainMock.NewMockCourseDomainItf(ctrl)
					mock.EXPECT().GetCourseByID(gomock.Any(), int64(101)).Return(&courseDomain.Course{ID: int64(101), Name: "Course Name"}, nil)
					return mock
				}(),
				courseEnrollmentService: func() courseEnrollmentDomain.CourseEnrollmentDomainItf {
					mock := courseEnrollmentDomainMock.NewMockCourseEnrollmentDomainItf(ctrl)
					mock.EXPECT().GetEnrollmentByStudentID(gomock.Any(), studentID).Return([]courseEnrollmentDomain.CourseEnrollment{
						{
							ID:         1,
							StudentID:  studentID,
							CourseID:   int64(101),
							Status:     courseEnrollmentDomain.StatusActive,
							CreateTime: timestamp,
							UpdateTime: timestamp,
						},
					}, nil)
					return mock
				}(),
			},
			args: args{
				ctx:       context.Background(),
				studentID: studentID,
			},
			want: ListCoursesResp{
				Status: common.StatusSuccess,
				Courses: []CourseDetail{
					{
						CourseID:   int64(101),
						CourseName: "Course Name",
						Status:     courseEnrollmentDomain.StatusActive,
						CreateTime: timestamp,
						UpdateTime: timestamp,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Failed to Retrieve Enrollments",
			fields: fields{
				studentService: func() studentDomain.StudentDomainItf {
					mock := studentDomainMock.NewMockStudentDomainItf(ctrl)
					mock.EXPECT().GetStudentByID(gomock.Any(), studentID).Return(&studentDomain.Student{ID: studentID, Email: "student@example.com"}, nil)
					return mock
				}(),
				courseEnrollmentService: func() courseEnrollmentDomain.CourseEnrollmentDomainItf {
					mock := courseEnrollmentDomainMock.NewMockCourseEnrollmentDomainItf(ctrl)
					mock.EXPECT().GetEnrollmentByStudentID(gomock.Any(), studentID).Return(nil, errors.New("enrollments error"))
					return mock
				}(),
				courseService: func() courseDomain.CourseDomainItf {
					return courseDomainMock.NewMockCourseDomainItf(ctrl)
				}(),
			},
			args: args{
				ctx:       context.Background(),
				studentID: studentID,
			},
			want: ListCoursesResp{
				Status:  common.StatusFailure,
				Message: "failed to retrieve enrollments",
			},
			wantErr: true,
		},
		{
			name: "Failed to Retrieve Course Data",
			fields: fields{
				studentService: func() studentDomain.StudentDomainItf {
					mock := studentDomainMock.NewMockStudentDomainItf(ctrl)
					mock.EXPECT().GetStudentByID(gomock.Any(), studentID).Return(&studentDomain.Student{ID: studentID, Email: "student@example.com"}, nil)
					return mock
				}(),
				courseEnrollmentService: func() courseEnrollmentDomain.CourseEnrollmentDomainItf {
					mock := courseEnrollmentDomainMock.NewMockCourseEnrollmentDomainItf(ctrl)
					mock.EXPECT().GetEnrollmentByStudentID(gomock.Any(), studentID).Return([]courseEnrollmentDomain.CourseEnrollment{
						{
							ID:         1,
							StudentID:  studentID,
							CourseID:   int64(101),
							Status:     courseEnrollmentDomain.StatusActive,
							CreateTime: timestamp,
							UpdateTime: timestamp,
						},
					}, nil)
					return mock
				}(),
				courseService: func() courseDomain.CourseDomainItf {
					mock := courseDomainMock.NewMockCourseDomainItf(ctrl)
					mock.EXPECT().GetCourseByID(gomock.Any(), int64(101)).Return(nil, errors.New("course data error"))
					return mock
				}(),
			},
			args: args{
				ctx:       context.Background(),
				studentID: studentID,
			},
			want: ListCoursesResp{
				Status:  common.StatusFailure,
				Message: "failed to retrieve course data",
			},
			wantErr: true,
		},
		{
			name: "Course Data Not Found",
			fields: fields{
				studentService: func() studentDomain.StudentDomainItf {
					mock := studentDomainMock.NewMockStudentDomainItf(ctrl)
					mock.EXPECT().GetStudentByID(gomock.Any(), studentID).Return(&studentDomain.Student{ID: studentID, Email: "student@example.com"}, nil)
					return mock
				}(),
				courseEnrollmentService: func() courseEnrollmentDomain.CourseEnrollmentDomainItf {
					mock := courseEnrollmentDomainMock.NewMockCourseEnrollmentDomainItf(ctrl)
					mock.EXPECT().GetEnrollmentByStudentID(gomock.Any(), studentID).Return([]courseEnrollmentDomain.CourseEnrollment{
						{
							ID:         1,
							StudentID:  studentID,
							CourseID:   int64(101),
							Status:     courseEnrollmentDomain.StatusActive,
							CreateTime: timestamp,
							UpdateTime: timestamp,
						},
					}, nil)
					return mock
				}(),
				courseService: func() courseDomain.CourseDomainItf {
					mock := courseDomainMock.NewMockCourseDomainItf(ctrl)
					mock.EXPECT().GetCourseByID(gomock.Any(), int64(101)).Return(nil, nil)
					return mock
				}(),
			},
			args: args{
				ctx:       context.Background(),
				studentID: studentID,
			},
			want: ListCoursesResp{
				Status:  common.StatusFailure,
				Message: "course data is not found for courseID: 101",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enrollmentUC := &EnrollmentUseCase{
				studentService:          tt.fields.studentService,
				courseService:           tt.fields.courseService,
				courseEnrollmentService: tt.fields.courseEnrollmentService,
			}
			got, err := enrollmentUC.ListCourses(tt.args.ctx, tt.args.studentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnrollmentUseCase.ListCourses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EnrollmentUseCase.ListCourses() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnrollmentUseCase_CancelCourse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const studentID int64 = 1
	const courseID int64 = 101

	type fields struct {
		studentService          studentDomain.StudentDomainItf
		courseService           courseDomain.CourseDomainItf
		courseEnrollmentService courseEnrollmentDomain.CourseEnrollmentDomainItf
	}
	type args struct {
		ctx       context.Context
		studentID int64
		courseID  int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    CancelCourseResp
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				courseEnrollmentService: func() courseEnrollmentDomain.CourseEnrollmentDomainItf {
					mock := courseEnrollmentDomainMock.NewMockCourseEnrollmentDomainItf(ctrl)
					mock.EXPECT().UpdateCourseEnrollmentStatus(gomock.Any(), studentID, courseID, courseEnrollmentDomain.StatusCancelled).Return(nil)
					return mock
				}(),
				studentService: func() studentDomain.StudentDomainItf {
					return studentDomainMock.NewMockStudentDomainItf(ctrl)
				}(),
				courseService: func() courseDomain.CourseDomainItf {
					return courseDomainMock.NewMockCourseDomainItf(ctrl)
				}(),
			},
			args: args{
				ctx:       context.Background(),
				studentID: studentID,
				courseID:  courseID,
			},
			want: CancelCourseResp{
				Status: common.StatusSuccess,
			},
			wantErr: false,
		},
		{
			name: "Failed to Cancel Course Enrollment",
			fields: fields{
				courseEnrollmentService: func() courseEnrollmentDomain.CourseEnrollmentDomainItf {
					mock := courseEnrollmentDomainMock.NewMockCourseEnrollmentDomainItf(ctrl)
					mock.EXPECT().UpdateCourseEnrollmentStatus(gomock.Any(), studentID, courseID, courseEnrollmentDomain.StatusCancelled).Return(errors.New("update error"))
					return mock
				}(),
				studentService: func() studentDomain.StudentDomainItf {
					return studentDomainMock.NewMockStudentDomainItf(ctrl)
				}(),
				courseService: func() courseDomain.CourseDomainItf {
					return courseDomainMock.NewMockCourseDomainItf(ctrl)
				}(),
			},
			args: args{
				ctx:       context.Background(),
				studentID: studentID,
				courseID:  courseID,
			},
			want: CancelCourseResp{
				Status:  common.StatusFailure,
				Message: "failed to cancel course enrollment",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enrollmentUC := &EnrollmentUseCase{
				studentService:          tt.fields.studentService,
				courseService:           tt.fields.courseService,
				courseEnrollmentService: tt.fields.courseEnrollmentService,
			}
			got, err := enrollmentUC.CancelCourse(tt.args.ctx, tt.args.studentID, tt.args.courseID)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnrollmentUseCase.CancelCourse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EnrollmentUseCase.CancelCourse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnrollmentUseCase_ListClassmates(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const studentID int64 = 1

	type fields struct {
		studentService          studentDomain.StudentDomainItf
		courseService           courseDomain.CourseDomainItf
		courseEnrollmentService courseEnrollmentDomain.CourseEnrollmentDomainItf
	}
	type args struct {
		ctx       context.Context
		studentID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ListClassmatesResp
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				courseEnrollmentService: func() courseEnrollmentDomain.CourseEnrollmentDomainItf {
					mock := courseEnrollmentDomainMock.NewMockCourseEnrollmentDomainItf(ctrl)
					mock.EXPECT().GetListClassmates(gomock.Any(), studentID).Return([]courseEnrollmentDomain.CourseEnrollment{
						{CourseID: 101, StudentID: 2},
						{CourseID: 101, StudentID: 3},
					}, nil)
					return mock
				}(),
				courseService: func() courseDomain.CourseDomainItf {
					mock := courseDomainMock.NewMockCourseDomainItf(ctrl)
					mock.EXPECT().GetCourseByID(gomock.Any(), int64(101)).Return(&courseDomain.Course{ID: 101, Name: "Course A"}, nil)
					return mock
				}(),
				studentService: func() studentDomain.StudentDomainItf {
					mock := studentDomainMock.NewMockStudentDomainItf(ctrl)
					mock.EXPECT().GetStudentByID(gomock.Any(), int64(1)).Return(&studentDomain.Student{ID: 1, Email: "student1@example.com"}, nil)
					mock.EXPECT().GetStudentByID(gomock.Any(), int64(2)).Return(&studentDomain.Student{ID: 2, Email: "student2@example.com"}, nil)
					mock.EXPECT().GetStudentByID(gomock.Any(), int64(3)).Return(&studentDomain.Student{ID: 3, Email: "student3@example.com"}, nil)
					return mock
				}(),
			},
			args: args{
				ctx:       context.Background(),
				studentID: studentID,
			},
			want: ListClassmatesResp{
				Status: common.StatusSuccess,
				Courses: []ListClassmatesCourseResp{
					{
						CourseID:   101,
						CourseName: "Course A",
						ClassMates: []ListClassmatesStudentsResp{
							{StudentID: "2", StudentEmail: "student2@example.com"},
							{StudentID: "3", StudentEmail: "student3@example.com"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Failed to Retrieve Enrollments",
			fields: fields{
				courseEnrollmentService: func() courseEnrollmentDomain.CourseEnrollmentDomainItf {
					mock := courseEnrollmentDomainMock.NewMockCourseEnrollmentDomainItf(ctrl)
					mock.EXPECT().GetListClassmates(gomock.Any(), studentID).Return(nil, errors.New("fetch enrollments error"))
					return mock
				}(),
				courseService: func() courseDomain.CourseDomainItf {
					return courseDomainMock.NewMockCourseDomainItf(ctrl)
				}(),
				studentService: func() studentDomain.StudentDomainItf {
					mock := studentDomainMock.NewMockStudentDomainItf(ctrl)
					mock.EXPECT().GetStudentByID(gomock.Any(), int64(1)).Return(&studentDomain.Student{ID: 1, Email: "student1@example.com"}, nil)
					return mock
				}(),
			},
			args: args{
				ctx:       context.Background(),
				studentID: studentID,
			},
			want: ListClassmatesResp{
				Status:  common.StatusFailure,
				Message: "failed to get list of classmates",
			},
			wantErr: true,
		},
		{
			name: "Failed to Retrieve Course Data",
			fields: fields{
				courseEnrollmentService: func() courseEnrollmentDomain.CourseEnrollmentDomainItf {
					mock := courseEnrollmentDomainMock.NewMockCourseEnrollmentDomainItf(ctrl)
					mock.EXPECT().GetListClassmates(gomock.Any(), studentID).Return([]courseEnrollmentDomain.CourseEnrollment{
						{CourseID: 101, StudentID: 2},
					}, nil)
					return mock
				}(),
				courseService: func() courseDomain.CourseDomainItf {
					mock := courseDomainMock.NewMockCourseDomainItf(ctrl)
					mock.EXPECT().GetCourseByID(gomock.Any(), int64(101)).Return(nil, errors.New("course not found"))
					return mock
				}(),
				studentService: func() studentDomain.StudentDomainItf {
					mock := studentDomainMock.NewMockStudentDomainItf(ctrl)
					mock.EXPECT().GetStudentByID(gomock.Any(), int64(1)).Return(&studentDomain.Student{ID: 1, Email: "student1@example.com"}, nil)
					return mock
				}(),
			},
			args: args{
				ctx:       context.Background(),
				studentID: studentID,
			},
			want: ListClassmatesResp{
				Status:  common.StatusFailure,
				Message: "failed to retrieve course data",
			},
			wantErr: true,
		},
		{
			name: "Failed to Retrieve Student Data",
			fields: fields{
				courseEnrollmentService: func() courseEnrollmentDomain.CourseEnrollmentDomainItf {
					mock := courseEnrollmentDomainMock.NewMockCourseEnrollmentDomainItf(ctrl)
					mock.EXPECT().GetListClassmates(gomock.Any(), studentID).Return([]courseEnrollmentDomain.CourseEnrollment{
						{CourseID: 101, StudentID: 2},
					}, nil)
					return mock
				}(),
				courseService: func() courseDomain.CourseDomainItf {
					mock := courseDomainMock.NewMockCourseDomainItf(ctrl)
					mock.EXPECT().GetCourseByID(gomock.Any(), int64(101)).Return(&courseDomain.Course{ID: 101, Name: "Course A"}, nil)
					return mock
				}(),
				studentService: func() studentDomain.StudentDomainItf {
					mock := studentDomainMock.NewMockStudentDomainItf(ctrl)
					mock.EXPECT().GetStudentByID(gomock.Any(), int64(1)).Return(&studentDomain.Student{ID: 1, Email: "student1@example.com"}, nil)
					mock.EXPECT().GetStudentByID(gomock.Any(), int64(2)).Return(nil, errors.New("student not found"))
					return mock
				}(),
			},
			args: args{
				ctx:       context.Background(),
				studentID: studentID,
			},
			want: ListClassmatesResp{
				Status:  common.StatusFailure,
				Message: "failed to retrieve student data: 2",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enrollmentUC := &EnrollmentUseCase{
				studentService:          tt.fields.studentService,
				courseService:           tt.fields.courseService,
				courseEnrollmentService: tt.fields.courseEnrollmentService,
			}
			got, err := enrollmentUC.ListClassmates(tt.args.ctx, tt.args.studentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnrollmentUseCase.ListClassmates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EnrollmentUseCase.ListClassmates() = %v, want %v", got, tt.want)
			}
		})
	}
}
