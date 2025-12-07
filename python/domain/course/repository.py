from typing import Optional
from .models import Course

class CourseRepository:
    def __init__(self, db):
        self.db = db

    def get_course_by_id(self, course_id: int) -> Optional[Course]:
        conn = self.db.get_connection()
        try:
            cur = conn.cursor()
            cur.execute(
                "SELECT id, name, create_time, update_time FROM courses WHERE id = %s",
                (course_id,),
            )
            row = cur.fetchone()
            if not row:
                return None
            return Course(id=row[0], name=row[1], create_time=row[2], update_time=row[3])
        finally:
            conn.close()
