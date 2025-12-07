from dataclasses import dataclass
from datetime import datetime

@dataclass
class Student:
    id: int
    email: str
    create_time: datetime
    update_time: datetime
