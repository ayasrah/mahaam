
using Mahaam.Infra;
using Twilio;
using Twilio.Rest.Verify.V2.Service;

namespace Mahaam.Feat.Users;

public interface IEmail
{
	void Init();
	string? SendOtp(string email);
	string? VerifyOtp(string otp, string sid, string email);
}
public class Email(ILog log, Settings settings) : IEmail
{

	public void Init()
	{
		try
		{
			TwilioClient.Init(settings.Email.AccountSid, settings.Email.AuthToken);
		}
		catch (Exception e)
		{
			log.Error(e.ToString());
		}
	}

	public string? SendOtp(string email)
	{
		try
		{
			var verification = VerificationResource.Create(pathServiceSid: settings.Email.VerificationServiceSid, to: email, channel: "email");
			return verification.Sid;
		}
		catch (Exception e)
		{
			log.Error(e.ToString());
		}
		return null;
	}

	public string VerifyOtp(string otp, string sid, string email)
	{
		try
		{
			var check = VerificationCheckResource.Create(to: email, code: otp, verificationSid: sid, pathServiceSid: settings.Email.VerificationServiceSid
			// ,verificationSid:sid was not there
			);
			return check.Status;
		}
		catch (Exception e)
		{
			log.Error(e.ToString());
		}
		return null;
	}
}