import socket
import os
import time
import uuid
import re
from infra import cache, email, configs
from infra.log import Log
from infra.db import DB
from infra.monitor.monitor_models import Health
from infra.factory import App


def init(app_instance: App):
    """Initialize the application with proper startup sequence"""
    init_db()
    email.init()
    
    # Create health object
    health = Health(
        id=uuid.uuid4(),
        api_name=configs.data.apiName,
        api_version=configs.data.apiVersion,
        node_ip=get_node_ip(),
        node_name=get_node_name(),
        env_name=configs.data.envName,
        started=None,  # Will be set by health service
        pulse=None,
        stopped=None
    )
    
    # Initialize health service and cache
    app_instance.health_service.server_started(health)
    cache.init(health)
    
    # Log startup message
    start_msg = f"✓ {cache.api_name}-v{cache.api_version}/{cache.node_ip}-{cache.node_name} started with healthID={cache.health_id}"
    Log.info(start_msg)
    time.sleep(2)
    app_instance.health_service.start_sending_pulses()


def init_db():
    """Test database connection and log connection info"""
    try:
        # Test connection by creating engine and connecting
        engine = DB.get_engine()
        with engine.connect() as connection:
            print("connection",connection)
        
        # Extract host from database URL for logging
        pattern = r"://[^@]+@([^:/]+)"
        match = re.search(pattern, configs.data.dbUrl)
        host = match.group(1) if match else "unknown"
        
        Log.info(f"✓ Connected to DB on server {host}")
    except Exception as e:
        Log.error(f"Failed to connect to database: {e}")
        raise


def get_node_ip() -> str:
    """Get the local IP address by connecting to a public IP"""
    try:
        with socket.socket(socket.AF_INET, socket.SOCK_DGRAM) as s:
            s.connect(("8.8.8.8", 10002))
            return s.getsockname()[0]
    except Exception as e:
        Log.error(f"An error occurred while getting the local IP address: {e}")
        return "127.0.0.1"


def get_node_name() -> str:
    """Get the machine name"""
    try:
        return os.uname().nodename
    except Exception as e:
        Log.error(f"Failed to get node name: {e}")
        return "unknown"
