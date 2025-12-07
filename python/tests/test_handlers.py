import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient
from python.handlers.handlers import create_router
from datetime import datetime
from python.use_case.enrollment.types import (
    CourseSignUpResp,
    CourseEnrollmentResp,
    ListCoursesResp,
    CancelCourseResp,
    ListClassmatesResp,
    ListClassmatesCourseResp,
    ListClassmatesStudentsResp,
    CourseDetail,
)

class FakeUC:
    def course_sign_up(self, req):
        now = datetime.utcnow()
        return CourseSignUpResp(status="success", enrollment_data=CourseEnrollmentResp(id=1, student_id=req.student_id, student_email="student1@example.com", course_id=req.course_id, course_name="Mathematics 101", status=1, create_time=now, update_time=now))
    def list_courses(self, student_id: int):
        now = datetime.utcnow()
        return ListCoursesResp(status="success", courses=[CourseDetail(course_id=1, course_name="Mathematics 101", status=1, create_time=now, update_time=now)])
    def cancel_course(self, student_id: int, course_id: int):
        return CancelCourseResp(status="success")
    def list_classmates(self, student_id: int):
        mates = [ListClassmatesStudentsResp(student_id="2", student_email="student2@example.com")]
        return ListClassmatesResp(status="success", courses=[ListClassmatesCourseResp(course_id=1, course_name="Mathematics 101", class_mates=mates)])

def make_client():
    app = FastAPI()
    app.include_router(create_router(FakeUC()))
    return TestClient(app)

def test_post_enrollments():
    c = make_client()
    r = c.post("/v1/enrollments", json={"student_id": 1, "course_id": 1})
    assert r.status_code == 201
    data = r.json()
    assert data["status"] == "success"
    assert data["enrollment_data"]["course_id"] == 1

def test_get_courses():
    c = make_client()
    r = c.get("/v1/students/1/courses")
    assert r.status_code == 200
    assert r.json()["status"] == "success"

def test_delete_enrollments():
    c = make_client()
    r = c.request("DELETE", "/v1/enrollments", json={"student_id": 1, "course_id": 1})
    assert r.status_code == 200
    assert r.json()["status"] == "success"

def test_get_classmates():
    c = make_client()
    r = c.get("/v1/students/1/classmates")
    assert r.status_code == 200
    assert r.json()["status"] == "success"
