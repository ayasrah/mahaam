from uuid import UUID, uuid4
from typing import Protocol
from infra.db import DB
from feat.user.user_model import User
from infra.validation import ProtocolEnforcer

class UserRepo(Protocol):
    def create(self, conn) -> UUID: ...
    def update_name(self, id: UUID, name: str) -> None: ...
    def update_email(self, id: UUID, email: str, conn) -> None: ...
    def select_one_by_email(self, email: str, conn=None) -> User | None: ...
    def select_one(self, id: UUID, conn=None) -> User | None: ...
    def delete(self, id: UUID, conn) -> int: ...

class DefaultUserRepo(metaclass=ProtocolEnforcer, protocol=UserRepo):
    def create(self, conn) -> UUID:
        sql = "INSERT INTO users (id, created_at) VALUES (:id, current_timestamp)"
        id = uuid4()
        DB.insert(sql, {"id": str(id)}, conn)
        return id

    def update_name(self, id: UUID, name: str) -> None:
        sql = "UPDATE users SET name = :name, updated_at = current_timestamp WHERE id = :id"
        DB.update(sql, {"id": str(id), "name": name})

    def update_email(self, id: UUID, email: str, conn) -> None:
        sql = "UPDATE users SET email = :email, updated_at = current_timestamp WHERE id = :id"
        DB.update(sql, {"id": str(id), "email": email}, conn)

    def select_one_by_email(self, email: str, conn=None) -> User | None:
        sql = "SELECT id, name, email FROM users WHERE email = :email"
        return DB.select_one(User, sql, {"email": email}, conn)

    def select_one(self, id: UUID, conn=None) -> User | None:
        sql = "SELECT id, name, email FROM users WHERE id = :id"
        return DB.select_one(User, sql, {"id": str(id)}, conn)

    def delete(self, id: UUID, conn) -> int:
        sql = "DELETE FROM users WHERE id = :id"
        return DB.delete(sql, {"id": str(id)}, conn)
