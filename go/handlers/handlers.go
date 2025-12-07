package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	common "github/rakadityas/course-management-system/go/common"
	enrollmentUseCase "github/rakadityas/course-management-system/go/use-case/enrollment"

	"github.com/gorilla/mux"
)

// responseError sends a JSON error response with consistent format
func responseError(w http.ResponseWriter, message string, statusCode int) {
	response := HandlerStatus{Status: common.StatusFailure, Message: message}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// Handler struct holds the services required for handling requests.
type Handler struct {
	EnrollmentUseCase enrollmentUseCase.EnrollmentUseCaseItf
}

// NewHandler creates a new Handler instance with the provided services.
func NewHandler(enrollmentUC enrollmentUseCase.EnrollmentUseCaseItf) *Handler {
	return &Handler{
		EnrollmentUseCase: enrollmentUC,
	}
}

// CourseSignUpHandler handles POST /v1/enrollments to enroll a student in a course.
func (h *Handler) CourseSignUpHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		w.Header().Set("Content-Type", "application/json")

		var requestPayload enrollmentUseCase.CourseSignUpRequest
		if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
			responseError(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		if requestPayload.StudentID == 0 || requestPayload.CourseID == 0 {
			responseError(w, "Student ID and Course ID are required", http.StatusBadRequest)
			return
		}

		resp, err := h.EnrollmentUseCase.CourseSignUp(ctx, requestPayload)
		if err != nil {
			responseError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}

// ListCoursesHandler handles GET /v1/students/{studentId}/courses to list enrolled courses.
func (h *Handler) ListCoursesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)
		studentID, err := strconv.ParseInt(vars["studentId"], 10, 64)
		if err != nil {
			responseError(w, "Invalid student ID format", http.StatusBadRequest)
			return
		}

		if studentID == 0 {
			responseError(w, "Invalid student ID", http.StatusBadRequest)
			return
		}

		resp, err := h.EnrollmentUseCase.ListCourses(ctx, studentID)
		if err != nil {
			respByte, _ := json.Marshal(resp)
			http.Error(w, string(respByte), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

// CancelCourseHandler handles DELETE /v1/enrollments to cancel a course enrollment.
func (h *Handler) CancelCourseHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		w.Header().Set("Content-Type", "application/json")

		var requestPayload enrollmentUseCase.CancelCourseRequest
		if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
			responseError(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		if requestPayload.CourseID == 0 || requestPayload.StudentID == 0 {
			responseError(w, "Student ID and Course ID are required", http.StatusBadRequest)
			return
		}

		resp, err := h.EnrollmentUseCase.CancelCourse(ctx, requestPayload.StudentID, requestPayload.CourseID)
		if err != nil {
			responseError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

// ListClassmatesHandler handles GET /v1/students/{studentId}/classmates to list student's classmates.
func (h *Handler) ListClassmatesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)
		studentID, err := strconv.ParseInt(vars["studentId"], 10, 64)
		if err != nil {
			responseError(w, "Invalid student ID format", http.StatusBadRequest)
			return
		}

		if studentID == 0 {
			responseError(w, "Invalid student ID", http.StatusBadRequest)
			return
		}

		resp, err := h.EnrollmentUseCase.ListClassmates(ctx, studentID)
		if err != nil {
			responseError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}
