from dataclasses import dataclass
from datetime import datetime

StatusActive = 1
StatusCancelled = 0

@dataclass
class CourseEnrollment:
    id: int
    student_id: int
    course_id: int
    status: int
    create_time: datetime
    update_time: datetime
