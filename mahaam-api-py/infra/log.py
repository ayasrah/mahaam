import logging

from uuid import UUID
from typing import Optional, Callable, Any
from logging.handlers import RotatingFileHandler

from infra.monitor.log_repo import LogRepo
from infra import configs
from infra.req import req


# Create logger
logger = logging.getLogger(__name__)
logger.setLevel(logging.INFO)

# Define CreateLogFunc type
CreateLogFunc = Callable[[UUID, str, str], None]

# Global variable for log creation function
create_log_func: Optional[CreateLogFunc] = None

# Custom formatter to match Go format
formatter = logging.Formatter(
    '%(asctime)s.%(msecs)03d %(levelname)s %(message)s',
    datefmt='%Y-%m-%d %H:%M:%S'
)


def _get_file_handler():
    handler = RotatingFileHandler(
        configs.data.logFile, 
        maxBytes=configs.data.logFileSizeLimit,
        backupCount=configs.data.logFileCountLimit
    )
    handler.setFormatter(formatter)
    return handler


def _get_console_handler():
    console_handler = logging.StreamHandler()
    console_handler.setFormatter(formatter)
    return console_handler


def init(create_log_fn: CreateLogFunc):
    """Initialize logging with provided log creation function"""
    global create_log_func
    create_log_func = create_log_fn
    
    if getattr(logger, '_is_initialized', False):
        return
        
    # Configure for different environments
    if configs.data.envName == "local":
        logger.setLevel(logging.DEBUG)
    else:
        logger.setLevel(logging.INFO)
    
    logger.addHandler(_get_file_handler())
    logger.addHandler(_get_console_handler())
    logger._is_initialized = True


class Log:
    @staticmethod
    def info(info: str, traffic_id: Optional[UUID] = None):
        current_traffic_id = traffic_id or req.traffic_id
        message = f"TrafficId: {current_traffic_id}, {info}" if current_traffic_id else info
        
        logger.info(message)
        if create_log_func and current_traffic_id:
            create_log_func(UUID(current_traffic_id) if isinstance(current_traffic_id, str) else current_traffic_id, "Info", info)

    @staticmethod
    def error(error: str, traffic_id: Optional[UUID] = None):
        current_traffic_id = traffic_id or req.traffic_id
        message = f"TrafficId: {current_traffic_id}, {error}" if current_traffic_id else error
        
        logger.error(message)
        if create_log_func and current_traffic_id:
            create_log_func(UUID(current_traffic_id) if isinstance(current_traffic_id, str) else current_traffic_id, "Error", error)


