import jwt
import uuid
from datetime import datetime, timedelta
from fastapi import Request
from typing import Tuple, Protocol
from infra import configs
from feat.user.device_repo import DeviceRepo
from feat.user.user_repo import UserRepo
from infra.validation import ProtocolEnforcer, Rule
from infra.exceptions import UnauthorizedException, ForbiddenException
from jwt import decode as jwt_decode
from infra.req import Req


class AuthService(Protocol):
	def validate_and_extract_jwt(self, request: Request) -> Tuple[uuid.UUID, uuid.UUID]: ...
	def create_token(self, user_id: str, device_id: str) -> str: ...


class DefaultAuthService(metaclass=ProtocolEnforcer, protocol=AuthService):
	def __init__(self, device_repo: DeviceRepo, user_repo: UserRepo) -> None:
		self.device_repo = device_repo
		self.user_repo = user_repo

	def validate_and_extract_jwt(self, request: Request) -> Tuple[uuid.UUID, uuid.UUID]:
		path = request.url.path
		authorization = request.headers.get("Authorization")
		
		if not authorization:
			raise UnauthorizedException("Authorization header not exists")
		
		if not authorization.startswith("Bearer "):
			raise UnauthorizedException("Invalid Authorization header format")
		
		token_string = authorization[7:]  # Remove 'Bearer ' to get the jwt
		
		JWT.validate(token_string)
		
		# Decode token to extract claims
		token_payload = jwt_decode(token_string, options={"verify_signature": False})
		
		user_id = token_payload.get("userId")
		self.non_empty_uuid(user_id, "userId")
		
		device_id = token_payload.get("deviceId")
		self.non_empty_uuid(device_id, "deviceId")
		
		is_logged_in = False
		try:
			device = self.device_repo.select_one(uuid.UUID(device_id))
			if (device is None or uuid.UUID(user_id) != device.user_id) and path != "/user/logout":
				raise UnauthorizedException("Invalid user info")
			
			user = self.user_repo.select_one(uuid.UUID(user_id))
			is_logged_in = user.email is not None
		except ValueError:
			raise UnauthorizedException("Invalid UUID format")
		
		return (uuid.UUID(user_id), uuid.UUID(device_id), is_logged_in)

	def create_token(self, user_id: str, device_id: str) -> str:
		return JWT.create(user_id, device_id)

	
	def non_empty_uuid(self,uuid_string: str | None, name: str) -> None:
		EMPTY_UUID = "00000000-0000-0000-0000-000000000000"
		Rule.required(uuid_string, name)
		if not uuid_string or uuid_string.strip() == "" or uuid_string == EMPTY_UUID:
			raise ForbiddenException(f"{name} is Empty")
		
		# Validate UUID format
		try:
			uuid.UUID(uuid_string)
		except ValueError:
			raise UnauthorizedException(f"{name} is not a valid UUID")


# Keep the original Auth class for backward compatibility
class Auth:
	@staticmethod
	def validate_and_extract_jwt(request: Request) -> Tuple[uuid.UUID, uuid.UUID]:
		from infra.factory import App
		return App.auth_service.validate_and_extract_jwt(request)

	@staticmethod
	def create_token(user_id: str, device_id: str) -> str:
		from infra.factory import App
		return App.auth_service.create_token(user_id, device_id)


class JWT:
	@staticmethod
	def create(user_id: str, device_id: str) -> str:
		try:
			payload = {
				"userId": user_id,
				"deviceId": device_id,
				"exp": datetime.utcnow() + timedelta(days=7),
				"iat": datetime.utcnow(),
				"iss": "mahaam-api"
			}
			token = jwt.encode(payload, JWT._security_key(), algorithm="HS256")
			return token
		except Exception as e:
			raise UnauthorizedException(f"Error creating JWT token: {str(e)}")

	@staticmethod
	def validate(token: str) -> None:
		try:
			validation_params = JWT._get_validation_params()
			jwt_decode(token, JWT._security_key(), algorithms=["HS256"], **validation_params)
		except Exception as e:
			raise UnauthorizedException(f"JWT validation failed: {str(e)}")

	@staticmethod
	def _get_validation_params() -> dict:
		return {
			"verify_exp": True,
			"verify_aud": False,
			"verify_iss": True,
			"issuer": "mahaam-api"
		}

	@staticmethod
	def _security_key() -> str:
		return configs.data.tokenSecretKey
