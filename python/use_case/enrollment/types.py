from datetime import datetime
from typing import List, Optional
from pydantic import BaseModel

class CourseSignUpRequest(BaseModel):
    student_id: int
    course_id: int

class CourseEnrollmentResp(BaseModel):
    id: int
    student_id: int
    student_email: str
    course_id: int
    course_name: str
    status: int
    create_time: datetime
    update_time: datetime

class CourseSignUpResp(BaseModel):
    status: str
    message: Optional[str] = None
    enrollment_data: Optional[CourseEnrollmentResp] = None

class CourseDetail(BaseModel):
    course_id: int
    course_name: str
    status: int
    create_time: datetime
    update_time: datetime

class ListCoursesResp(BaseModel):
    status: str
    message: Optional[str] = None
    courses: List[CourseDetail] = []

class CancelCourseRequest(BaseModel):
    student_id: int
    course_id: int

class CancelCourseResp(BaseModel):
    status: str
    message: Optional[str] = None

class ListClassmatesStudentsResp(BaseModel):
    student_id: str
    student_email: str

class ListClassmatesCourseResp(BaseModel):
    course_id: int
    course_name: str
    class_mates: List[ListClassmatesStudentsResp]

class ListClassmatesResp(BaseModel):
    status: str
    message: Optional[str] = None
    courses: List[ListClassmatesCourseResp] = []
