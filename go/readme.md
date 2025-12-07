# Course Management System

A Go-based course management system with clean architecture and RESTful API design.

for refreshing go skills.

## Tech Stack

- **Golang**: Programming language used for development.
- **MySQL**: Database system for storing data.
- **Docker**: Containerization tool used for hosting MySQL.

## Project Overview

This system manages course enrollments with features for students to:
- Sign up for courses
- View enrolled courses
- Cancel course enrollments
- View classmates in enrolled courses

## Entities

The application features three main entities:

### Student
Represents student data with the following fields:
```
type Student struct {
	ID         int64
	Email      string
	CreateTime time.Time
	UpdateTime time.Time
}
```

### Course
Represents course data with the following fields:
```
type Course struct {
	ID         int64
	Name       string
	CreateTime time.Time
	UpdateTime time.Time
}
```

### Course Enrollment
Tracks student course enrollments with the following fields:
```
type CourseEnrollment struct {
	ID         int64
	StudentID  int64
	CourseID   int64
	Status     int
	CreateTime time.Time
	UpdateTime time.Time
}
```

## API Endpoints

### Course Enrollment Management

#### Enroll in a Course
```
POST /v1/enrollments
Content-Type: application/json

Request Body:
{
    "studentId": <int>,
    "courseId": <int>
}
```

#### List Enrolled Courses
```
GET /v1/students/{studentId}/courses
```

#### Cancel Course Enrollment
```
DELETE /v1/enrollments
Content-Type: application/json

Request Body:
{
    "studentId": <int>,
    "courseId": <int>
}
```

#### List Classmates
```
GET /v1/students/{studentId}/classmates
```

## Project Structure

```
├── cmd/            # Application entry point
├── common/         # Shared constants and utilities
├── db/            # Database schemas and migrations
├── domain/        # Business domain models and repositories
├── handlers/      # HTTP request handlers
├── routes/        # API route definitions
├── use-case/      # Business logic implementation
└── vendor/        # Dependencies (managed by Go modules)
```

## Makefile Commands

Once the dependencies are installed and configured, you can use the following commands to get started:

### Build and Run the Application
```
make run
```
This command will:
- Build the Go application and place the binary in the bin directory.
- The app will run on port 8991 (configured in etc/development.json)
- Run the binary with the specified configuration file.

### Start Docker Containers
```
make compose-up
```
This command will:
- Start the Docker containers defined in docker-compose.yml in detached mode.

### Stop Docker Containers
```
make compose-down
```
This command will:
- Stop the Docker containers defined in docker-compose.yml.

### Building dockerfile
```
make compose-build:
```
This command will:
- building the dockerfile

