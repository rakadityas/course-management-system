from typing import List
from .types import (
    CourseSignUpRequest,
    CourseSignUpResp,
    CourseEnrollmentResp,
    ListCoursesResp,
    CourseDetail,
    CancelCourseResp,
    ListClassmatesResp,
    ListClassmatesCourseResp,
    ListClassmatesStudentsResp,
)
from python.common.const import StatusFailure, StatusSuccess
from python.domain.student.service import StudentService
from python.domain.course.service import CourseService
from python.domain.course_enrollment.service import CourseEnrollmentService
from python.domain.course_enrollment.models import StatusActive, StatusCancelled

class EnrollmentUseCase:
    def __init__(self, student_service: StudentService, course_service: CourseService, course_enrollment_service: CourseEnrollmentService):
        self.student_service = student_service
        self.course_service = course_service
        self.course_enrollment_service = course_enrollment_service

    def course_sign_up(self, req: CourseSignUpRequest) -> CourseSignUpResp:
        student = self.student_service.get_student_by_id(req.student_id)
        if student is None:
            return CourseSignUpResp(status=StatusFailure, message="student data not found")

        course = self.course_service.get_course_by_id(req.course_id)
        if course is None:
            return CourseSignUpResp(status=StatusFailure, message="course data not found")

        existing = self.course_enrollment_service.get_enrollment_by_student_and_course(req.student_id, req.course_id)
        if len(existing) > 0:
            return CourseSignUpResp(status=StatusFailure, message="student has enrolled before")

        new_enrollment = self.course_enrollment_service.create_enrollment(req.student_id, req.course_id, StatusActive)
        return CourseSignUpResp(
            status=StatusSuccess,
            enrollment_data=CourseEnrollmentResp(
                id=new_enrollment.id,
                student_id=new_enrollment.student_id,
                student_email=student.email,
                course_id=new_enrollment.course_id,
                course_name=course.name,
                status=StatusActive,
                create_time=new_enrollment.create_time,
                update_time=new_enrollment.update_time,
            ),
        )

    def list_courses(self, student_id: int) -> ListCoursesResp:
        student = self.student_service.get_student_by_id(student_id)
        if student is None:
            return ListCoursesResp(status=StatusFailure, message="student data not found")

        enrollments = self.course_enrollment_service.get_enrollment_by_student_id(student_id)
        courses: List[CourseDetail] = []
        for e in enrollments:
            course = self.course_service.get_course_by_id(e.course_id)
            if course is None:
                return ListCoursesResp(status=StatusFailure, message=f"course data is not found for courseID: {e.course_id}")
            courses.append(
                CourseDetail(
                    course_id=course.id,
                    course_name=course.name,
                    status=e.status,
                    create_time=e.create_time,
                    update_time=e.update_time,
                )
            )
        return ListCoursesResp(status=StatusSuccess, courses=courses)

    def cancel_course(self, student_id: int, course_id: int) -> CancelCourseResp:
        updated = self.course_enrollment_service.update_course_enrollment_status(student_id, course_id, StatusCancelled)
        if not updated:
            return CancelCourseResp(status=StatusFailure, message="failed to cancel course enrollment")
        return CancelCourseResp(status=StatusSuccess)

    def list_classmates(self, student_id: int) -> ListClassmatesResp:
        student = self.student_service.get_student_by_id(student_id)
        if student is None:
            return ListClassmatesResp(status=StatusFailure, message="student data not found")
        enrollments = self.course_enrollment_service.get_list_classmates(student_id)
        grouped = {}
        for e in enrollments:
            grouped.setdefault(e.course_id, []).append(e.student_id)
        courses_resp = []
        for course_id, student_ids in grouped.items():
            course = self.course_service.get_course_by_id(course_id)
            if course is None:
                return ListClassmatesResp(status=StatusFailure, message=f"course data is not found for courseID: {course_id}")
            classmates = []
            for sid in student_ids:
                if sid == student_id:
                    continue
                s = self.student_service.get_student_by_id(sid)
                if s is None:
                    return ListClassmatesResp(status=StatusFailure, message=f"student data is not found for studentID: {sid}")
                classmates.append(ListClassmatesStudentsResp(student_id=str(s.id), student_email=s.email))
            courses_resp.append(ListClassmatesCourseResp(course_id=course.id, course_name=course.name, class_mates=classmates))
        return ListClassmatesResp(status=StatusSuccess, courses=courses_resp)
