package mahaam.infra;

import com.twilio.Twilio;
import com.twilio.rest.verify.v2.service.Verification;
import com.twilio.rest.verify.v2.service.VerificationCheck;

public class Email {
	private static String verificationServiceSid;

	public static void init(String accountSid, String verificationServiceSid, String authToken) {
		try {
			Email.verificationServiceSid = verificationServiceSid;
			Twilio.init(accountSid, authToken);
		} catch (Exception e) {
			Log.error(e.toString());
		}
	}

	public static String sendOtp(String email) {
		try {
			Verification verification = Verification.creator(verificationServiceSid, email, "email").create();
			return verification.getSid();
		} catch (Exception e) {
			Log.error(e.toString());
		}
		return null;
	}

	public static String verifyOtp(String otp, String sid, String email) {
		try {
			VerificationCheck check = VerificationCheck.creator(verificationServiceSid).setTo(email).setCode(otp)
					.create();
			return check.getStatus().toString();
		} catch (Exception e) {
			Log.error(e.toString());
		}
		return null;
	}
}
