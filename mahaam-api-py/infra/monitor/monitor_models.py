from uuid import UUID
from datetime import datetime
from pydantic import BaseModel, Field


class Traffic(BaseModel):
    id: UUID
    health_id: UUID = Field(serialization_alias="healthId")
    method: str
    path: str
    code: int | None = None
    elapsed: int | None = None
    headers: str | None = None
    request: str | None = None
    response: str | None = None


class TrafficHeaders(BaseModel):
    user_id: UUID | None = Field(default=None, serialization_alias="userId")
    device_id: UUID | None = Field(default=None, serialization_alias="deviceId")
    app_version: str | None = Field(default=None, serialization_alias="appVersion")
    app_store: str | None = Field(default=None, serialization_alias="appStore")


class Health(BaseModel):
    id: UUID
    api_name: str = Field(serialization_alias="apiName")
    api_version: str = Field(serialization_alias="apiVersion") 
    node_ip: str = Field(serialization_alias="nodeIp")
    node_name: str = Field(serialization_alias="nodeName")
    env_name: str = Field(serialization_alias="envName")
    started: datetime | None = None
    pulse: datetime | None = None
    stopped: datetime | None = None
    
    model_config = {"populate_by_name": True, "use_enum_values": True}


