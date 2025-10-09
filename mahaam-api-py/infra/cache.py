from infra.monitor.monitor_models import Health
from uuid import UUID

# Simple module-level variables like Go implementation
node_ip = ""
node_name = ""
api_name = ""
api_version = ""
env_name = ""
health_id = UUID(int=0)

def init(health: Health):
    """Initialize cache with health object"""
    global node_ip, node_name, api_name, api_version, env_name, health_id
    node_ip = health.node_ip
    node_name = health.node_name
    api_name = health.api_name
    api_version = health.api_version
    env_name = health.env_name
    health_id = health.id
