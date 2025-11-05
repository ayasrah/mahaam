from infra.monitor.monitor_models import Health
from uuid import UUID

class Cache:
    _health: Health | None = None

    @classmethod
    def init(cls, health: Health) -> None:
        cls._health = health

    @classmethod
    def node_ip(cls) -> str:
        return cls._health.node_ip if cls._health else ""

    @classmethod
    def node_name(cls) -> str:
        return cls._health.node_name if cls._health else ""

    @classmethod
    def health_id(cls) -> UUID:
        return cls._health.id if cls._health else UUID(int=0)