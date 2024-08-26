# Course Management System

## Tech Stack

- **Golang**: Programming language used for development.
- **MySQL**: Database system for storing data.
- **Docker**: Containerization tool used for hosting MySQL.

## Directory Structure

This project is structured based on Clean Architecture principles:

- **`bin`**: Contains the compiled binary files.
- **`cmd`**: Contains `main.go` file and entry point for the application.
- **`config`**: Contains configuration code.
- **`domain`**: Contains core entities such as students, courses, and course enrollment.
- **`etc`**: Contains plain configuration files.
- **`handlers`**: Contains API handlers.
- **`routes`**: Contains API route definitions.
- **`scripts`**: Contains DDL and DML scripts for database queries.
- **`use-case`**: Contains core business logic and use cases combining one or more domains.

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
- Docker MySQL Configuration:
    - Host: mysql
    - Username: myuser
    - Password: mypassword
    - Database: course_management


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



# API Documentation

## Endpoints

### 1. Sign Up for a Course
**Endpoint:** `POST /signup`

**Description:** Enroll a student in a course.

**Request Payload:**
```
{
  "student_id": 123,
  "course_id": 456
}
```
- student_id (int64): ID of the student.
- course_id (int64): ID of the course.


**Response:**

Success response
```
{
  "status": "success",
  "message": "Successfully enrolled",
  "enrollment_data": {
    "id": 1,
    "student_id": 123,
    "student_email": "student@example.com",
    "course_id": 456,
    "course_name": "Course Name",
    "status": 1,
    "create_time": "2024-08-25T12:34:56Z",
    "update_time": "2024-08-25T12:34:56Z"
  }
}
```

Failed response: request empty
```
{
  "status": "failure",
  "message": "Request Data is empty"
}
```

Failed response: student not found
```
{
  "status": "failure",
  "message": "student data not found"
}
```

Failed response: course not found
```
{
  "status": "failure",
  "message": "course data not found"
}
```

Failed response: student has enrolled before
```
{
  "status": "failure",
  "message": "student has enrolled before"
}
```


### 2. List Courses for a Student
**Endpoint:** `GET /courses`

**Description:** Retrieve a list of courses for a specific student.

**Query Parameters:**
- student_id (int64): ID of the student.

**Response:**

Success response
```
{
  "status": "success",
  "courses": [
    {
      "course_id": 456,
      "course_name": "Course Name",
      "status": 1,
      "create_time": "2024-08-25T12:34:56Z",
      "update_time": "2024-08-25T12:34:56Z"
    }
  ]
}
```

Failure response: invalid student id
```
{
  "status": "failure",
  "message": "Invalid student ID"
}
```

Failure response: Student ID Zero
```
{
  "status": "failure",
  "message": "Student ID Zero"
}
```

Failure response: course data is not found
```
{
  "status": "failure",
  "message": "course data is not found for courseID: 10"
}
```



### 3. Cancel a Course Enrollment
**Endpoint:** `POST /cancel`

**Description:** Cancel a student's enrollment in a course.

**Request Payload:**
```
{
  "student_id": 123,
  "course_id": 456
}
```
- student_id (int64): ID of the student.
- course_id (int64): ID of the course.

**Response:**

Success response
```
{
  "status": "success"
}
```

Failure response: Invalid Request Payload (empty)
```
{
  "status": "failure",
  "message": "Invalid request payload (empty)"
}
```

### 4. List Classmates
**Endpoint:** `GET /classmates`

**Description:** Get a list of classmates enrolled in the same courses as the given student.

**Query Parameter:**
- student_id (int64): ID of the student.

**Response:**

success response:
```
{
  "status": "success",
  "courses": [
    {
      "course_id": 101,
      "course_name": "Course 101",
      "class_mates": [
        {
          "student_id": "2",
          "student_email": "student2@example.com"
        },
        {
          "student_id": "3",
          "student_email": "student3@example.com"
        }
      ]
    }
  ]
}

```

Failure response: Missing student_id
```
{
  "status": "failure",
  "message": "student_id is required"
}
```


Failure response: Invalid student_id
```
{
  "status": "failure",
  "message": "Invalid student_id"
}
```


Failure response: Invalid request payload (empty)
```
{
  "status": "failure",
  "message": "Invalid request payload (empty)"
}
```

Failure response: course data is not found
```
{
  "status": "failure",
  "message": "course data is not found for courseID: 10"
}
```

Failure response: student data is not found
```
{
  "status": "failure",
  "message": "student data is not found for studentID: 10"
}
```