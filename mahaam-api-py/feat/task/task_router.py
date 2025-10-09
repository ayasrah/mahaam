from uuid import UUID
from fastapi import Depends, Form, Body, Query, Path
from fastapi.responses import Response, JSONResponse
from feat.task.task_model import Task
from infra import http
from infra.log import Log
from typing import Annotated, Protocol
from infra.validation import ProtocolEnforcer, Rule
from fastapi import APIRouter
from feat.task.task_service import TaskService
from fastapi_utils.cbv import cbv

class TaskRouter(Protocol):
	def create(self, plan_id: UUID = Path(...), title: str = Form(...)) -> Response: ...
	def delete(self, plan_id: UUID = Path(...), id: UUID = Path(...)) -> Response: ...
	def update_done(self, plan_id: UUID = Path(...), id: UUID = Path(...), done: bool = Form(...)) -> Response: ...
	def update_title(self, id: UUID = Path(...), title: str = Form(...)) -> Response: ...
	def reorder(self, plan_id: UUID = Path(...), old_order: int = Form(...), new_order: int = Form(...)) -> Response: ...
	def select_many(self, plan_id: UUID = Path(...)) -> Response: ...

router = APIRouter(prefix="/plans/{plan_id}", tags=["Tasks"])

def get_task_service() -> TaskService:
	from infra.factory import App
	return App.task_service

@cbv(router)
class DefaultTaskRouter(metaclass=ProtocolEnforcer, protocol=TaskRouter):
	def __init__(self, task_service: TaskService = Depends(get_task_service)):
		self.task_service = task_service
	
	@router.post("/tasks")
	def create(self, plan_id: UUID = Path(...), title: str = Form(...)) -> Response:
		Rule.required(plan_id, "planId")
		Rule.required(title, "title")
		id = self.task_service.create(plan_id, title)
		return JSONResponse(status_code=http.CREATED, content=id)

	@router.delete("/tasks/{id}")
	def delete(self, plan_id: UUID = Path(...), id: UUID = Path(...)) -> Response:
		Rule.required(plan_id, "planId")
		Rule.required(id, "id")
		self.task_service.delete(plan_id, id)
		return Response(status_code=http.NO_CONTENT)

	@router.patch("/tasks/{id}/done")
	def update_done(self, plan_id: UUID = Path(...), id: UUID = Path(...), done: bool = Form(...)) -> Response:
		Rule.required(plan_id, "planId")
		Rule.required(id, "id")
		Rule.required(done, "done")
		self.task_service.update_done(plan_id, id, done)
		return Response(status_code=http.OK)

	@router.patch("/tasks/{id}/title")
	def update_title(self, id: UUID = Path(...), title: str = Form(...)) -> Response:
		Rule.required(id, "id")
		Rule.required(title, "title")
		self.task_service.update_title(id, title)
		return Response(status_code=http.OK)

	@router.patch("/tasks/reorder")
	def reorder(self, plan_id: UUID = Path(...), oldOrder: int = Form(...), newOrder: int = Form(...)) -> Response:
		Rule.required(plan_id, "planId")
		Rule.required(oldOrder, "oldOrder")
		Rule.required(newOrder, "newOrder")
		self.task_service.reorder(plan_id, oldOrder, newOrder)
		return Response(status_code=http.OK)

	@router.get("/tasks", response_model=list[Task], response_model_exclude_none=True)
	# @router.get("")
	def select_many(self, plan_id: UUID = Path(...)) -> Response:
		Rule.required(plan_id, "planId")
		tasks = self.task_service.select_many(plan_id)
		return tasks
