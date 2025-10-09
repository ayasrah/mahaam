using System.Transactions;
using Mahaam.Infra;

namespace Mahaam.Feat.Users;

public interface IUserService
{
	CreatedUser Create(Device device);
	string? SendMeOtp(string email);
	VerifiedUser VerifyOtp(string email, string sid, string otp);
	public VerifiedUser RefreshToken();
	void UpdateName(string name);
	void Logout(Guid deviceId);
	void Delete(string sid, string otp);
	List<Device> GetDevices();
	List<SuggestedEmail> GetSuggestedEmails();
	public void DeleteSuggestedEmail(Guid suggestedEmailId);
}

class UserService : IUserService
{

	public CreatedUser Create(Device device)
	{
		using var scope = new TransactionScope();
		var userId = App.UserRepo.Create();

		// add device
		device.UserId = userId;
		App.DeviceRepo.DeleteByFingerprint(device.Fingerprint);
		var deviceId = App.DeviceRepo.Create(device);

		string jwt = Auth.CreateToken(userId!.ToString(), deviceId.ToString());
		scope.Complete();

		Log.Info($"User Created with id:{userId}, deviceId:{device.Id}.");
		return new CreatedUser { Id = userId, DeviceId = deviceId, Jwt = jwt };
	}

	public string? SendMeOtp(string email)
	{
		string? verifySid;
		if (Config.TestEmails.Contains(email)) verifySid = Config.TestSID;
		else verifySid = Email.SendOtp(email);

		if (verifySid != null) Log.Info($"OTP sent to {email}");

		return verifySid;
	}

	public VerifiedUser VerifyOtp(string email, string sid, string otp)
	{
		string otpStatus;
		if (Config.TestEmails.Contains(email) && sid.Equals(Config.TestSID) && otp.Equals(Config.TestOTP))
			otpStatus = "approved";
		else
			otpStatus = Email.VerifyOtp(otp, sid, email);

		if (!"approved".Equals(otpStatus))
			throw new ArgumentException($"OTP not verified for {email}, status: {otpStatus}");

		using var scope = new TransactionScope();
		var user = App.UserRepo.GetOne(email);
		var deviceId = Req.DeviceId;
		if (user is null)
		{
			App.UserRepo.UpdateEmail(Req.UserId, email);
			Log.Info($"User loggedIn for {email}");
		}
		else
		{
			// move plans of current user to the one with email
			App.PlanRepo.UpdateUserId(Req.UserId, user.Id);
			var devices = App.DeviceRepo.GetMany(user.Id);
			if (devices != null && devices.Count >= 5)
			{
				App.DeviceRepo.Delete(devices.Last().Id);
			}

			App.DeviceRepo.UpdateUserId(deviceId, user.Id);
			App.UserRepo.Delete(Req.UserId);
			Log.Info($"Merging userId:{Req.UserId} to {user.Id}");
		}


		var newUserId = user is null ? Req.UserId! : user.Id;
		string jwt = Auth.CreateToken(newUserId.ToString(), deviceId.ToString());
		scope.Complete();

		Log.Info($"OTP verified for {email}");
		return new VerifiedUser { UserId = newUserId, DeviceId = deviceId, Jwt = jwt, UserFullName = user?.Name, Email = email };
	}

	public VerifiedUser RefreshToken()
	{
		var user = App.UserRepo.GetOne(Req.UserId);
		string jwt = Auth.CreateToken(Req.UserId.ToString(), Req.DeviceId.ToString());

		return new VerifiedUser { UserId = Req.UserId, DeviceId = Req.DeviceId, Jwt = jwt, UserFullName = user?.Name, Email = user?.Email };
	}

	public void UpdateName(string name)
	{
		App.UserRepo.UpdateName(Req.UserId, name);
	}

	public void Logout(Guid deviceId)
	{
		var device = App.DeviceRepo.GetOne(deviceId);
		if (device is null || !device.UserId.Equals(Req.UserId))
			throw new UnauthorizedException("Invalid deviceId");
		App.DeviceRepo.Delete(deviceId);
	}

	public void DeleteSuggestedEmail(Guid suggestedEmailId)
	{
		var suggestedEmail = App.SuggestedEmailsRepo.GetOne(suggestedEmailId);
		if (suggestedEmail is null || !suggestedEmail.UserId.Equals(Req.UserId))
			throw new UnauthorizedException("Invalid suggestedEmailId");
		App.SuggestedEmailsRepo.Delete(suggestedEmailId);
	}

	public void Delete(string sid, string otp)
	{
		var user = App.UserRepo.GetOne(Req.UserId);

		string otpStatus;
		if (user.Email != null && Config.TestEmails.Contains(user.Email) && sid.Equals(Config.TestSID) && otp.Equals(Config.TestOTP))
			otpStatus = "approved";
		else
			otpStatus = Email.VerifyOtp(otp, sid, user.Email ?? string.Empty);

		if (!"approved".Equals(otpStatus))
			throw new ArgumentException($"OTP not approved for {user.Email ?? "unknown"}, status: {otpStatus}");

		using var scope = new TransactionScope();

		if (user.Email != null)
			App.SuggestedEmailsRepo.DeleteManyByEmail(user.Email);
		App.UserRepo.Delete(Req.UserId);
		scope.Complete();
	}

	public List<Device> GetDevices()
	{
		return App.DeviceRepo.GetMany(Req.UserId);
	}

	public List<SuggestedEmail> GetSuggestedEmails()
	{
		return App.SuggestedEmailsRepo.GetMany(Req.UserId);
	}
}