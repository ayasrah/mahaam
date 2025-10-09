import uuid
from fastapi import APIRouter, Form, status
from fastapi.responses import Response
from infra import http
from infra.log import Log
from typing import Annotated, Protocol
from infra.validation import ProtocolEnforcer
from fastapi_utils.cbv import cbv

class AuditRouter(Protocol):
    def error(self, error: str = Form(...)) -> Response: ...
    def info(self, info: str = Form(...)) -> Response: ...

router = APIRouter(tags=["Audit"]) 

@cbv(router)
class DefaultAuditRouter(metaclass=ProtocolEnforcer, protocol=AuditRouter):
	@router.post("/audit/error")
	def error(self, error: str = Form(...)) -> Response:
		Log.error("mahaam-mb: " + error)
		return Response(status_code=http.CREATED)

	@router.post("/audit/info")
	def info(self, info: str = Form(...)) -> Response:
		Log.info("mahaam-mb: " + info)
		return Response(status_code=http.CREATED)