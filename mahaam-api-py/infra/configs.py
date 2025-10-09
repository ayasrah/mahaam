from pydantic import BaseModel
import json
import logging


class Config(BaseModel):
    apiName: str
    apiVersion: str
    envName: str
    dbUrl: str
    logFile: str
    logFileSizeLimit: int
    logFileCountLimit: int
    logFileOutputTemplate: str
    logFileRollingInterval: str
    httpPort: int
    tokenSecretKey: str
    emailAccountSid: str
    emailVerificationServiceSid: str
    emailAuthToken: str
    testEmails: list[str]
    testSID: str
    testOTP: str
    logReqEnabled: bool


data: Config = None


def init(path: str):
    global data
    try:
        with open(path) as f:
            config_data = json.load(f)
            data = Config(**config_data)
    except Exception as e:
        logging.fatal("Failed to load config: %s", e)
        raise
