from uuid import UUID
from typing import List, Protocol
from feat.plan.plan_repo import PlanRepo
from feat.user.user_repo import UserRepo
from feat.user.device_repo import DeviceRepo
from feat.user.suggested_emails_repo import SuggestedEmailsRepo
from feat.user.user_model import Device, SuggestedEmail, VerifiedUser, CreatedUser
from infra.exceptions import InputException, UnauthorizedException
from infra.req import req
from infra.security import Auth
from infra.email import send_otp, verify_otp
from infra.validation import ProtocolEnforcer
from infra import db, configs
from infra.log import Log

class UserService(Protocol):
    def create(self, device: Device) -> CreatedUser: ...
    def send_me_otp(self, email: str) -> str | None: ...
    def verify_otp(self, email: str, sid: str, otp: str) -> VerifiedUser: ...
    def refresh_token(self) -> VerifiedUser: ...
    def update_name(self, name: str) -> None: ...
    def logout(self, device_id: UUID) -> None: ...
    def delete_suggested_email(self, suggested_email_id: UUID) -> None: ...
    def delete(self, sid: str, otp: str) -> None: ...
    def get_devices(self) -> List[Device]: ...
    def get_suggested_emails(self) -> List[SuggestedEmail]: ...

class DefaultUserService(metaclass=ProtocolEnforcer, protocol=UserService):
    def __init__(self, user_repo: UserRepo, plan_repo: PlanRepo, device_repo: DeviceRepo, suggested_emails_repo: SuggestedEmailsRepo) -> None:
        self.user_repo = user_repo
        self.plan_repo = plan_repo
        self.device_repo = device_repo
        self.suggested_emails_repo = suggested_emails_repo



    def create(self, device: Device) -> CreatedUser:
        with db.DB.transaction_scope() as conn:
            user_id = self.user_repo.create(conn)
            device.user_id = user_id
            self.device_repo.delete_by_fingerprint(device.fingerprint, conn)
            device_id = self.device_repo.create(device, conn)
            jwt = Auth.create_token(str(user_id), str(device_id))
        Log.info(f"User Created with id:{user_id}, deviceId:{device_id}.")
        return CreatedUser(id=user_id, device_id=device_id, jwt=jwt)

    def send_me_otp(self, email: str) -> str | None:
        if email in configs.data.testEmails:
            verify_sid = configs.data.testSID
        else:
            verify_sid, err = send_otp(email)
            if err is not None:
                Log.info(f"Error sending OTP to {email}")
                return None
        if verify_sid is not None:
            Log.info(f"OTP sent to {email}")
        return verify_sid

    def verify_otp(self, email: str, sid: str, otp: str) -> VerifiedUser:
        if email in configs.data.testEmails and sid == configs.data.testSID and otp == configs.data.testOTP:
            otp_status = "approved"
        else:
            otp_status, err = verify_otp(otp, sid, email)
            if err is not None or otp_status != "approved":
                raise InputException(f"OTP not verified for {email}, status: {otp_status}")

        with db.DB.transaction_scope() as conn:
            user = self.user_repo.select_one_by_email(email, conn)
            if user is None:
                self.user_repo.update_email(req.user_id, email, conn)
                Log.info(f"User loggedIn for {email}")
            else:
				# move plans of current user to the one with email
                self.plan_repo.update_user_id(req.user_id, user.id, conn)
                devices = self.device_repo.select_many(user.id, conn)
                if devices and len(devices) >= 5:
                    self.device_repo.delete(devices[-1].id, conn)
                self.device_repo.update_user_id(req.device_id, user.id, conn)
                self.user_repo.delete(req.user_id, conn)
                Log.info(f"Merging userId:{req.user_id} to {user.id}")
				
            new_user_id = req.user_id if user is None else user.id
            jwt = Auth.create_token(str(new_user_id), str(req.device_id))
        Log.info(f"OTP verified for {email}")
        return VerifiedUser(user_id=new_user_id, device_id=req.device_id, jwt=jwt, user_full_name=getattr(user, 'name', None), email=email)

    def refresh_token(self) -> VerifiedUser:
        user = self.user_repo.select_one(req.user_id)
        jwt = Auth.create_token(str(req.user_id), str(req.device_id))
        return VerifiedUser(user_id=req.user_id, device_id=req.device_id, jwt=jwt, user_full_name=getattr(user, 'name', None), email=getattr(user, 'email', None))

    def update_name(self, name: str) -> None:
        self.user_repo.update_name(req.user_id, name)

    def logout(self, device_id: UUID) -> None:
        device = self.device_repo.select_one(device_id)
        if device is None or device.user_id != req.user_id:
            raise UnauthorizedException("Invalid deviceId")
        self.device_repo.delete(device_id)

    def delete_suggested_email(self, suggested_email_id: UUID) -> None:
        suggested_email = self.suggested_emails_repo.select_one(suggested_email_id)
        if suggested_email is None or suggested_email.user_id != req.user_id:
            raise UnauthorizedException("Invalid suggestedEmailId")
        self.suggested_emails_repo.delete(suggested_email_id)

    def delete(self, sid: str, otp: str) -> None:
        user = self.user_repo.select_one(req.user_id)
        email = getattr(user, 'email', None)
        if email in configs.data.testEmails and sid == configs.data.testSID and otp == configs.data.testOTP:
            otp_status = "approved"
        else:
            otp_status, err = verify_otp(otp, sid, email)
            if err is not None or otp_status != "approved":
                raise InputException(f"OTP not approved for {email}, status: {otp_status}")
        with db.DB.transaction_scope() as conn:
            if email is not None:
                self.suggested_emails_repo.delete_many_by_email(email, conn)
            self.user_repo.delete(req.user_id, conn)

    def get_devices(self) -> List[Device]:
        return self.device_repo.select_many(req.user_id)

    def get_suggested_emails(self) -> List[SuggestedEmail]:
        return self.suggested_emails_repo.select_many(req.user_id)
