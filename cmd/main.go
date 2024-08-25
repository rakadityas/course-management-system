package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github/rakadityas/course-management-system/config"
	coursedomain "github/rakadityas/course-management-system/domain/course"
	courseenrollmentdomain "github/rakadityas/course-management-system/domain/course-enrollment"
	studentdomain "github/rakadityas/course-management-system/domain/student"

	handlers "github/rakadityas/course-management-system/handlers"
	"github/rakadityas/course-management-system/routes"
	enrollmentusecase "github/rakadityas/course-management-system/use-case/enrollment"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	// initialize configuration
	configurationPath := flag.String("configuration_path", "etc/development.json", "file configuration path")
	flag.Parse()
	configuration := config.NewConfiguration(*configurationPath)

	// initialize database connection
	db, err := sql.Open("mysql", configuration.GetConfiguration().Resource.PrimaryDatabase)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

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

	fmt.Printf("Starting server on port %s\n", configuration.GetConfiguration().Server.HttpPort)
	if err := http.ListenAndServe(configuration.GetConfiguration().Server.HttpPort, router); err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
