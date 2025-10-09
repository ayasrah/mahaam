using Twilio;
using Twilio.Rest.Verify.V2.Service;

namespace Mahaam.Infra;

public class Email
{
	public static void Init()
	{
		try
		{
			TwilioClient.Init(Config.EmailAccountSid, Config.EmailAuthToken);
		}
		catch (Exception e)
		{
			Log.Error(e.ToString());
		}
	}

	public static string? SendOtp(string email)
	{
		try
		{
			var verification = VerificationResource.Create(pathServiceSid: Config.EmailVerificationServiceSid, to: email, channel: "email");
			return verification.Sid;
		}
		catch (Exception e)
		{
			Log.Error(e.ToString());
		}
		return null;
	}

	public static string VerifyOtp(string otp, string sid, string email)
	{
		try
		{
			var check = VerificationCheckResource.Create(to: email, code: otp, verificationSid: sid, pathServiceSid: Config.EmailVerificationServiceSid
			// ,verificationSid:sid was not there
			);
			return check.Status;
		}
		catch (Exception e)
		{
			Log.Error(e.ToString());
		}
		return null;
	}
}