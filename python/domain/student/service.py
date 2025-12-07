from typing import Optional
from .models import Student
from .repository import StudentRepository

class StudentService:
    def __init__(self, repo: StudentRepository):
        self.repo = repo

    def get_student_by_id(self, student_id: int) -> Optional[Student]:
        return self.repo.get_student_by_id(student_id)
