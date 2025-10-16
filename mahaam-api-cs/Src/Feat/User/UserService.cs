using System.Transactions;
using Mahaam.Feat.Plans;
using Mahaam.Infra;

namespace Mahaam.Feat.Users;

public interface IUserService
{
	CreatedUser Create(Device device);
	string? SendMeOtp(string email);
	VerifiedUser VerifyOtp(string email, string sid, string otp);
	VerifiedUser RefreshToken();
	void UpdateName(string name);
	void Logout(Guid deviceId);
	void Delete(string sid, string otp);
	List<Device> GetDevices();
	List<SuggestedEmail> GetSuggestedEmails();
	void DeleteSuggestedEmail(Guid suggestedEmailId);
}

class UserService(IUserRepo userRepo, IDeviceRepo deviceRepo, IPlanRepo planRepo, ISuggestedEmailsRepo suggestedEmailsRepo, ILog log, IAuth auth, IEmail email) : IUserService
{
	private readonly IUserRepo _userRepo = userRepo;
	private readonly IDeviceRepo _deviceRepo = deviceRepo;
	private readonly IPlanRepo _planRepo = planRepo;
	private readonly ISuggestedEmailsRepo _suggestedEmailsRepo = suggestedEmailsRepo;
	private readonly ILog _log = log;
	private readonly IAuth _auth = auth;
	private readonly IEmail _email = email;
	public CreatedUser Create(Device device)
	{
		using var scope = new TransactionScope();
		var userId = _userRepo.Create();

		// add device
		device.UserId = userId;
		_deviceRepo.DeleteByFingerprint(device.Fingerprint);
		var deviceId = _deviceRepo.Create(device);

		string jwt = _auth.CreateToken(userId!.ToString(), deviceId.ToString());
		scope.Complete();

		_log.Info($"User Created with id:{userId}, deviceId:{device.Id}.");
		return new CreatedUser { Id = userId, DeviceId = deviceId, Jwt = jwt };
	}

	public string? SendMeOtp(string email)
	{
		string? verifySid;
		if (Config.TestEmails.Contains(email)) verifySid = Config.TestSID;
		else verifySid = _email.SendOtp(email);

		if (verifySid != null) _log.Info($"OTP sent to {email}");

		return verifySid;
	}

	public VerifiedUser VerifyOtp(string email, string sid, string otp)
	{
		string otpStatus;
		if (Config.TestEmails.Contains(email) && sid.Equals(Config.TestSID) && otp.Equals(Config.TestOTP))
			otpStatus = "approved";
		else
			otpStatus = _email.VerifyOtp(otp, sid, email);

		if (!"approved".Equals(otpStatus))
			throw new ArgumentException($"OTP not verified for {email}, status: {otpStatus}");

		using var scope = new TransactionScope();
		var user = _userRepo.GetOne(email);
		var deviceId = Req.DeviceId;
		if (user is null)
		{
			_userRepo.UpdateEmail(Req.UserId, email);
			_log.Info($"User loggedIn for {email}");
		}
		else
		{
			// move plans of current user to the one with email
			_planRepo.UpdateUserId(Req.UserId, user.Id);
			var devices = _deviceRepo.GetMany(user.Id);
			if (devices != null && devices.Count >= 5)
			{
				_deviceRepo.Delete(devices.Last().Id);
			}

			_deviceRepo.UpdateUserId(deviceId, user.Id);
			_userRepo.Delete(Req.UserId);
			_log.Info($"Merging userId:{Req.UserId} to {user.Id}");
		}


		var newUserId = user is null ? Req.UserId! : user.Id;
		string jwt = _auth.CreateToken(newUserId.ToString(), deviceId.ToString());
		scope.Complete();

		_log.Info($"OTP verified for {email}");
		return new VerifiedUser { UserId = newUserId, DeviceId = deviceId, Jwt = jwt, UserFullName = user?.Name, Email = email };
	}

	public VerifiedUser RefreshToken()
	{
		var user = _userRepo.GetOne(Req.UserId);
		string jwt = _auth.CreateToken(Req.UserId.ToString(), Req.DeviceId.ToString());

		return new VerifiedUser { UserId = Req.UserId, DeviceId = Req.DeviceId, Jwt = jwt, UserFullName = user?.Name, Email = user?.Email };
	}

	public void UpdateName(string name)
	{
		_userRepo.UpdateName(Req.UserId, name);
	}

	public void Logout(Guid deviceId)
	{
		var device = _deviceRepo.GetOne(deviceId);
		if (device is null || !device.UserId.Equals(Req.UserId))
			throw new UnauthorizedException("Invalid deviceId");
		_deviceRepo.Delete(deviceId);
	}

	public void DeleteSuggestedEmail(Guid suggestedEmailId)
	{
		var suggestedEmail = _suggestedEmailsRepo.GetOne(suggestedEmailId);
		if (suggestedEmail is null || !suggestedEmail.UserId.Equals(Req.UserId))
			throw new UnauthorizedException("Invalid suggestedEmailId");
		_suggestedEmailsRepo.Delete(suggestedEmailId);
	}

	public void Delete(string sid, string otp)
	{
		var user = _userRepo.GetOne(Req.UserId);

		string otpStatus;
		if (user.Email != null && Config.TestEmails.Contains(user.Email) && sid.Equals(Config.TestSID) && otp.Equals(Config.TestOTP))
			otpStatus = "approved";
		else
			otpStatus = _email.VerifyOtp(otp, sid, user.Email ?? string.Empty);

		if (!"approved".Equals(otpStatus))
			throw new ArgumentException($"OTP not approved for {user.Email ?? "unknown"}, status: {otpStatus}");

		using var scope = new TransactionScope();

		if (user.Email != null)
			_suggestedEmailsRepo.DeleteManyByEmail(user.Email);
		_userRepo.Delete(Req.UserId);
		scope.Complete();
	}

	public List<Device> GetDevices()
	{
		return _deviceRepo.GetMany(Req.UserId);
	}

	public List<SuggestedEmail> GetSuggestedEmails()
	{
		return _suggestedEmailsRepo.GetMany(Req.UserId);
	}
}