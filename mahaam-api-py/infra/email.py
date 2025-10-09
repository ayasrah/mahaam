import uuid
from twilio.rest import Client
from twilio.base.exceptions import TwilioRestException
from infra import db, configs
from infra.log import Log

client: Client = None


def init():
    global client
    client = Client(configs.data.emailAccountSid, configs.data.emailAuthToken)


def send_otp(email: str) -> tuple[str, Exception]:
    try:
        verification = client.verify.v2.services(configs.data.emailVerificationServiceSid) \
            .verifications \
            .create(to=email, channel='email')
        return verification.sid, None
    except TwilioRestException as e:
        Log.error(str(uuid.UUID(int=0)),
                  f"Error sending OTP to {email}: {str(e)}")
        return "", e


def verify_otp(otp: str, sid: str, email: str) -> tuple[str, Exception]:
    try:
        verification_check = client.verify.v2.services(configs.data.emailVerificationServiceSid) \
            .verification_checks \
            .create(to=email, code=otp, verification_sid=sid)
        return verification_check.status, None
    except TwilioRestException as e:
        Log.info(str(uuid.UUID(int=0)),
                 f"Error verifying OTP for {email}: {str(e)}")
        return "", e
