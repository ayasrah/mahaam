from typing import Protocol
import uuid
from infra import db
from infra.validation import ProtocolEnforcer
from infra.monitor.monitor_models import Health


class HealthRepo(Protocol):
    def create(self, health: Health) -> int: ...
    def update_pulse(self, id: uuid.UUID) -> int: ...
    def update_stopped(self, id: uuid.UUID) -> int: ...


class DefaultHealthRepo(metaclass=ProtocolEnforcer, protocol=HealthRepo):

    def create(self, health: Health) -> int:
        health_data = {
            "id": str(health.id),
            "api_name": health.api_name,
            "api_version": health.api_version,
            "env_name": health.env_name,
            "node_ip": health.node_ip,
            "node_name": health.node_name
        }

        sql = """
            INSERT INTO x_health (id, api_name, api_version, env_name, node_ip, node_name, started_at) 
            VALUES(:id, :api_name, :api_version, :env_name, :node_ip, :node_name, current_timestamp)"""

        return db.DB.insert(sql, health_data)

    def update_pulse(self, id: uuid.UUID) -> int:
        sql = "UPDATE x_health SET pulsed_at = current_timestamp WHERE id = :id"
        return db.DB.update(sql, {"id": str(id)})

    def update_stopped(self, id: uuid.UUID) -> int:
        sql = "UPDATE x_health SET stopped_at = current_timestamp WHERE id = :id"
        return db.DB.update(sql, {"id": str(id)})
