from uuid import UUID
from datetime import datetime
from pydantic import BaseModel, Field


class Task(BaseModel):
    id: UUID
    plan_id: UUID = Field(serialization_alias="planId")
    title: str
    done: bool = False
    sort_order: int = Field(default=0, serialization_alias="sortOrder")
    created_at: datetime | None = Field(default=None, serialization_alias="createdAt")
    updated_at: datetime | None = Field(default=None, serialization_alias="updatedAt")
