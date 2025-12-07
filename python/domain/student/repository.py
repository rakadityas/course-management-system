from typing import Optional
from .models import Student

class StudentRepository:
    def __init__(self, db):
        self.db = db

    def get_student_by_id(self, student_id: int) -> Optional[Student]:
        conn = self.db.get_connection()
        try:
            cur = conn.cursor()
            cur.execute(
                "SELECT id, email, create_time, update_time FROM students WHERE id = %s",
                (student_id,),
            )
            row = cur.fetchone()
            if not row:
                return None
            return Student(id=row[0], email=row[1], create_time=row[2], update_time=row[3])
        finally:
            conn.close()
