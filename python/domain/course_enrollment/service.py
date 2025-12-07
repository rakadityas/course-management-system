from typing import List
from .repository import CourseEnrollmentRepository
from .models import CourseEnrollment

class CourseEnrollmentService:
    def __init__(self, repo: CourseEnrollmentRepository):
        self.repo = repo

    def create_enrollment(self, student_id: int, course_id: int, status: int) -> CourseEnrollment:
        return self.repo.create_enrollment(student_id, course_id, status)

    def get_enrollment_by_student_id(self, student_id: int) -> List[CourseEnrollment]:
        return self.repo.get_enrollment_by_student_id(student_id)

    def update_course_enrollment_status(self, student_id: int, course_id: int, new_status: int) -> bool:
        return self.repo.update_course_enrollment_status(student_id, course_id, new_status)

    def get_enrollment_by_student_and_course(self, student_id: int, course_id: int) -> List[CourseEnrollment]:
        return self.repo.get_enrollment_by_student_and_course(student_id, course_id)

    def get_list_classmates(self, student_id: int) -> List[CourseEnrollment]:
        return self.repo.get_list_classmates(student_id)
