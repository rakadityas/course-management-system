from typing import Optional
from .models import Course
from .repository import CourseRepository

class CourseService:
    def __init__(self, repo: CourseRepository):
        self.repo = repo

    def get_course_by_id(self, course_id: int) -> Optional[Course]:
        return self.repo.get_course_by_id(course_id)
