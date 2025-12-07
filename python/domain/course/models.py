from dataclasses import dataclass
from datetime import datetime

@dataclass
class Course:
    id: int
    name: str
    create_time: datetime
    update_time: datetime
