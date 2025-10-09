import contextvars
from typing import Any, Optional

# Store a dict in context var
_request_context: contextvars.ContextVar[dict[str, Any]] = contextvars.ContextVar("_request_context", default={})


class ReqCtx:
    @staticmethod
    def run(callback):
        # Create a new context for the request
        token = _request_context.set({})
        try:
            return callback()
        finally:
            _request_context.reset(token)

    @staticmethod
    def set(name: str, value: Any) -> None:
        ctx = _request_context.get()
        ctx[name] = value

    @staticmethod
    def get(name: str) -> Optional[Any]:
        return _request_context.get().get(name)


class Req:
    @staticmethod
    def run(callback):
        return ReqCtx.run(callback)

    @property
    def traffic_id(self) -> str:
        return ReqCtx.get("trafficId") or ""

    @traffic_id.setter
    def traffic_id(self, value: str) -> None:
        ReqCtx.set("trafficId", value)

    @property
    def user_id(self) -> str:
        return ReqCtx.get("userId") or ""

    @user_id.setter
    def user_id(self, value: str) -> None:
        ReqCtx.set("userId", value)

    @property
    def device_id(self) -> str:
        return ReqCtx.get("deviceId") or ""

    @device_id.setter
    def device_id(self, value: str) -> None:
        ReqCtx.set("deviceId", value)

    @property
    def app_store(self) -> str:
        return ReqCtx.get("appStore") or ""

    @app_store.setter
    def app_store(self, value: str) -> None:
        ReqCtx.set("appStore", value)

    @property
    def app_version(self) -> str:
        return ReqCtx.get("appVersion") or ""

    @app_version.setter
    def app_version(self, value: str) -> None:
        ReqCtx.set("appVersion", value)

    @property
    def is_logged_in(self) -> bool:
        return ReqCtx.get("isLoggedIn") or False

    @is_logged_in.setter
    def is_logged_in(self, value: bool) -> None:
        ReqCtx.set("isLoggedIn", value)


req = Req()
