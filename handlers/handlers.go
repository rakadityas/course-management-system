package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	common "github/rakadityas/course-management-system/common"
	enrollmentUseCase "github/rakadityas/course-management-system/use-case/enrollment"
)

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

// CourseSignUpHandler handles the course sign-up process.
func (h *Handler) CourseSignUpHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var requestPayload enrollmentUseCase.CourseSignUpRequest
		if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
			statusByte, _ := json.Marshal(HandlerStatus{Status: common.StatusFailure, Message: "Invalid request payload"})
			http.Error(w, string(statusByte), http.StatusBadRequest)
			return
		}
		if requestPayload.StudentID == 0 || requestPayload.CourseID == 0 {
			statusByte, _ := json.Marshal(HandlerStatus{Status: common.StatusFailure, Message: "Request Data is empty"})
			http.Error(w, string(statusByte), http.StatusBadRequest)
			return
		}

		resp, err := h.EnrollmentUseCase.CourseSignUp(ctx, requestPayload)
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

// ListCoursesHandler handles the listing of courses for a student.
func (h *Handler) ListCoursesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		studentIDParam := r.URL.Query().Get("student_id")
		studentID, err := strconv.ParseInt(studentIDParam, 10, 64)
		if err != nil {
			statusByte, _ := json.Marshal(HandlerStatus{Status: common.StatusFailure, Message: "Invalid student ID"})
			http.Error(w, string(statusByte), http.StatusBadRequest)
			return
		}
		if studentID == 0 {
			statusByte, _ := json.Marshal(HandlerStatus{Status: common.StatusFailure, Message: "Student ID Zero"})
			http.Error(w, string(statusByte), http.StatusBadRequest)
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

// CancelCourseHandler handles the cancellation of a course enrollment.
func (h *Handler) CancelCourseHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var requestPayload enrollmentUseCase.CancelCourseRequest
		if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
			statusByte, _ := json.Marshal(HandlerStatus{Status: common.StatusFailure, Message: "Invalid request payload"})
			http.Error(w, string(statusByte), http.StatusBadRequest)
			return
		}
		if requestPayload.CourseID == 0 || requestPayload.StudentID == 0 {
			statusByte, _ := json.Marshal(HandlerStatus{Status: common.StatusFailure, Message: "Invalid request payload (empty)"})
			http.Error(w, string(statusByte), http.StatusBadRequest)
			return
		}

		resp, err := h.EnrollmentUseCase.CancelCourse(ctx, requestPayload.StudentID, requestPayload.CourseID)
		if err != nil {
			statusResp, _ := json.Marshal(resp)
			http.Error(w, string(statusResp), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

// ListClassmatesHandler handles requests to list the classmates of a student.
func (h *Handler) ListClassmatesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		studentIDStr := r.URL.Query().Get("student_id")
		if studentIDStr == "" {
			statusByte, _ := json.Marshal(HandlerStatus{Status: common.StatusFailure, Message: "student_id is required"})
			http.Error(w, string(statusByte), http.StatusBadRequest)
			return
		}

		studentID, err := strconv.ParseInt(studentIDStr, 10, 64)
		if err != nil {
			statusByte, _ := json.Marshal(HandlerStatus{Status: common.StatusFailure, Message: "Invalid student_id"})
			http.Error(w, string(statusByte), http.StatusBadRequest)
			return
		}
		if studentID == 0 {
			statusByte, _ := json.Marshal(HandlerStatus{Status: common.StatusFailure, Message: "Invalid request payload (empty)"})
			http.Error(w, string(statusByte), http.StatusBadRequest)
			return
		}

		resp, err := h.EnrollmentUseCase.ListClassmates(ctx, studentID)
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
