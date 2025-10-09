from uuid import UUID
from datetime import datetime
from dataclasses import dataclass
from pydantic import BaseModel, Field


@dataclass
class User:
    id: UUID
    email: str | None = None
    status: str | None = None
    name: str | None = None


class User2(BaseModel):
    id: UUID
    email: str | None = None
    status: str | None = None
    name: str | None = None
    user_id: UUID | None = Field(default=None, serialization_alias="userId")


class Device(BaseModel):
    id: UUID | None = None
    user_id: UUID | None = Field(default=None, serialization_alias="userId")
    platform: str | None = None
    fingerprint: str = None
    info: str | None = None
    created_at: datetime | None = Field(default=None, serialization_alias="createdAt")


class SuggestedEmail(BaseModel):
    id: UUID
    user_id: UUID = Field(serialization_alias="userId")
    email: str | None = None
    created_at: datetime | None = Field(default=None, serialization_alias="createdAt")


class VerifiedUser(BaseModel):
    user_id: UUID = Field(serialization_alias="userId")
    device_id: UUID = Field(serialization_alias="deviceId")
    jwt: str
    user_full_name: str | None = Field(default=None, serialization_alias="userFullName")
    email: str | None = None


class CreatedUser(BaseModel):
    id: UUID
    device_id: UUID = Field(serialization_alias="deviceId")
    jwt: str
