# Course Management System

A Python-based course management system with clean architecture and RESTful API design.

for refreshing python skills.

## Tech Stack

- **Python (FastAPI)**: Web framework for the HTTP API.
- **MySQL**: Database system for storing data.
- **Docker**: Containerization tool used for hosting MySQL and the app.

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
from dataclasses import dataclass
from datetime import datetime

@dataclass
class Student:
    id: int
    email: str
    create_time: datetime
    update_time: datetime
```

### Course
Represents course data with the following fields:
```
from dataclasses import dataclass
from datetime import datetime

@dataclass
class Course:
    id: int
    name: str
    create_time: datetime
    update_time: datetime
```

### Course Enrollment
Tracks student course enrollments with the following fields:
```
from dataclasses import dataclass
from datetime import datetime

@dataclass
class CourseEnrollment:
    id: int
    student_id: int
    course_id: int
    status: int
    create_time: datetime
    update_time: datetime
```

## API Endpoints

### Course Enrollment Management

#### Enroll in a Course
```
POST /v1/enrollments
Content-Type: application/json

Request Body:
{
    "student_id": <int>,
    "course_id": <int>
}
```

#### List Enrolled Courses
```
GET /v1/students/{student_id}/courses
```

#### Cancel Course Enrollment
```
DELETE /v1/enrollments
Content-Type: application/json

Request Body:
{
    "student_id": <int>,
    "course_id": <int>
}
```

#### List Classmates
```
GET /v1/students/{student_id}/classmates
```

## Project Structure

```
python/
├── cmd/            # Application entry point
├── common/         # Shared constants
├── domain/         # Business domain models, repositories, services
├── handlers/       # HTTP request handlers
├── routes/         # API route definitions and app wiring
├── use_case/       # Business logic implementation
├── utils/          # Utilities (e.g., DB connector)
├── Dockerfile      # Container image for the app
├── docker-compose.yaml
└── requirements.txt
```

## How to Run

### Build and Run with Docker Compose
```
cd python
docker compose up -d --build
```
This command will:
- Start the MySQL and Python app containers defined in `docker-compose.yaml`.
- Mount the schema and seed data from `../go/db` into MySQL.
- Run the FastAPI app on port 8992.

### Stop Docker Containers
```
docker compose down
```
This command will stop and remove the containers defined in `docker-compose.yaml`.

### Example Requests
```
curl -X POST http://localhost:8992/v1/enrollments \
  -H "Content-Type: application/json" \
  -d '{"student_id":1,"course_id":1}'

curl http://localhost:8992/v1/students/1/courses

curl -X DELETE http://localhost:8992/v1/enrollments \
  -H "Content-Type: application/json" \
  -d '{"student_id":1,"course_id":1}'

curl http://localhost:8992/v1/students/1/classmates
```
