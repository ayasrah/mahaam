from uuid import UUID, uuid4
from typing import Protocol

from infra import db
from infra.req import req
from infra.log import Log
from feat.plan.plan_model import Plan, PlanIn, PlanType
from feat.user.user_model import User
from infra.validation import ProtocolEnforcer


class PlanRepo(Protocol):
	def select_one(self, id: UUID) -> Plan | None: ...
	def select_many(self, user_id: UUID, type: str) -> list[Plan]: ...
	def create(self, plan: PlanIn, conn=None) -> UUID: ...
	def update(self, plan: PlanIn) -> None: ...
	def delete(self, id: UUID, conn=None) -> None: ...
	def update_done_percent(self, id: UUID, conn) -> None: ...
	def remove_from_order(self, user_id: UUID, id: UUID, conn=None) -> None: ...
	def update_order(self, user_id: UUID, type: str,
					 old_order: int, new_order: int, conn=None) -> None: ...
	def update_type(self, user_id: UUID, id: UUID, type: str, conn=None) -> None: ...
	def get_count(self, user_id: UUID, type: str, conn=None) -> int: ...
	def update_user_id(self, old_user_id: UUID, new_user_id: UUID, conn=None) -> int: ...


class DefaultPlanRepo(metaclass=ProtocolEnforcer, protocol=PlanRepo):
	def create(self, plan: PlanIn, conn=None) -> UUID:
		query = """INSERT INTO plans (id, user_id, title, starts, ends, type, status, done_percent, sort_order, created_at)
			VALUES (:id, :user_id, :title, :starts, :ends, :type, :status, '0/0', 
			(SELECT COUNT(1) FROM plans WHERE user_id = :user_id AND type = :type), 
			current_timestamp)"""
		id = uuid4()
		db.DB.insert(query, {
			"id": str(id),
			"title": plan.title,
			"starts": plan.starts,
			"ends": plan.ends,
			"user_id": req.user_id,
			"type": PlanType.MAIN,
			"status": "Open"
		}, conn)
		return id

	def update(self, plan: PlanIn) -> None:
		query = "UPDATE plans SET title = :title, starts = :starts, ends = :ends, updated_at = current_timestamp WHERE id = :id"
		db.DB.update(query, {"id": str(plan.id), "title": plan.title, "starts": plan.starts, "ends": plan.ends})

	def select_one(self, id: UUID) -> Plan | None:
		query = """
			SELECT c.id, c.title, c.starts, c.ends, c.type, c.done_percent, c.sort_order,
				EXISTS(SELECT 1 FROM plan_members cm WHERE cm.plan_id = c.id) AS is_shared,
				u.id user_id, u.email user_email, u.name user_name
			FROM plans c
			LEFT JOIN users u ON c.user_id = u.id
			WHERE c.id = :id"""

		coll = db.DB.select_one(dict, query, {"id": str(id)})
		if coll is None:
			return None
		
		plan = Plan(
			id=coll["id"],
			title=coll["title"],
			starts=coll["starts"],
			ends=coll["ends"],
			type=coll["type"],
			done_percent=coll["done_percent"],
			sort_order=coll["sort_order"],
			is_shared=coll["is_shared"],
			user=User(id=coll["user_id"], email=coll["user_email"], name=coll["user_name"])
		)
		return plan


	def select_many(self, user_id: UUID, type: str) -> list[Plan]:
		try:
			query = """
			SELECT c.id, c.title, c.starts, c.ends, c.type, c.done_percent, c.sort_order,
				EXISTS(SELECT 1 FROM plan_members cm WHERE cm.plan_id = c.id) AS is_shared,
				u.id user_id, u.email user_email, u.name user_name
			FROM plans c
			LEFT JOIN users u ON c.user_id = u.id
			WHERE c.user_id = :user_id AND c.type = :type
			ORDER BY c.sort_order DESC"""

			dic = db.DB.select_many(dict, query, {"user_id": str(user_id), "type": type})
			plans = []
			for row in dic:
				plan = Plan(
					id=row["id"],
					title=row["title"],
					starts=row["starts"],
					ends=row["ends"],
					type=row["type"],
					done_percent=row["done_percent"],
					sort_order=row["sort_order"],
					is_shared=row["is_shared"],
					user=User(id=row["user_id"], email=row["user_email"],name=row["user_name"])
				)
				plans.append(plan)
			return plans
		except Exception as e:
			Log.error(str(e))
			return []

	def delete(self, id: UUID, conn=None) -> None:
		count = db.DB.delete(
			"DELETE FROM plans WHERE id = :id", {"id": str(id)}, conn)
		if count > 0:
			Log.info(f"Plan {id} deleted")

	def update_done_percent(self, id: UUID, conn) -> None:
		query = "SELECT * FROM tasks WHERE plan_id = :id"
		tasks = db.DB.select_many(dict, query, {"id": str(id)}, conn)

		done = sum(1 for task in tasks if task["done"])
		not_done = len(tasks)
		done_percent = f"{done}/{not_done}"
		update_query = "UPDATE plans SET done_percent = :donePercent WHERE id = :id"
		db.DB.update(update_query, {"donePercent": done_percent, "id": str(id)}, conn)

	def remove_from_order(self, user_id: UUID, id: UUID, conn=None) -> None:
		query = """
			UPDATE plans SET sort_order = sort_order - 1 
			WHERE user_id = :user_id AND type = (SELECT type FROM Plans WHERE id =:id) 
				AND sort_order > (SELECT sort_order FROM plans WHERE id =:id)"""
		db.DB.update(query, {"user_id": str(user_id), "id": str(id)}, conn)

	def update_order(self, user_id: UUID, type: str, old_order: int, new_order: int, conn=None) -> None:
		query = """
			UPDATE plans SET sort_order = 
				CASE 
					WHEN sort_order = :old_order THEN :new_order
					WHEN sort_order > :old_order AND sort_order <= :new_order THEN sort_order - 1
					WHEN sort_order >= :new_order AND sort_order < :old_order THEN sort_order + 1
					ELSE sort_order
				END
			WHERE 
				user_id = :user_id AND 
				type = :type"""
		db.DB.update(query, {
			"user_id": str(user_id), 
			"type": type, 
			"old_order": old_order, 
			"new_order": new_order
		}, conn)

	def update_type(self, user_id: UUID, id: UUID, type: str, conn=None) -> None:
		query = """
			UPDATE plans SET type = :type, 
			sort_order = (SELECT COUNT(1) FROM plans WHERE user_id = :user_id AND type = :type),
			updated_at = current_timestamp WHERE id = :id"""
		db.DB.update(query, {"user_id": str(user_id), "id": str(id), "type": type}, conn)

	def get_count(self, user_id: UUID, type: str, conn=None) -> int:
		query_count = "SELECT COUNT(*) FROM plans WHERE user_id = :user_id and type = :type"
		result = db.DB.select_one(dict, query_count, {"user_id": str(user_id), "type": type}, conn)
		
		if result is None:
			return 0
		else:
			return next(iter(result.values()))

	def update_user_id(self, old_user_id: UUID, new_user_id: UUID, conn=None) -> int:
		query = """
		 	UPDATE plans SET user_id = :new_user_id,
			sort_order = (sort_order + (Select count(1) from plans where user_id=:new_user_id)),
			updated_at = current_timestamp 
			WHERE user_id = :old_user_id"""
		return db.DB.update(query, {"old_user_id": str(old_user_id), "new_user_id": str(new_user_id)}, conn)
