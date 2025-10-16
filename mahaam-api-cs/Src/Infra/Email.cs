using Twilio;
using Twilio.Rest.Verify.V2.Service;

namespace Mahaam.Infra;

public interface IEmail
{
	void Init();
	string? SendOtp(string email);
	string? VerifyOtp(string otp, string sid, string email);
}
public class Email(ILog log) : IEmail
{
	private readonly ILog _log = log;
	public void Init()
	{
		try
		{
			TwilioClient.Init(Config.EmailAccountSid, Config.EmailAuthToken);
		}
		catch (Exception e)
		{
			_log.Error(e.ToString());
		}
	}

	public string? SendOtp(string email)
	{
		try
		{
			var verification = VerificationResource.Create(pathServiceSid: Config.EmailVerificationServiceSid, to: email, channel: "email");
			return verification.Sid;
		}
		catch (Exception e)
		{
			_log.Error(e.ToString());
		}
		return null;
	}

	public string VerifyOtp(string otp, string sid, string email)
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
			_log.Error(e.ToString());
		}
		return null;
	}
}