from uuid import UUID
from typing import List, Protocol
from feat.plan.plan_model import Plan, PlanType
from feat.plan.plan_members_repo import PlanMembersRepo
from feat.user.suggested_emails_repo import SuggestedEmailsRepo
from feat.user.user_repo import UserRepo
from infra.exceptions import InputException, LogicException, UnauthorizedException
from infra.log import Log
from infra.validation import ProtocolEnforcer
from feat.plan.plan_model import PlanIn
from feat.plan.plan_repo import PlanRepo
from infra import db
from infra.req import req


class PlanService(Protocol):
    def select_one(self, plan_id: UUID) -> Plan: ...
    def select_many(self, type: str | None) -> list[Plan]: ...
    def create(self, plan: PlanIn) -> UUID: ...
    def update(self, plan: PlanIn) -> None: ...
    def delete(self, id: UUID) -> None: ...
    def share(self, id: UUID, email: str) -> None: ...
    def unshare(self, id: UUID, email: str) -> None: ...
    def leave(self, id: UUID) -> None: ...
    def update_type(self, id: UUID, type: str) -> None: ...
    def reorder(self, type: str, old_order: int, new_order: int) -> None: ...
    def validate_user_owns_the_plan(self, plan_id: UUID) -> None: ...
    def _validate_user_logged_in(self) -> None: ...


class DefaultPlanService(metaclass=ProtocolEnforcer, protocol=PlanService):
    def __init__(self, plan_repo: PlanRepo, plan_members_repo: PlanMembersRepo, user_repo: UserRepo, suggested_emails_repo: SuggestedEmailsRepo) -> None:
        self.plan_repo = plan_repo  
        self.plan_members_repo = plan_members_repo
        self.user_repo = user_repo
        self.suggested_emails_repo = suggested_emails_repo
        
              
    def select_one(self, plan_id: UUID) -> Plan:
        plan = self.plan_repo.select_one(plan_id)
        if plan and plan.is_shared:
            plan.members = self.plan_members_repo.get_users(plan_id)
        return plan

    def select_many(self, type: str | None) -> list[Plan]:
        # plans of the user shared or not
        plans = self.plan_repo.select_many(req.user_id, type)
        if req.is_logged_in:
            # plans of others shared with the user
            plan_members = self.plan_members_repo.get_other_plans(req.user_id)
            plans.extend(plan_members)
        return plans

    def create(self, plan: PlanIn) -> UUID:
        count = self.plan_repo.get_count(req.user_id, PlanType.MAIN)
        if count >= 100:
            raise LogicException("max_is_100", "Max is 100")

        return self.plan_repo.create(plan)

    def update(self, plan: PlanIn) -> None:
        self.validate_user_owns_the_plan(plan.id)
        self.plan_repo.update(plan)

    def delete(self, id: UUID) -> None:
        self.validate_user_owns_the_plan(id)
        with db.DB.transaction_scope() as conn:
            self.plan_repo.remove_from_order(req.user_id, id, conn)
            # This will delete all related records as it is a cascade delete
            self.plan_repo.delete(id, conn)

    def share(self, id: UUID, email: str) -> None:
        self._validate_user_logged_in()
        self.validate_user_owns_the_plan(id)
        user = self.user_repo.select_one_by_email(email)
        if user is None:
            raise InputException(f"email:{email} was not found")

        if user.id == req.user_id:
            raise LogicException(
                "not_allowed_to_share_with_creator", "Not allowed to share with creator")

        limit = 20
        plan = self.plan_repo.select_one(id)
        if plan and plan.is_shared:
            user_id_count = self.plan_members_repo.get_users_count(id)
            if user_id_count >= limit:
                raise LogicException("max_is_20", "Max is 20")
        else:
            plans_count = self.plan_members_repo.get_plans_count(req.user_id)
            if plans_count >= limit:
                raise LogicException("max_is_20", "Max is 20")

        self.plan_members_repo.create(id, user.id)
        self.suggested_emails_repo.create(req.user_id, email)
        usr = self.user_repo.select_one(req.user_id)
        self.suggested_emails_repo.create(user.id, usr.email)

    def unshare(self, id: UUID, email: str) -> None:
        self._validate_user_logged_in()
        self.validate_user_owns_the_plan(id)
        user = self.user_repo.select_one_by_email(email)
        if user is None:
            raise InputException(f"email:{email} was not found")

        self.plan_members_repo.delete(id, user.id)

    def leave(self, id: UUID) -> None:
        self._validate_user_logged_in()
        deleted_records = self.plan_members_repo.delete(id, req.user_id)
        if deleted_records == 1:
            Log.info(f"user {req.user_id} left plan {id}")
        else:
            raise InputException(f"userId={req.user_id} cannot leave planId={id}")

    def update_type(self, id: UUID, type: str) -> None:
        self.validate_user_owns_the_plan(id)
        count = self.plan_repo.get_count(req.user_id, type)
        if count >= 100:
            raise LogicException("max_is_100", "Max is 100")

        with db.DB.transaction_scope() as conn:
            self.plan_repo.remove_from_order(req.user_id, id, conn)
            self.plan_repo.update_type(req.user_id, id, type, conn)

    def reorder(self, type: str, old_order: int, new_order: int) -> None:
        count = self.plan_repo.get_count(req.user_id, type)
        if old_order > count or new_order > count:
            raise InputException("oldOrder and newOrder should be less than " + count)
        self.plan_repo.update_order(req.user_id, type, old_order, new_order)

    def validate_user_owns_the_plan(self, plan_id: UUID) -> None:
        plan = self.plan_repo.select_one(plan_id)
        if plan is None:
            raise InputException("planId not found")
        if plan.user.id != req.user_id:
            raise UnauthorizedException("User does not own this plan")

    def _validate_user_logged_in(self) -> None:
        user = self.user_repo.select_one(req.user_id)
        if user.email is None:
            raise LogicException("you_are_not_logged_in",
                                 "You are not logged in")


