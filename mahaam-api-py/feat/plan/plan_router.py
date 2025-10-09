import json
from uuid import UUID
from fastapi import Depends, Form, Body, Query, Path
from fastapi.responses import Response, JSONResponse
from infra import http
from infra.log import Log
from typing import Annotated, Protocol
from infra.validation import ProtocolEnforcer, Rule
from fastapi import APIRouter
from feat.plan.plan_model import Plan, PlanIn, PlanType
from feat.plan.plan_service import PlanService
from fastapi_utils.cbv import cbv

class PlanRouter(Protocol):
    def create(self, plan: PlanIn = Body(...)) -> Response: ...
    def update(self, plan: PlanIn = Body(...)) -> Response: ...
    def delete(self, id: UUID = Path(...)) -> Response: ...
    def share(self, id: UUID = Path(...), email: str = Form(...)) -> Response: ...
    def unshare(self, id: UUID = Path(...), email: str = Form(...)) -> Response: ...
    def leave(self, id: UUID = Path(...)) -> Response: ...
    def update_type(self, id: UUID = Path(...), type: str = Form(...)) -> Response: ...
    def reorder(self, type: str = Form(...), old_order: int = Form(...), new_order: int = Form(...)) -> Response: ...
    def select_one(self, plan_id: UUID = Path(...)) -> Response: ...
    def select_many(self, type: str | None = Query(None)) -> Response: ...

router = APIRouter(tags=["Plans"])

def get_plan_service() -> PlanService:
    from infra.factory import App
    return App.plan_service

@cbv(router)
class DefaultPlanRouter(metaclass=ProtocolEnforcer, protocol=PlanRouter):
    def __init__(self, plan_service: PlanService = Depends(get_plan_service)):
        self.plan_service = plan_service
    
    @router.post("/plans")
    def create(self, plan: PlanIn = Body(...)) -> Response:
        Rule.one_at_least_required([plan.title, plan.starts, plan.ends], "title or starts or ends is required")
        id = self.plan_service.create(plan)  # Uses self.service
        return JSONResponse(status_code=http.CREATED, content=id)

    @router.put("/plans")
    def update(self, plan: PlanIn = Body(...)) -> Response:
        Rule.required(plan.id, "Id")
        Rule.one_at_least_required([plan.title, plan.starts, plan.ends], "title or starts or ends is required")
        self.plan_service.update(plan)
        return Response(status_code=http.OK)

    @router.delete("/plans/{id}")
    def delete(self, id: UUID = Path(...)) -> Response:
        Rule.required(id, "id")
        self.plan_service.delete(id)
        return Response(status_code=http.NO_CONTENT)

    @router.patch("/plans/{id}/share")
    def share(self, id: UUID = Path(...), email: str = Form(...)) -> Response:
        Rule.required(id, "id")
        Rule.required(email, "email")
        # Service logic handles user lookup and sharing
        self.plan_service.share(id, email)
        return Response(status_code=http.OK)

    @router.patch("/plans/{id}/unshare")
    def unshare(self, id: UUID = Path(...), email: str = Form(...)) -> Response:
        Rule.required(id, "id")
        Rule.required(email, "email")
        self.plan_service.unshare(id, email)
        return Response(status_code=http.OK)

    @router.patch("/plans/{id}/leave")
    def leave(self, id: UUID = Path(...)) -> Response:
        Rule.required(id, "id")
        self.plan_service.leave(id)
        return Response(status_code=http.OK)

    @router.patch("/plans/{id}/type")
    def update_type(self, id: UUID = Path(...), type: str = Form(...)) -> Response:
        Rule.required(id, "id")
        Rule.required(type, "type")
        Rule.contains(PlanType.ALL, type)
        self.plan_service.update_type(id, type)
        return Response(status_code=http.OK)

    @router.patch("/plans/reorder")
    def reorder(self, type: str = Form(...), oldOrder: int = Form(...), newOrder: int = Form(...)) -> Response:
        Rule.required(type, "type")
        Rule.contains(PlanType.ALL, type)
        Rule.required(oldOrder, "oldOrder")
        Rule.required(newOrder, "newOrder")
        self.plan_service.reorder(type, oldOrder, newOrder)
        return Response(status_code=http.OK)

    @router.get("/plans/{plan_id}", response_model=Plan, response_model_exclude_none=True)
    def select_one(self, plan_id: UUID = Path(...)) -> Response:
        Rule.required(plan_id, "planId")
        plan = self.plan_service.select_one(plan_id)
        
        return plan

    @router.get("/plans", response_model=list[Plan], response_model_exclude_none=True)
    # @router.get("")
    def select_many(self, type: str | None = Query(None)) -> Response:
        if type is not None:
            Rule.contains(PlanType.ALL, type)
        else:
            type = PlanType.MAIN
        plans = self.plan_service.select_many(type)
        return plans
