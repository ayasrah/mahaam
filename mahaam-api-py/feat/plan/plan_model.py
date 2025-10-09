from uuid import UUID
from datetime import datetime
from feat.user.user_model import User
from pydantic import BaseModel, Field

class Plan(BaseModel):
    id: UUID
    title: str | None = None
    type: str | None = None
    sort_order: int = Field(default=0, serialization_alias="sortOrder")
    starts: datetime | None = None
    ends: datetime | None = None
    done_percent: str | None = Field(default=None, serialization_alias="donePercent")
    created_at: datetime | None = Field(default=None, serialization_alias="createdAt")
    updated_at: datetime | None = Field(default=None, serialization_alias="updatedAt")
    members: list[User] | None = Field(default=None, serialization_alias="members")
    is_shared: bool = Field(default=False, serialization_alias="isShared")
    user: User | None = None


class PlanIn(BaseModel):
    id: UUID | None = None
    title: str | None = None
    starts: datetime | None = None
    ends: datetime | None = None


class PlanType:
    MAIN = "Main"
    ARCHIVED = "Archived"

    ALL = [MAIN, ARCHIVED]
