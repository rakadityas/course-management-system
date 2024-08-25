package routes

import (
	"github/rakadityas/course-management-system/handlers"

	"github.com/gorilla/mux"
)

// SetupRoutes initializes the routes and returns the router.
func SetupRoutes(handler *handlers.Handler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/signup", handler.CourseSignUpHandler()).Methods("POST")
	r.HandleFunc("/courses", handler.ListCoursesHandler()).Methods("GET")
	r.HandleFunc("/cancel", handler.CancelCourseHandler()).Methods("POST")
	r.HandleFunc("/classmates", handler.ListClassmatesHandler()).Methods("GET")

	return r
}
