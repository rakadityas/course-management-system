from fastapi import APIRouter, HTTPException
from fastapi.responses import JSONResponse
from python.common.const import StatusFailure
from python.use_case.enrollment.enrollment import EnrollmentUseCase
from python.use_case.enrollment.types import (
    CourseSignUpRequest,
    CancelCourseRequest,
)

def create_router(enrollment_uc: EnrollmentUseCase) -> APIRouter:
    router = APIRouter()

    @router.post("/v1/enrollments")
    def course_sign_up(req: CourseSignUpRequest):
        if req.student_id == 0 or req.course_id == 0:
            raise HTTPException(status_code=400, detail={"status": StatusFailure, "message": "Student ID and Course ID are required"})
        resp = enrollment_uc.course_sign_up(req)
        if resp.status == StatusFailure:
            return JSONResponse(status_code=500, content=resp.model_dump(mode="json"))
        return JSONResponse(status_code=201, content=resp.model_dump(mode="json"))

    @router.get("/v1/students/{student_id}/courses")
    def list_courses(student_id: int):
        resp = enrollment_uc.list_courses(student_id)
        if resp.status == StatusFailure:
            return JSONResponse(status_code=500, content=resp.model_dump(mode="json"))
        return JSONResponse(status_code=200, content=resp.model_dump(mode="json"))

    @router.delete("/v1/enrollments")
    def cancel_course(req: CancelCourseRequest):
        if req.student_id == 0 or req.course_id == 0:
            raise HTTPException(status_code=400, detail={"status": StatusFailure, "message": "Student ID and Course ID are required"})
        resp = enrollment_uc.cancel_course(req.student_id, req.course_id)
        if resp.status == StatusFailure:
            return JSONResponse(status_code=500, content=resp.model_dump(mode="json"))
        return JSONResponse(status_code=200, content=resp.model_dump(mode="json"))

    @router.get("/v1/students/{student_id}/classmates")
    def list_classmates(student_id: int):
        resp = enrollment_uc.list_classmates(student_id)
        if resp.status == StatusFailure:
            return JSONResponse(status_code=500, content=resp.model_dump(mode="json"))
        return JSONResponse(status_code=200, content=resp.model_dump(mode="json"))

    return router
