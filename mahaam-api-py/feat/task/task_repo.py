from uuid import UUID, uuid4
from infra import db
from feat.task.task_model import Task
from typing import Protocol
from infra.validation import ProtocolEnforcer
from infra.exceptions import InputException

class TaskRepo(Protocol):
    def select_many(self, plan_id: UUID, conn=None) -> list[Task]: ...
    def select_one(self, id: UUID, conn=None) -> Task | None: ...
    def create(self, plan_id: UUID, title: str, conn=None) -> UUID: ...
    def delete_one(self, id: UUID, conn=None) -> None: ...
    def update_done(self, id: UUID, done: bool, conn=None) -> None: ...
    def update_title(self, id: UUID, title: str) -> None: ...
    def update_order(self, plan_id: UUID, old_order: int, new_order: int, conn=None) -> None: ...
    def update_order_before_delete(self, plan_id: UUID, id: UUID, conn=None) -> None: ...
    def get_count(self, plan_id: UUID, conn=None) -> int: ...


class DefaultTaskRepo(metaclass=ProtocolEnforcer, protocol=TaskRepo):
    def select_many(self, plan_id: UUID, conn=None) -> list[Task]:
        sql = """
        SELECT id, plan_id, title, done, sort_order, created_at, updated_at 
        FROM tasks WHERE plan_id = :plan_id ORDER BY sort_order DESC"""
        return db.DB.select_many(Task, sql, {"plan_id": str(plan_id)}, conn)

    def select_one(self, id: UUID, conn=None) -> Task | None:
        sql = """
        SELECT id, plan_id, title, done, sort_order, created_at, updated_at 
        FROM tasks WHERE id = :id"""
        return db.DB.select_one(Task, sql, {"id": str(id)}, conn)

    def create(self, plan_id: UUID, title: str, conn=None) -> UUID:
        id = uuid4()
        sql = """
        INSERT INTO tasks (id, plan_id, title, done, sort_order, created_at) 
        VALUES (:id, :plan_id, :title, :done, 
        (SELECT COUNT(1) FROM tasks WHERE plan_id = :plan_id), 
        current_timestamp)"""

        task_data = {
            "id": str(id),
            "plan_id": str(plan_id),
            "title": title,
            "done": False
        }

        db.DB.insert(sql, task_data, conn)
        return id

    def delete_one(self, id: UUID, conn=None) -> None:
        sql = "DELETE FROM tasks WHERE id = :id"

        deleted_rows = db.DB.delete(sql, {"id": str(id)}, conn)
        if deleted_rows == 0:
            raise InputException(f"Task id={id} not found")

    def update_done(self, id: UUID, done: bool, conn=None) -> None:
        sql = "UPDATE tasks SET done = :done, updated_at = current_timestamp WHERE id = :id"

        updated_rows = db.DB.update(sql, {"id": str(id), "done": done}, conn)
        if updated_rows == 0:
            raise InputException(f"Task id={id} not found")

    def update_title(self, id: UUID, title: str) -> None:
        sql = "UPDATE tasks SET title = :title, updated_at = current_timestamp WHERE id = :id"

        updated_rows = db.DB.update(sql, {"id": str(id), "title": title})
        if updated_rows == 0:
            raise InputException(f"Task id={id} not found")

    def update_order_before_delete(self, plan_id: UUID, id: UUID, conn=None) -> None:
        sql = """
        UPDATE tasks SET sort_order = sort_order - 1 
        WHERE plan_id = :plan_id AND sort_order > (SELECT sort_order FROM tasks WHERE id = :id)"""

        db.DB.update(sql, {"plan_id": str(plan_id), "id": str(id)}, conn)

    def update_order(self, plan_id: UUID, old_order: int, new_order: int, conn=None) -> None:
        sql = """
        UPDATE tasks SET sort_order = 
            CASE 
                WHEN sort_order = :old_order THEN :new_order 
                WHEN sort_order > :old_order AND sort_order <= :new_order THEN sort_order - 1 
                WHEN sort_order >= :new_order AND sort_order < :old_order THEN sort_order + 1 
                ELSE sort_order 
            END 
        WHERE plan_id = :plan_id"""

        db.DB.update(sql, {
            "plan_id": str(plan_id),
            "old_order": old_order,
            "new_order": new_order
        }, conn)

    def get_count(self, plan_id: UUID, conn=None) -> int:
        sql_count = "SELECT COUNT(1) count FROM tasks WHERE plan_id = :plan_id"
        result = db.DB.select_one(dict, sql_count, {"plan_id": str(plan_id)}, conn)
        return result["count"]
