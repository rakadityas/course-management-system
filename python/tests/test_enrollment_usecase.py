import pytest
from datetime import datetime
from python.use_case.enrollment.enrollment import EnrollmentUseCase
from python.domain.student.models import Student
from python.domain.course.models import Course
from python.domain.course_enrollment.models import CourseEnrollment, StatusActive, StatusCancelled

class FakeStudentService:
    def __init__(self, data):
        self.data = data
    def get_student_by_id(self, student_id: int):
        return self.data.get(student_id)

class FakeCourseService:
    def __init__(self, data):
        self.data = data
    def get_course_by_id(self, course_id: int):
        return self.data.get(course_id)

class FakeCourseEnrollmentService:
    def __init__(self):
        self.items = []
    def create_enrollment(self, student_id: int, course_id: int, status: int):
        now = datetime.utcnow()
        item = CourseEnrollment(id=len(self.items)+1, student_id=student_id, course_id=course_id, status=status, create_time=now, update_time=now)
        self.items.append(item)
        return item
    def get_enrollment_by_student_id(self, student_id: int):
        return [e for e in self.items if e.student_id == student_id and e.status == StatusActive]
    def update_course_enrollment_status(self, student_id: int, course_id: int, new_status: int):
        updated = False
        for e in self.items:
            if e.student_id == student_id and e.course_id == course_id:
                e.status = new_status
                e.update_time = datetime.utcnow()
                updated = True
        return updated
    def get_enrollment_by_student_and_course(self, student_id: int, course_id: int):
        return [e for e in self.items if e.student_id == student_id and e.course_id == course_id]
    def get_list_classmates(self, student_id: int):
        courses = {e.course_id for e in self.items if e.student_id == student_id and e.status == StatusActive}
        return [e for e in self.items if e.course_id in courses and e.student_id != student_id and e.status == StatusActive]

def make_usecase():
    students = {
        1: Student(id=1, email="student1@example.com", create_time=datetime.utcnow(), update_time=datetime.utcnow()),
        2: Student(id=2, email="student2@example.com", create_time=datetime.utcnow(), update_time=datetime.utcnow()),
        3: Student(id=3, email="student3@example.com", create_time=datetime.utcnow(), update_time=datetime.utcnow()),
    }
    courses = {
        1: Course(id=1, name="Mathematics 101", create_time=datetime.utcnow(), update_time=datetime.utcnow()),
        2: Course(id=2, name="Introduction to Programming", create_time=datetime.utcnow(), update_time=datetime.utcnow()),
        3: Course(id=3, name="History of Art", create_time=datetime.utcnow(), update_time=datetime.utcnow()),
    }
    ss = FakeStudentService(students)
    cs = FakeCourseService(courses)
    es = FakeCourseEnrollmentService()
    return EnrollmentUseCase(ss, cs, es), es

def test_course_sign_up_success():
    uc, es = make_usecase()
    resp = uc.course_sign_up(type("R", (), {"student_id": 1, "course_id": 1})())
    assert resp.status == "success"
    assert resp.enrollment_data is not None
    assert resp.enrollment_data.student_id == 1
    assert resp.enrollment_data.course_id == 1

def test_course_sign_up_duplicate():
    uc, es = make_usecase()
    uc.course_sign_up(type("R", (), {"student_id": 1, "course_id": 1})())
    resp = uc.course_sign_up(type("R", (), {"student_id": 1, "course_id": 1})())
    assert resp.status == "failure"

def test_list_courses_success():
    uc, es = make_usecase()
    uc.course_sign_up(type("R", (), {"student_id": 2, "course_id": 1})())
    uc.course_sign_up(type("R", (), {"student_id": 2, "course_id": 2})())
    resp = uc.list_courses(2)
    assert resp.status == "success"
    assert len(resp.courses) == 2

def test_cancel_course_success():
    uc, es = make_usecase()
    uc.course_sign_up(type("R", (), {"student_id": 1, "course_id": 2})())
    resp = uc.cancel_course(1, 2)
    assert resp.status == "success"
    active = es.get_enrollment_by_student_id(1)
    assert all(e.course_id != 2 or e.status == StatusCancelled for e in es.items)

def test_list_classmates_success():
    uc, es = make_usecase()
    uc.course_sign_up(type("R", (), {"student_id": 1, "course_id": 1})())
    uc.course_sign_up(type("R", (), {"student_id": 2, "course_id": 1})())
    uc.course_sign_up(type("R", (), {"student_id": 3, "course_id": 1})())
    resp = uc.list_classmates(1)
    assert resp.status == "success"
    assert len(resp.courses) == 1
    course = resp.courses[0]
    ids = {m.student_id for m in course.class_mates}
    assert ids == {"2", "3"}
