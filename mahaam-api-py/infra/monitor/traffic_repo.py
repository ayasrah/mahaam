from typing import Protocol
from concurrent.futures import ThreadPoolExecutor
from infra import db
from infra.monitor.monitor_models import Traffic
from infra.validation import ProtocolEnforcer
from infra.log import Log


class TrafficRepo(Protocol):
    def create(self, traffic: Traffic) -> None: ...


class DefaultTrafficRepo(metaclass=ProtocolEnforcer, protocol=TrafficRepo):

    def create(self, traffic: Traffic) -> None:
        def insert_traffic():
            try:
                traffic_data = {
                    "id": str(traffic.id),
                    "health_id": str(traffic.health_id),
                    "method": traffic.method,
                    "path": traffic.path,
                    "code": traffic.code,
                    "elapsed": traffic.elapsed,
                    "headers": traffic.headers,
                    "request": traffic.request,
                    "response": traffic.response
                }

                sql = """
                    INSERT INTO x_traffic (id, health_id, method, path, code, elapsed, headers, request, response, created_at) 
                    VALUES(:id, :health_id, :method, :path, :code, :elapsed, :headers, :request, :response, current_timestamp)"""

                db.DB.insert(sql, traffic_data)
            except Exception as e:
                Log.error(f"error creating traffic record: {e}")

        ThreadPoolExecutor(max_workers=1).submit(insert_traffic)
