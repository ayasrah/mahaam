from typing import Protocol
from uuid import UUID, uuid4
from typing import List
from infra.db import DB
from feat.user.user_model import SuggestedEmail
import logging
from infra.validation import ProtocolEnforcer

class SuggestedEmailsRepo(Protocol):
    def create(self, user_id: UUID, email: str) -> UUID: ...
    def delete(self, id: UUID) -> int: ...
    def delete_many_by_email(self, email: str, conn) -> int: ...
    def select_many(self, user_id: UUID) -> List[SuggestedEmail]: ...
    def select_one(self, id: UUID) -> SuggestedEmail | None: ...

class DefaultSuggestedEmailsRepo(metaclass=ProtocolEnforcer, protocol=SuggestedEmailsRepo):
    def create(self, user_id: UUID, email: str) -> UUID:
        sql = """
        INSERT INTO suggested_emails (id, user_id, email, created_at) 
        VALUES (:id, :user_id, :email, current_timestamp)
        ON CONFLICT (user_id, email) DO NOTHING"""
        id = uuid4()
        params = {"id": str(id), "user_id": str(user_id), "email": email}
        updated = DB.insert(sql, params)
        return id if updated > 0 else UUID(int=0)

    def delete(self, id: UUID) -> int:
        sql = "DELETE FROM suggested_emails WHERE id = :id"
        return DB.delete(sql, {"id": str(id)})

    def delete_many_by_email(self, email: str, conn) -> int:
        sql = "DELETE FROM suggested_emails WHERE email = :email"
        return DB.delete(sql, {"email": email}, conn)

    def select_many(self, user_id: UUID) -> List[SuggestedEmail]:
        sql = """
            SELECT id, user_id, email, created_at
            FROM suggested_emails WHERE user_id = :user_id ORDER BY created_at DESC"""
        return DB.select_many(SuggestedEmail, sql, {"user_id": str(user_id)})

    def select_one(self, id: UUID) -> SuggestedEmail | None:
        sql = "SELECT id, user_id, email, created_at FROM suggested_emails WHERE id = :id"
        return DB.select_one(SuggestedEmail, sql, {"id": str(id)})
