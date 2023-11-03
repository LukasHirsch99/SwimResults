from datetime import datetime
from enum import Enum
from dataclasses import dataclass, fields
from dataclasses_json import DataClassJsonMixin


class Gender(Enum):
    male = "male"
    female = "female"
    mixed = "mixed"


@dataclass
class Meet(DataClassJsonMixin):
    id: int
    name: str
    image: str
    invitations: list[str]
    deadline: str
    address: str
    googlemapslink: str
    startdate: str
    enddate: str


@dataclass
class Session:
    id: int
    meetId: int
    warmupStart: datetime
    sessionStart: datetime
    day: datetime


@dataclass
class SessionItem:
    id: int
    sessionId: int
    event: str
    sessionItemNr: int


@dataclass
class Club(DataClassJsonMixin):
    id: int
    name: str
    nationality: str

    def __hash__(self) -> int:
        return hash(str(self.id) + self.name + self.nationality)


@dataclass
class Swimmer:
    id: int
    name: str
    birthyear: int
    clubId: int
    gender: Gender


@dataclass
class Start:
    id: int
    sessionItemId: int
    heat: int
    lane: int
    swimmerId: int
    time: datetime

    def __str__(self) -> str:
        if self.time == None:
            timeString = "null"
        else:
            timeString = self.time.strftime("'%H:%M:%S.%f'")
        return f"({self.id}, {self.sessionItemId}, {self.heat}, {self.lane}, {self.swimmerId}, {timeString})"


@dataclass
class Result:
    id: int
    sessionItemId: int
    ageClass: str
    swimmerId: int
    position: int
    time: datetime
    additionalTimeInfo: str
    splits: str

    def __str__(self) -> str:
        if self.time == None:
            timeString = "null"
        else:
            timeString = self.time.strftime("'%H:%M:%S.%f'")

        if self.position == None:
            positionString = "null"
        else:
            positionString = str(self.position)
        return f"({self.id}, {self.sessionItemId}, '{self.ageClass}', {positionString}, {self.swimmerId}, {timeString}, '{self.additionalTimeInfo}', '{self.splits}')"
