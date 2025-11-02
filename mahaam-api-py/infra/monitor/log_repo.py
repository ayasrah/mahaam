import uuid
from concurrent.futures import ThreadPoolExecutor
from infra import cache, db
from typing import Protocol
from infra.validation import ProtocolEnforcer

class LogRepo(Protocol):
    def create(self, traffic_id: uuid.UUID | None, type: str, message: str): ...

class DefaultLogRepo(metaclass=ProtocolEnforcer, protocol=LogRepo):
    def create(self, traffic_id: uuid.UUID | None, type: str, message: str):
        def insert_log():
            try:
                data = {
                    "trafficId": str(traffic_id) if traffic_id else None,
                    "type": type,
                    "message": message,
                    "node_ip": cache.node_ip,
                }

                sql = """
                INSERT INTO monitor.logs (traffic_id, type, message, node_ip, created_at)
                VALUES (:trafficId, :type, :message, :node_ip, current_timestamp)
                """
                db.DB.insert(sql, data)
            except Exception as ex:
                Log.error(f"Unable to create log record:({type}: {message}). Cause: {ex}")

        ThreadPoolExecutor(max_workers=1).submit(insert_log)
