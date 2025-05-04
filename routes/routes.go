package routes

import (
	"github/rakadityas/course-management-system/handlers"

	"github.com/gorilla/mux"
)

// SetupRoutes initializes the routes and returns the router.
func SetupRoutes(handler *handlers.Handler) *mux.Router {
	r := mux.NewRouter()

	// Course enrollment endpoints
	r.HandleFunc("/v1/enrollments", handler.CourseSignUpHandler()).Methods("POST")
	r.HandleFunc("/v1/students/{studentId}/courses", handler.ListCoursesHandler()).Methods("GET")
	r.HandleFunc("/v1/enrollments", handler.CancelCourseHandler()).Methods("DELETE")
	r.HandleFunc("/v1/students/{studentId}/classmates", handler.ListClassmatesHandler()).Methods("GET")

	return r
}
