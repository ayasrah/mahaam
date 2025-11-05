import uuid
import time
from typing import Callable
from fastapi import Request, Response
from starlette.middleware.base import BaseHTTPMiddleware
from starlette.types import ASGIApp

import json
from infra.cache import Cache
from infra import log, http, configs
from infra.factory import App
from infra.monitor.monitor_models import Traffic, TrafficHeaders
from infra.security import Auth
from infra.exceptions import AppException, NotFoundException, UnauthorizedException
from infra.req import req

async def authenticate_req (request: Request, path: str):
	"""Authenticate request and validate required headers"""
	path_base = request.scope.get("root_path", "")
	
	if path_base != "/mahaam-api":
		raise NotFoundException("Invalid path base")

	app_store = request.headers.get("x-app-store")
	app_version = request.headers.get("x-app-version")

	if not app_store or not app_version:
		raise UnauthorizedException("Required headers not exists")

	req.app_store = app_store
	req.app_version = app_version

	bypass_auth_paths = ["/swagger", "/health", "/users/create", "/audit"]
	if not any(path.startswith(p) for p in bypass_auth_paths):
		user_id, device_id, is_logged_in = Auth.validate_and_extract_jwt(request)
		req.user_id = user_id
		req.device_id = device_id
		req.is_logged_in = is_logged_in

def handle_exception( e: Exception, traffic_id: uuid.UUID) -> tuple[str, int]:
	"""Handle exceptions and return response body and status code"""
	response_status = http.SERVER_ERROR
	res_body = json.dumps(str(e))

	if isinstance(e, AppException):
		response_status = e.http_code
		if e.key:
			res_body = json.dumps(
				{"key": e.key, "error": str(e)})

	log.Log.error(str(e), traffic_id=traffic_id)
	return res_body, response_status


def create_traffic(traffic_id: uuid.UUID, request: Request, path: str, 
					response_status: int, duration: int, req_body: str, res_body: str):
	"""Create and store traffic record"""
	no_traffic_paths = ["/swagger", "/health", "/audit"]
	if any(path.startswith(p) for p in no_traffic_paths):
		return

	try:
		# Don't log response for user endpoints or if empty
		if path.startswith("/user") or (res_body and not res_body.strip()):
			res_body = None
		
		# Don't log request/response for successful requests
		is_success_response = response_status < 400
		if is_success_response:
			req_body = None
			res_body = None
		
		# Create traffic headers
		traffic_headers = TrafficHeaders(
			user_id=req.user_id,
			device_id=req.device_id,
			app_store=req.app_store,
			app_version=req.app_version
		)
		
		traffic = Traffic(
			id=traffic_id,
			health_id=Cache.health_id(),
			method=request.method,
			path=path,
			code=response_status,
			elapsed=duration,
			headers=json.dumps(traffic_headers.model_dump(by_alias=True)),
			request=req_body,
			response=res_body
		)
		# Ideally wrap this in try-catch inside the task
		try:
			App.traffic_repo.create(traffic)
		except Exception as traffic_error:
			log.Log.error("error creating traffic record: " + str(traffic_error), traffic_id=traffic_id)
	except Exception as e:
		log.Log.error("error creating traffic record: " + str(e), traffic_id=traffic_id)


class AppMW(BaseHTTPMiddleware):
	def __init__(self, app: ASGIApp):
		super().__init__(app)

	async def dispatch(self, request: Request, call_next: Callable):
		return await req.run(lambda: self._handle_request(request, call_next))

	async def _handle_request(self, request: Request, call_next):
		start = time.time()
		traffic_id = uuid.uuid4()
		req.traffic_id = str(traffic_id)

		req_body = None
		if configs.data.logReqEnabled:
			req_body = await self.get_payload(request)
		
		path_base = request.scope.get("root_path", "")
		path = request.url.path.replace(path_base, "")

		res_body = None
		response_status = 500
		original_body = b""

		try:
			await authenticate_req(request, path)

			if configs.data.logReqEnabled:
				# Temporarily replace request stream to be able to read again
				request._receive = self._receive_override(request)

			response = await call_next(request)

			body = [section async for section in response.body_iterator]
			original_body = b"".join(body)
			response.body_iterator = iter([original_body])
			res_body = original_body.decode()
			response_status = response.status_code

		except Exception as e:
			res_body, response_status = handle_exception(e, traffic_id)
			if configs.data.logReqEnabled and req_body is None:
				req_body = await self.get_payload(request)
			return Response(content=res_body, status_code=response_status, media_type=http.JSON)

		finally:
			duration = int((time.time() - start) * 1000)
			create_traffic(traffic_id, request, path, response_status, duration, req_body, res_body)

		return Response(content=original_body, status_code=response_status, media_type="application/json")

	async def get_payload(self, request: Request) -> str | None:
		body = await request.body()
		return body.decode("utf-8") if body else None

	def _receive_override(self, request: Request):
		# Make body reusable by resetting the internal receive function
		body = request._body

		async def receive():
			return {"type": "http.request", "body": body}

		return receive

