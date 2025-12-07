from typing import List
from .models import CourseEnrollment
from datetime import datetime

class CourseEnrollmentRepository:
    def __init__(self, db):
        self.db = db

    def create_enrollment(self, student_id: int, course_id: int, status: int) -> CourseEnrollment:
        now = datetime.utcnow()
        conn = self.db.get_connection()
        try:
            cur = conn.cursor()
            cur.execute(
                "INSERT INTO course_enrollments (student_id, course_id, status, create_time, update_time) VALUES (%s, %s, %s, %s, %s)",
                (student_id, course_id, status, now, now),
            )
            conn.commit()
            enrollment_id = cur.lastrowid
            return CourseEnrollment(id=enrollment_id, student_id=student_id, course_id=course_id, status=status, create_time=now, update_time=now)
        finally:
            conn.close()

    def get_enrollment_by_student_id(self, student_id: int) -> List[CourseEnrollment]:
        conn = self.db.get_connection()
        try:
            cur = conn.cursor()
            cur.execute(
                "SELECT id, student_id, course_id, status, create_time, update_time FROM course_enrollments WHERE student_id = %s AND status = 1",
                (student_id,),
            )
            rows = cur.fetchall()
            return [
                CourseEnrollment(id=r[0], student_id=r[1], course_id=r[2], status=r[3], create_time=r[4], update_time=r[5])
                for r in rows
            ]
        finally:
            conn.close()

    def update_course_enrollment_status(self, student_id: int, course_id: int, new_status: int) -> bool:
        conn = self.db.get_connection()
        try:
            cur = conn.cursor()
            cur.execute(
                "UPDATE course_enrollments SET status = %s, update_time = %s WHERE student_id = %s AND course_id = %s",
                (new_status, datetime.utcnow(), student_id, course_id),
            )
            conn.commit()
            return cur.rowcount > 0
        finally:
            conn.close()

    def get_enrollment_by_student_and_course(self, student_id: int, course_id: int) -> List[CourseEnrollment]:
        conn = self.db.get_connection()
        try:
            cur = conn.cursor()
            cur.execute(
                "SELECT id, student_id, course_id, status, create_time, update_time FROM course_enrollments WHERE student_id = %s AND course_id = %s",
                (student_id, course_id),
            )
            rows = cur.fetchall()
            return [
                CourseEnrollment(id=r[0], student_id=r[1], course_id=r[2], status=r[3], create_time=r[4], update_time=r[5])
                for r in rows
            ]
        finally:
            conn.close()

    def get_list_classmates(self, student_id: int) -> List[CourseEnrollment]:
        conn = self.db.get_connection()
        try:
            cur = conn.cursor()
            cur.execute(
                """
                SELECT ce.id, ce.student_id, ce.course_id, ce.status, ce.create_time, ce.update_time
                FROM course_enrollments ce
                JOIN course_enrollments ce2 ON ce.course_id = ce2.course_id
                WHERE ce2.student_id = %s AND ce.student_id != %s AND ce2.status = 1 AND ce.status = 1
                """,
                (student_id, student_id),
            )
            rows = cur.fetchall()
            return [
                CourseEnrollment(id=r[0], student_id=r[1], course_id=r[2], status=r[3], create_time=r[4], update_time=r[5])
                for r in rows
            ]
        finally:
            conn.close()
