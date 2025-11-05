# removed unused import
from infra.cache import Cache
from infra import configs
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
            id=Cache.health_id(),
            api_name=configs.data.apiName,
            api_version=configs.data.apiVersion,
            node_ip=Cache.node_ip(),
            node_name=Cache.node_name(),
            env_name=configs.data.envName
        )
    
