from typing import Protocol
from uuid import UUID, uuid4
from typing import List
from infra.db import DB
from feat.user.user_model import Device
from infra.validation import ProtocolEnforcer

class DeviceRepo(Protocol):
    def create(self, device: Device, conn) -> UUID: ...
    def delete(self, id: UUID, conn=None) -> int: ...
    def delete_by_user(self, user_id: UUID, except_device_id: UUID) -> int: ...
    def delete_by_fingerprint(self, fingerprint: str, conn) -> int: ...
    def select_one(self, id: UUID) -> Device | None: ...
    def select_many(self, user_id: UUID, conn=None) -> List[Device]: ...
    def update_user_id(self, device_id: UUID, user_id: UUID, conn) -> int: ...

class DefaultDeviceRepo(metaclass=ProtocolEnforcer, protocol=DeviceRepo):
    def create(self, device: Device, conn) -> UUID:
        device_id = uuid4()
        sql = """
        INSERT INTO devices (id, user_id, platform, fingerprint, info, created_at) 
        VALUES (:id, :user_id, :platform, :fingerprint, :info, current_timestamp)"""
        params = {
            "id": str(device_id),
            "user_id": str(device.user_id),
            "platform": device.platform,
            "fingerprint": device.fingerprint,
            "info": device.info
        }
        DB.insert(sql, params, conn)
        return device_id

    def delete(self, id: UUID, conn=None) -> int:
        sql = "DELETE FROM devices WHERE id = :id"
        return DB.delete(sql, {"id": str(id)}, conn)

    def delete_by_user(self, user_id: UUID, except_device_id: UUID) -> int:
        sql = "DELETE FROM devices WHERE user_id = :user_id AND id != :except_device_id"
        return DB.delete(sql, {"user_id": str(user_id), "except_device_id": str(except_device_id)})

    def delete_by_fingerprint(self, fingerprint: str, conn) -> int:
        sql = "DELETE FROM devices WHERE fingerprint = :fingerprint"
        return DB.delete(sql, {"fingerprint": fingerprint}, conn)

    def select_one(self, id: UUID) -> Device | None:
        sql = """
            SELECT id, user_id, platform, fingerprint, info, created_at
            FROM devices WHERE id = :id ORDER BY created_at DESC"""
        return DB.select_one(Device, sql, {"id": str(id)})

    def select_many(self, user_id: UUID, conn=None) -> List[Device]:
        sql = """
            SELECT id, user_id, platform, fingerprint, info, created_at
            FROM devices WHERE user_id = :user_id ORDER BY created_at DESC"""
        return DB.select_many(Device, sql, {"user_id": str(user_id)}, conn)

    def update_user_id(self, device_id: UUID, user_id: UUID, conn) -> int:
        sql = "UPDATE devices SET user_id = :user_id, updated_at = current_timestamp WHERE id = :device_id"
        return DB.update(sql, {"device_id": str(device_id), "user_id": str(user_id)}, conn)
