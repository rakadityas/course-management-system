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

## Setup and Development

1. Clone the repository
2. Install dependencies:
   ```
   go mod download
   ```
3. Set up the database:
   ```
   make migrate
   ```
4. Run the application:
   ```
   make run
   ```

## Testing

Run tests with:
```
go test ./...
```

## Docker Support

Build and run with Docker:
```
docker-compose up --build
```
