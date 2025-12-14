package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	coursedomain "github.com/rakadityas/course-management-system/go/domain/course"
	courseenrollmentdomain "github.com/rakadityas/course-management-system/go/domain/course-enrollment"
	studentdomain "github.com/rakadityas/course-management-system/go/domain/student"

	"log"
	"net/http"

	handlers "github.com/rakadityas/course-management-system/go/handlers"
	"github.com/rakadityas/course-management-system/go/routes"
	enrollmentusecase "github.com/rakadityas/course-management-system/go/use-case/enrollment"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Get the database connection string from the environment
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// Get the app port string from the environment
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		log.Fatal("APP_PORT is not set")
	}

	// TODO: find more elegant way for solving racing issue during docker compose up due to DB not yet ready
	time.Sleep(10 * time.Second)

	// initialize database connection
	db, err := sql.Open("mysql", databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// initialize domains
	studentService := studentdomain.NewStudentService(studentdomain.NewSQLStudentRepository(db))
	courseService := coursedomain.NewCourseService(coursedomain.NewSQLCourseRepository(db))
	courseEnrollmentService := courseenrollmentdomain.NewCourseEnrollmentService(courseenrollmentdomain.NewSQLCourseEnrollmentRepository(db))

	// initialize use cases
	enrollmentUseCase := enrollmentusecase.NewEnrollmentUseCase(studentService, courseService, courseEnrollmentService)

	// init http service
	handler := handlers.NewHandler(enrollmentUseCase)

	// Setup routes
	router := routes.SetupRoutes(handler)

	fmt.Printf("Starting server on port %s\n", appPort)
	if err := http.ListenAndServe(appPort, router); err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
