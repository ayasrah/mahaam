# removed unused import
from infra import cache
from typing import Protocol
from infra.validation import ProtocolEnforcer
from fastapi import APIRouter
from infra.monitor.monitor_models import Health
from fastapi_utils.cbv import cbv

class HealthRouter(Protocol):
    def get_info(self) -> Health: ...

router = APIRouter(tags=["Health"])

@cbv(router)
class DefaultHealthRouter(metaclass=ProtocolEnforcer, protocol=HealthRouter):
    
    @router.get("/health", response_model=Health, response_model_exclude_none=True)
    def get_info(self) -> Health:
        return Health(
            id=cache.health_id,
            api_name=cache.api_name,
            api_version=cache.api_version,
            node_ip=cache.node_ip,
            node_name=cache.node_name,
            env_name=cache.env_name
        )
    
