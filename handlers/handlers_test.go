package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	enrollmentUseCase "github/rakadityas/course-management-system/use-case/enrollment"
	enrollmentUseCaseMock "github/rakadityas/course-management-system/use-case/enrollment/mocks"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestHandler_CourseSignUpHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		studentID int64 = 1
		courseID  int64 = 101
		status          = 1
	)
	type fields struct {
		EnrollmentUseCase enrollmentUseCase.EnrollmentUseCaseItf
	}
	tests := []struct {
		name           string
		fields         fields
		requestPayload enrollmentUseCase.CourseSignUpRequest
		mockResp       enrollmentUseCase.CourseSignUpResp
		mockErr        error
		wantStatusCode int
		wantBody       string
	}{
		{
			name: "Success",
			fields: fields{
				EnrollmentUseCase: func() enrollmentUseCase.EnrollmentUseCaseItf {
					mockEnrollmentUC := enrollmentUseCaseMock.NewMockEnrollmentUseCaseItf(ctrl)
					mockEnrollmentUC.EXPECT().CourseSignUp(gomock.Any(), enrollmentUseCase.CourseSignUpRequest{StudentID: studentID, CourseID: courseID}).Return(enrollmentUseCase.CourseSignUpResp{
						Status: "success",
						EnrollmentData: &enrollmentUseCase.CourseEnrollment{
							ID:           1,
							StudentID:    studentID,
							StudentEmail: "student@example.com",
							CourseID:     courseID,
							CourseName:   "Course Name",
							Status:       status,
						},
					}, nil)
					return mockEnrollmentUC
				}(),
			},
			requestPayload: enrollmentUseCase.CourseSignUpRequest{
				StudentID: studentID,
				CourseID:  courseID,
			},
			wantStatusCode: http.StatusOK,
			wantBody:       `{"status":"success","enrollment_data":{"id":1,"student_id":1,"student_email":"student@example.com","course_id":101,"course_name":"Course Name","status":1,"create_time":"0001-01-01T00:00:00Z","update_time":"0001-01-01T00:00:00Z"}}`,
		},
		{
			name: "Empty Request Data",
			fields: fields{
				EnrollmentUseCase: nil,
			},
			requestPayload: enrollmentUseCase.CourseSignUpRequest{
				StudentID: 0,
				CourseID:  courseID,
			},
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "Request Data is empty\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				EnrollmentUseCase: tt.fields.EnrollmentUseCase,
			}

			body, _ := json.Marshal(tt.requestPayload)
			req := httptest.NewRequest(http.MethodPost, "/course-sign-up", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			handler := h.CourseSignUpHandler()
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Errorf("Status code = %v, want %v", rec.Code, tt.wantStatusCode)
			}

			if rec.Code == http.StatusOK {
				var gotBody, wantBody map[string]interface{}
				if err := json.Unmarshal(rec.Body.Bytes(), &gotBody); err != nil {
					t.Fatalf("Failed to unmarshal response body: %v", err)
				}
				if err := json.Unmarshal([]byte(tt.wantBody), &wantBody); err != nil {
					t.Fatalf("Failed to unmarshal expected body: %v", err)
				}
				if !reflect.DeepEqual(gotBody, wantBody) {
					t.Errorf("Response body = %v, want %v", gotBody, wantBody)
				}
			} else {
				if rec.Body.String() != tt.wantBody {
					t.Errorf("Response body = %v, want %v", rec.Body.String(), tt.wantBody)
				}
			}
		})
	}
}

func TestHandler_ListCoursesHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const studentID int64 = 1
	type fields struct {
		EnrollmentUseCase enrollmentUseCase.EnrollmentUseCaseItf
	}
	tests := []struct {
		name           string
		fields         fields
		queryParams    map[string]string
		mockResp       enrollmentUseCase.ListCoursesResp
		mockErr        error
		wantStatusCode int
		wantBody       string
	}{
		{
			name: "Success",
			fields: fields{
				EnrollmentUseCase: func() enrollmentUseCase.EnrollmentUseCaseItf {
					mockEnrollmentUC := enrollmentUseCaseMock.NewMockEnrollmentUseCaseItf(ctrl)
					mockEnrollmentUC.EXPECT().ListCourses(gomock.Any(), studentID).Return(enrollmentUseCase.ListCoursesResp{
						Status: "success",
						Courses: []enrollmentUseCase.CourseDetail{
							{
								CourseID:   101,
								CourseName: "Course A",
								Status:     1,
								CreateTime: time.Time{}, // Use a fixed time or time.Now() in actual tests
								UpdateTime: time.Time{}, // Use a fixed time or time.Now() in actual tests
							},
						},
					}, nil)
					return mockEnrollmentUC
				}(),
			},
			queryParams: map[string]string{
				"student_id": strconv.FormatInt(studentID, 10),
			},
			wantStatusCode: http.StatusOK,
			wantBody:       `{"status":"success","courses":[{"course_id":101,"course_name":"Course A","status":1,"create_time":"0001-01-01T00:00:00Z","update_time":"0001-01-01T00:00:00Z"}]}`,
		},
		{
			name: "Invalid Student ID",
			fields: fields{
				EnrollmentUseCase: nil,
			},
			queryParams: map[string]string{
				"student_id": "invalid",
			},
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "Invalid student ID\n",
		},
		{
			name: "Zero Student ID",
			fields: fields{
				EnrollmentUseCase: nil,
			},
			queryParams: map[string]string{
				"student_id": "0",
			},
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "Student ID Zero\n",
		},
		{
			name: "Error From UseCase",
			fields: fields{
				EnrollmentUseCase: func() enrollmentUseCase.EnrollmentUseCaseItf {
					mockEnrollmentUC := enrollmentUseCaseMock.NewMockEnrollmentUseCaseItf(ctrl)
					mockEnrollmentUC.EXPECT().ListCourses(gomock.Any(), studentID).Return(enrollmentUseCase.ListCoursesResp{
						Status:  "error",
						Message: "failed to retrieve courses",
					}, errors.New("some error"))
					return mockEnrollmentUC
				}(),
			},
			queryParams: map[string]string{
				"student_id": strconv.FormatInt(studentID, 10),
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       "failed to retrieve courses\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				EnrollmentUseCase: tt.fields.EnrollmentUseCase,
			}

			req := httptest.NewRequest(http.MethodGet, "/list-courses", nil)
			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()
			rec := httptest.NewRecorder()

			handler := h.ListCoursesHandler()
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Errorf("Status code = %v, want %v", rec.Code, tt.wantStatusCode)
			}

			if rec.Code == http.StatusOK {
				var gotBody, wantBody map[string]interface{}
				if err := json.Unmarshal(rec.Body.Bytes(), &gotBody); err != nil {
					t.Fatalf("Failed to unmarshal response body: %v", err)
				}
				if err := json.Unmarshal([]byte(tt.wantBody), &wantBody); err != nil {
					t.Fatalf("Failed to unmarshal expected body: %v", err)
				}
				if !reflect.DeepEqual(gotBody, wantBody) {
					t.Errorf("Response body = %v, want %v", gotBody, wantBody)
				}
			} else {
				if rec.Body.String() != tt.wantBody {
					t.Errorf("Response body = %v, want %v", rec.Body.String(), tt.wantBody)
				}
			}
		})
	}
}

func TestHandler_CancelCourseHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		studentID int64 = 1
		courseID  int64 = 101
	)
	type fields struct {
		EnrollmentUseCase enrollmentUseCase.EnrollmentUseCaseItf
	}
	tests := []struct {
		name           string
		fields         fields
		requestPayload enrollmentUseCase.CancelCourseRequest
		mockResp       enrollmentUseCase.CancelCourseResp
		mockErr        error
		wantStatusCode int
		wantBody       string
	}{
		{
			name: "Success",
			fields: fields{
				EnrollmentUseCase: func() enrollmentUseCase.EnrollmentUseCaseItf {
					mockEnrollmentUC := enrollmentUseCaseMock.NewMockEnrollmentUseCaseItf(ctrl)
					mockEnrollmentUC.EXPECT().CancelCourse(gomock.Any(), studentID, courseID).Return(enrollmentUseCase.CancelCourseResp{
						Status: "success",
					}, nil)
					return mockEnrollmentUC
				}(),
			},
			requestPayload: enrollmentUseCase.CancelCourseRequest{
				StudentID: studentID,
				CourseID:  courseID,
			},
			wantStatusCode: http.StatusOK,
			wantBody:       `{"status":"success"}`,
		},
		{
			name: "Invalid Request Payload",
			fields: fields{
				EnrollmentUseCase: nil,
			},
			requestPayload: enrollmentUseCase.CancelCourseRequest{
				StudentID: 0,
				CourseID:  courseID,
			},
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "Invalid request payload (empty)\n",
		},
		{
			name: "Error From UseCase",
			fields: fields{
				EnrollmentUseCase: func() enrollmentUseCase.EnrollmentUseCaseItf {
					mockEnrollmentUC := enrollmentUseCaseMock.NewMockEnrollmentUseCaseItf(ctrl)
					mockEnrollmentUC.EXPECT().CancelCourse(gomock.Any(), studentID, courseID).Return(enrollmentUseCase.CancelCourseResp{
						Status:  "error",
						Message: "failed to cancel course",
					}, errors.New("some error"))
					return mockEnrollmentUC
				}(),
			},
			requestPayload: enrollmentUseCase.CancelCourseRequest{
				StudentID: studentID,
				CourseID:  courseID,
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       "failed to cancel course\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				EnrollmentUseCase: tt.fields.EnrollmentUseCase,
			}

			body, _ := json.Marshal(tt.requestPayload)
			req := httptest.NewRequest(http.MethodPost, "/cancel-course", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			handler := h.CancelCourseHandler()
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Errorf("Status code = %v, want %v", rec.Code, tt.wantStatusCode)
			}

			if rec.Code == http.StatusOK {
				var gotBody, wantBody map[string]interface{}
				if err := json.Unmarshal(rec.Body.Bytes(), &gotBody); err != nil {
					t.Fatalf("Failed to unmarshal response body: %v", err)
				}
				if err := json.Unmarshal([]byte(tt.wantBody), &wantBody); err != nil {
					t.Fatalf("Failed to unmarshal expected body: %v", err)
				}
				if !reflect.DeepEqual(gotBody, wantBody) {
					t.Errorf("Response body = %v, want %v", gotBody, wantBody)
				}
			} else {
				if rec.Body.String() != tt.wantBody {
					t.Errorf("Response body = %v, want %v", rec.Body.String(), tt.wantBody)
				}
			}
		})
	}
}

func TestHandler_ListClassmatesHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		studentID int64 = 1
	)
	type fields struct {
		EnrollmentUseCase enrollmentUseCase.EnrollmentUseCaseItf
	}
	tests := []struct {
		name           string
		fields         fields
		queryParams    map[string]string
		mockResp       enrollmentUseCase.ListClassmatesResp
		mockErr        error
		wantStatusCode int
		wantBody       string
	}{
		{
			name: "Success",
			fields: fields{
				EnrollmentUseCase: func() enrollmentUseCase.EnrollmentUseCaseItf {
					mockEnrollmentUC := enrollmentUseCaseMock.NewMockEnrollmentUseCaseItf(ctrl)
					mockEnrollmentUC.EXPECT().ListClassmates(gomock.Any(), studentID).Return(enrollmentUseCase.ListClassmatesResp{
						Status: "success",
						Courses: []enrollmentUseCase.ListClassmatesCourseResp{
							{
								CourseID:   101,
								CourseName: "Course 101",
								ClassMates: []enrollmentUseCase.ListClassmatesStudentsResp{
									{
										StudentID:    "2",
										StudentEmail: "student2@example.com",
									},
									{
										StudentID:    "3",
										StudentEmail: "student3@example.com",
									},
								},
							},
						},
					}, nil)
					return mockEnrollmentUC
				}(),
			},
			queryParams: map[string]string{
				"student_id": strconv.FormatInt(studentID, 10),
			},
			wantStatusCode: http.StatusOK,
			wantBody: `{
				"status": "success",
				"courses": [
					{
						"course_id": 101,
						"course_name": "Course 101",
						"class_mates": [
							{"student_id": "2", "student_email": "student2@example.com"},
							{"student_id": "3", "student_email": "student3@example.com"}
						]
					}
				]
			}`,
		},
		{
			name: "Missing student_id",
			fields: fields{
				EnrollmentUseCase: nil,
			},
			queryParams:    map[string]string{},
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "student_id is required\n",
		},
		{
			name: "Invalid student_id",
			fields: fields{
				EnrollmentUseCase: nil,
			},
			queryParams: map[string]string{
				"student_id": "invalid",
			},
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "Invalid student_id\n",
		},
		{
			name: "Student ID Zero",
			fields: fields{
				EnrollmentUseCase: nil,
			},
			queryParams: map[string]string{
				"student_id": "0",
			},
			wantStatusCode: http.StatusBadRequest,
			wantBody:       "Invalid request payload (empty)\n",
		},
		{
			name: "Error from UseCase",
			fields: fields{
				EnrollmentUseCase: func() enrollmentUseCase.EnrollmentUseCaseItf {
					mockEnrollmentUC := enrollmentUseCaseMock.NewMockEnrollmentUseCaseItf(ctrl)
					mockEnrollmentUC.EXPECT().ListClassmates(gomock.Any(), studentID).Return(enrollmentUseCase.ListClassmatesResp{
						Status:  "error",
						Message: "failed to list classmates",
					}, errors.New("some error"))
					return mockEnrollmentUC
				}(),
			},
			queryParams: map[string]string{
				"student_id": strconv.FormatInt(studentID, 10),
			},
			wantStatusCode: http.StatusInternalServerError,
			wantBody:       "failed to list classmates\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				EnrollmentUseCase: tt.fields.EnrollmentUseCase,
			}

			query := "?"
			for key, value := range tt.queryParams {
				query += fmt.Sprintf("%s=%s&", key, value)
			}
			query = strings.TrimSuffix(query, "&")

			req := httptest.NewRequest(http.MethodGet, "/list-classmates"+query, nil)
			rec := httptest.NewRecorder()

			handler := h.ListClassmatesHandler()
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Errorf("Status code = %v, want %v", rec.Code, tt.wantStatusCode)
			}

			if rec.Code == http.StatusOK {
				var gotBody, wantBody map[string]interface{}
				if err := json.Unmarshal(rec.Body.Bytes(), &gotBody); err != nil {
					t.Fatalf("Failed to unmarshal response body: %v", err)
				}
				if err := json.Unmarshal([]byte(tt.wantBody), &wantBody); err != nil {
					t.Fatalf("Failed to unmarshal expected body: %v", err)
				}
				if !reflect.DeepEqual(gotBody, wantBody) {
					t.Errorf("Response body = %v, want %v", gotBody, wantBody)
				}
			} else {
				if rec.Body.String() != tt.wantBody {
					t.Errorf("Response body = %v, want %v", rec.Body.String(), tt.wantBody)
				}
			}
		})
	}
}
