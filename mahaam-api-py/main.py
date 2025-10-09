import uvicorn
from contextlib import asynccontextmanager
from fastapi import FastAPI, APIRouter
from infra.monitor import audit_router, health_router
from feat.plan import plan_router
from feat.task import task_router
from feat.user import user_router
from infra import log, starter
import infra.configs as configs
from infra.factory import App
from infra.middlewares import AppMW
from json import JSONEncoder
from uuid import UUID

configs.init("config.json")

@asynccontextmanager
async def lifespan(app: FastAPI):
    # Startup
    app_instance = App()
    log.init(app_instance.log_repo.create)
    starter.init(app_instance)    
    yield
    # Shutdown
    app_instance.health_service.server_stopped()

app = FastAPI(root_path="/mahaam-api", lifespan=lifespan)

# Add middlewares to the app
app.add_middleware(AppMW)


router = APIRouter()
router.include_router(audit_router.router)
router.include_router(plan_router.router)
router.include_router(task_router.router)
router.include_router(user_router.router)
router.include_router(health_router.router)
app.include_router(router)


old_default = JSONEncoder.default

def new_default(self, obj):
    if isinstance(obj, UUID):
        return str(obj)
    return old_default(self, obj)

JSONEncoder.default = new_default


if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=configs.data.httpPort)