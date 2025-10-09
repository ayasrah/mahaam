from uuid import UUID
from typing import Protocol, List
from contextlib import contextmanager
from feat.plan.plan_repo import PlanRepo
from feat.task.task_repo import TaskRepo
from feat.task.task_model import Task
from infra import db
from infra.exceptions import InputException, LogicException
from infra.validation import ProtocolEnforcer

class TaskService(Protocol):
	def create(self, plan_id: UUID, title: str) -> UUID: ...
	def select_many(self, plan_id: UUID) -> List[Task]: ...
	def delete(self, plan_id: UUID, id: UUID) -> None: ...
	def update_done(self, plan_id: UUID, id: UUID, done: bool) -> None: ...
	def update_title(self, id: UUID, title: str) -> None: ...
	def reorder(self, plan_id: UUID, old_order: int, new_order: int, conn=None) -> None: ...
	
class DefaultTaskService(metaclass=ProtocolEnforcer, protocol=TaskService):
	def __init__(self, task_repo: TaskRepo, plan_repo: PlanRepo) -> None:
		self.task_repo = task_repo
		self.plan_repo = plan_repo

	def create(self, plan_id: UUID, title: str) -> UUID:
		with db.DB.transaction_scope() as conn:
			count = self.task_repo.get_count(plan_id, conn)
			if count >= 100:
				raise LogicException("max_is_100", "Max is 100")
			id = self.task_repo.create(plan_id, title, conn)
			self.plan_repo.update_done_percent(plan_id, conn)
		return id

	def select_many(self, plan_id: UUID) -> List[Task]:
		return self.task_repo.select_many(plan_id)

	def delete(self, plan_id: UUID, id: UUID) -> None:
		with db.DB.transaction_scope() as conn:
			self.task_repo.update_order_before_delete(plan_id, id, conn)
			self.task_repo.delete_one(id, conn)
			self.plan_repo.update_done_percent(plan_id, conn)

	def update_done(self, plan_id: UUID, id: UUID, done: bool) -> None:
		with db.DB.transaction_scope() as conn:
			self.task_repo.update_done(id, done, conn)
			self.plan_repo.update_done_percent(plan_id, conn)
			count = self.task_repo.get_count(plan_id, conn)
			task = self.task_repo.select_one(id, conn)
			self.reorder(plan_id, task.sort_order, 0 if done else count - 1, conn)

	def update_title(self, id: UUID, title: str) -> None:
		self.task_repo.update_title(id, title)

	def reorder(self, plan_id: UUID, old_order: int, new_order: int, conn=None) -> None:
		if old_order == new_order:
			return
		count = self.task_repo.get_count(plan_id, conn)
		if old_order > count or new_order > count:
			raise InputException(f"oldOrder and newOrder should be less than {count}")

		self.task_repo.update_order(plan_id, old_order, new_order, conn)
