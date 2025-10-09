from uuid import UUID
from fastapi import Depends, Form
from fastapi.responses import Response, JSONResponse
from infra import http
from infra.log import Log
from typing import Annotated, Protocol
from infra.validation import ProtocolEnforcer, Rule
from fastapi import APIRouter
from feat.user.user_service import UserService
from feat.user.user_model import Device, CreatedUser, VerifiedUser, SuggestedEmail
from fastapi_utils.cbv import cbv

class UserRouter(Protocol):
    def create(self, platform: str = Form(...), is_physical_device: bool = Form(...), device_fingerprint: str = Form(...), device_info: str = Form(...)) -> Response: ...
    def send_me_otp(self, email: str = Form(...)) -> Response: ...
    def verify_otp(self, email: str = Form(...), sid: str = Form(...), otp: str = Form(...)) -> Response: ...
    def refresh_token(self) -> Response: ...
    def update_name(self, name: str = Form(...)) -> Response: ...
    def logout(self, device_id: UUID = Form(...)) -> Response: ...
    def delete(self, sid: str = Form(...), otp: str = Form(...)) -> Response: ...
    def get_devices(self) -> Response: ...
    def get_suggested_emails(self) -> Response: ...
    def delete_suggested_email(self, suggested_email_id: UUID = Form(...)) -> Response: ...

router = APIRouter(tags=["User"])

def get_user_service() -> UserService:
    from infra.factory import App
    return App.user_service

@cbv(router)
class DefaultUserRouter(metaclass=ProtocolEnforcer, protocol=UserRouter):
    def __init__(self, user_service: UserService = Depends(get_user_service)):
        self.user_service = user_service
    
    @router.post("/users/create", response_model=CreatedUser, response_model_exclude_none=True)
    def create(self, platform: str = Form(...), is_physical_device: bool = Form(..., alias="isPhysicalDevice"), device_fingerprint: str = Form(..., alias="deviceFingerprint"), device_info: str = Form(..., alias="deviceInfo")) -> Response:
        Rule.required(is_physical_device, "isPhysicalDevice")
        Rule.required(platform, "platform")
        Rule.required(device_fingerprint, "deviceFingerprint")
        Rule.required(device_info, "deviceInfo")
        Rule.fail_if(not is_physical_device, "Device should be real not simulator")
        
        device = Device(platform=platform, fingerprint=device_fingerprint, info=device_info)
        created_user = self.user_service.create(device)
        return created_user

    @router.post("/users/send-me-otp")
    def send_me_otp(self, email: str = Form(...)) -> Response:
        Rule.validate_email(email)
        verification_sid = self.user_service.send_me_otp(email)
        return JSONResponse(status_code=http.OK, content=verification_sid)

    @router.post("/users/verify-otp", response_model=VerifiedUser, response_model_exclude_none=True)
    def verify_otp(self, email: str = Form(...), sid: str = Form(...), otp: str = Form(...)) -> Response:
        Rule.required(email, "email")
        Rule.required(sid, "sid")
        Rule.required(otp, "otp")
        verified_user = self.user_service.verify_otp(email, sid, otp)
        return verified_user

    @router.post("/users/refresh-token", response_model=VerifiedUser, response_model_exclude_none=True)
    def refresh_token(self) -> Response:
        verified_user = self.user_service.refresh_token()
        return verified_user

    @router.patch("/users/name", response_model=None)
    def update_name(self, name: str = Form(...)) -> Response:
        Rule.required(name, "name")
        self.user_service.update_name(name)
        return Response(status_code=http.OK)

    @router.post("/users/logout", response_model=None)
    def logout(self, device_id: UUID = Form(..., alias="deviceId")) -> Response:
        Rule.required(device_id, "deviceId")
        self.user_service.logout(device_id)
        return JSONResponse(status_code=http.OK, content={})

    @router.delete("/users")
    def delete(self, sid: str = Form(...), otp: str = Form(...)) -> Response:
        Rule.required(sid, "sid")
        Rule.required(otp, "otp")
        self.user_service.delete(sid, otp)
        return Response(status_code=http.NO_CONTENT)

    @router.get("/users/devices", response_model=list[Device])
    def get_devices(self) -> Response:
        devices = self.user_service.get_devices()
        return devices

    @router.get("/users/suggested-emails", response_model=list[SuggestedEmail])
    def get_suggested_emails(self) -> Response:
        suggested_emails = self.user_service.get_suggested_emails()
        return suggested_emails

    @router.delete("/users/suggested-emails", response_model=None)
    def delete_suggested_email(self, suggested_email_id: UUID = Form(...)) -> Response:
        Rule.required(suggested_email_id, "suggestedEmailId")
        self.user_service.delete_suggested_email(suggested_email_id)
        return Response(status_code=http.NO_CONTENT)
