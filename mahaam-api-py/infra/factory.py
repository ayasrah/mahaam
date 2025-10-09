
from pydantic import BaseModel

# Repos
from feat.plan.plan_repo import DefaultPlanRepo, PlanRepo
from feat.plan.plan_members_repo import DefaultPlanMembersRepo, PlanMembersRepo
from infra.monitor.health_repo import DefaultHealthRepo, HealthRepo
from infra.monitor.log_repo import DefaultLogRepo, LogRepo
from infra.monitor.traffic_repo import DefaultTrafficRepo, TrafficRepo
from feat.task.task_repo import DefaultTaskRepo, TaskRepo
from feat.user.user_repo import DefaultUserRepo, UserRepo
from feat.user.device_repo import DefaultDeviceRepo, DeviceRepo
from feat.user.suggested_emails_repo import DefaultSuggestedEmailsRepo, SuggestedEmailsRepo

# Services
from infra.security import DefaultAuthService, AuthService
from feat.plan.plan_service import DefaultPlanService, PlanService
from infra.monitor.health_service import DefaultHealthService, HealthService
from feat.task.task_service import DefaultTaskService, TaskService
from feat.user.user_service import DefaultUserService, UserService

# Routers
from feat.plan.plan_router import DefaultPlanRouter, PlanRouter
from infra.monitor.audit_router import DefaultAuditRouter, AuditRouter


# class App(BaseModel):
class App():
    # Repos
    log_repo: LogRepo = DefaultLogRepo()
    health_repo: HealthRepo = DefaultHealthRepo()
    traffic_repo: TrafficRepo = DefaultTrafficRepo()
    plan_repo: PlanRepo = DefaultPlanRepo()
    plan_members_repo: PlanMembersRepo = DefaultPlanMembersRepo()
    task_repo: TaskRepo = DefaultTaskRepo()
    user_repo: UserRepo = DefaultUserRepo()
    device_repo: DeviceRepo = DefaultDeviceRepo()
    suggested_emails_repo: SuggestedEmailsRepo = DefaultSuggestedEmailsRepo()

    # Services
    auth_service: AuthService = DefaultAuthService(device_repo=device_repo, user_repo=user_repo)
    health_service: HealthService = DefaultHealthService(health_repo=health_repo)
    plan_service: PlanService = DefaultPlanService(plan_repo=plan_repo, plan_members_repo=plan_members_repo, user_repo=user_repo, suggested_emails_repo=suggested_emails_repo)
    task_service: TaskService = DefaultTaskService(task_repo=task_repo, plan_repo=plan_repo)
    user_service: UserService = DefaultUserService(user_repo=user_repo, plan_repo=plan_repo, device_repo=device_repo, suggested_emails_repo=suggested_emails_repo)

    # Routers
    plan_router: PlanRouter = DefaultPlanRouter()
    audit_router: AuditRouter = DefaultAuditRouter()
    
