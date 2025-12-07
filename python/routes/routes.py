from fastapi import FastAPI
from python.handlers.handlers import create_router
from python.use_case.enrollment.enrollment import EnrollmentUseCase
from python.domain.student.repository import StudentRepository
from python.domain.student.service import StudentService
from python.domain.course.repository import CourseRepository
from python.domain.course.service import CourseService
from python.domain.course_enrollment.repository import CourseEnrollmentRepository
from python.domain.course_enrollment.service import CourseEnrollmentService
from python.utils.db import DB

def setup_app(db: DB) -> FastAPI:
    student_service = StudentService(StudentRepository(db))
    course_service = CourseService(CourseRepository(db))
    ce_service = CourseEnrollmentService(CourseEnrollmentRepository(db))
    enrollment_uc = EnrollmentUseCase(student_service, course_service, ce_service)
    app = FastAPI()
    app.include_router(create_router(enrollment_uc))
    return app
