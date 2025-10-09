from typing import Protocol
from typing import List
from uuid import UUID
from infra import db
from feat.plan.plan_model import Plan
from feat.user.user_model import User
from infra.validation import ProtocolEnforcer


class PlanMembersRepo(Protocol):
    def create(self, plan_id: UUID, user_id: UUID) -> None: ...
    def delete(self, plan_id: UUID, user_id: UUID) -> int: ...
    def get_other_plans(self, user_id: UUID) -> List[Plan]: ...
    def get_users(self, plan_id: UUID) -> List[User]: ...
    def get_plans_count(self, user_id: UUID) -> int: ...
    def get_users_count(self, plan_id: UUID) -> int: ...


class DefaultPlanMembersRepo(metaclass=ProtocolEnforcer, protocol=PlanMembersRepo):
    def create(self, plan_id: UUID, user_id: UUID) -> None:
        sql = """
            INSERT INTO plan_members(plan_id, user_id, created_at) 
            VALUES (:plan_id, :user_id, current_timestamp)"""
        param = {'plan_id': str(plan_id), 'user_id': str(user_id)}
        db.DB.insert(sql, param)

    def delete(self, plan_id: UUID, user_id: UUID) -> int:
        sql = """ 
			DELETE FROM plan_members
            WHERE plan_id = :plan_id AND user_id = :user_id"""
        param = {'plan_id': str(plan_id), 'user_id': str(user_id)}
        return db.DB.delete(sql, param)

    def get_other_plans(self, user_id: UUID) -> List[Plan]:
        sql = """
            SELECT c.id, c.title, c.starts, c.ends, c.type, c.done_percent, c.sort_order,
                u.id user_id, u.email user_email, u.name user_name
            FROM plan_members cm
            LEFT JOIN plans c ON cm.plan_id = c.id
            LEFT JOIN users u ON c.user_id = u.id
            WHERE cm.user_id = :user_id
            ORDER BY c.created_at DESC
        """
        rows = db.DB.select_many(dict, sql, {'user_id': str(user_id)})
        plans = []
        for row in rows:
            plan = Plan(
                id=row['id'],
                title=row['title'],
                starts=row['starts'],
                ends=row['ends'],
                type=row['type'],
                done_percent=row['done_percent'],
                sort_order=row['sort_order'],
                user=User(id=row['user_id'], email=row['user_email'], name=row['user_name']),
                is_shared=True
            )
            plans.append(plan)
        return plans

    def get_users(self, plan_id: UUID) -> List[User]:
        sql = """
            SELECT u.id, u.email, u.name
            FROM plan_members cm
            LEFT JOIN users u ON cm.user_id = u.id
            WHERE cm.plan_id = :plan_id
            ORDER BY u.created_at DESC
        """
        return db.DB.select_many(User, sql, {'plan_id': str(plan_id)})

    def get_plans_count(self, user_id: UUID) -> int:
        sql = """
            SELECT COUNT(1) as count FROM plan_members cm
            LEFT JOIN plans c ON cm.plan_id = c.id
            WHERE c.user_id = :user_id"""
        result = db.DB.select_one(dict, sql, {'user_id': str(user_id)})
        return result["count"]

    def get_users_count(self, plan_id: UUID) -> int:
        sql = "SELECT COUNT(1) as count FROM plan_members WHERE plan_id = :plan_id"
        result = db.DB.select_one(dict, sql, {'plan_id': str(plan_id)})
        return result["count"]
