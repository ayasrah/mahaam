import threading
import time
from typing import Protocol
from infra.log import Log
from infra import cache, configs
from infra.validation import ProtocolEnforcer
from infra.monitor.health_repo import HealthRepo
from infra.monitor.monitor_models import Health

class HealthService(Protocol):
    def server_started(self, health: Health) -> None: ...
    def server_stopped(self) -> None: ...
    def start_sending_pulses(self) -> None: ...
    

class DefaultHealthService(metaclass=ProtocolEnforcer, protocol=HealthService):
    def __init__(self, health_repo: HealthRepo) -> None:
        self.health_repo = health_repo
    pulse_thread = None
    pulse_stop_event = threading.Event()

    def server_started(self, health: Health) -> None:
        self.health_repo.create(health)
        # Note: startup message and timing moved to starter.py for consistency with C# version

    def start_sending_pulses(self) -> None:
        def pulse_loop():
            while not DefaultHealthService.pulse_stop_event.is_set():
                try:
                    if cache.health_id:
                        self.health_repo.update_pulse(cache.health_id)
                    time.sleep(60)  # 1 minute
                except Exception as e:
                    Log.error(str(e))
        DefaultHealthService.pulse_thread = threading.Thread(target=pulse_loop, daemon=True)
        DefaultHealthService.pulse_thread.start()

    def server_stopped(self) -> None:
        def stop_thread():
            try:
                if cache.health_id:
                    self.health_repo.update_stopped(cache.health_id)
                    stop_msg = f"✓ {configs.data.apiName}-v{configs.data.apiVersion}/{cache.node_ip}-{cache.node_name} stopped with healthID={cache.health_id}"
                    Log.info(stop_msg)
            except Exception as e:
                Log.error(str(e))
        t = threading.Thread(target=stop_thread)
        t.start()
        DefaultHealthService.pulse_stop_event.set()
        if DefaultHealthService.pulse_thread:
            DefaultHealthService.pulse_thread.join(timeout=2)
