using System.Transactions;
using Mahaam.Feat.Plans;
using Mahaam.Infra;

namespace Mahaam.Feat.Users;

public interface IUserService
{
	Task<CreatedUser> Create(Device device);
	Task<string?> SendMeOtp(string email);
	Task<VerifiedUser> VerifyOtp(string email, string sid, string otp);
	Task<VerifiedUser> RefreshToken();
	Task UpdateName(string name);
	Task Logout(Guid deviceId);
	Task Delete(string sid, string otp);
	Task<List<Device>> GetDevices();
	Task<List<SuggestedEmail>> GetSuggestedEmails();
	Task DeleteSuggestedEmail(Guid suggestedEmailId);
}

class UserService(IUserRepo userRepo, IDeviceRepo deviceRepo, IPlanRepo planRepo, ISuggestedEmailsRepo suggestedEmailsRepo, ILog log, IAuth auth, IEmail email, Settings settings) : IUserService
{
	private readonly IUserRepo _userRepo = userRepo;
	private readonly IDeviceRepo _deviceRepo = deviceRepo;
	private readonly IPlanRepo _planRepo = planRepo;
	private readonly ISuggestedEmailsRepo _suggestedEmailsRepo = suggestedEmailsRepo;
	private readonly ILog _log = log;
	private readonly IAuth _auth = auth;
	private readonly IEmail _email = email;
	private readonly Settings _settings = settings;
	public async Task<CreatedUser> Create(Device device)
	{
		using var scope = new TransactionScope(TransactionScopeAsyncFlowOption.Enabled);
		var userId = await _userRepo.Create();

		// add device
		device.UserId = userId;
		await _deviceRepo.DeleteByFingerprint(device.Fingerprint);
		var deviceId = await _deviceRepo.Create(device);

		string jwt = _auth.CreateToken(userId!.ToString(), deviceId.ToString());
		scope.Complete();

		_log.Info($"User Created with id:{userId}, deviceId:{device.Id}.");
		return new CreatedUser { Id = userId, DeviceId = deviceId, Jwt = jwt };
	}

	public async Task<string?> SendMeOtp(string email)
	{
		string? verifySid;
		if (_settings.Email.TestEmails.Contains(email)) verifySid = _settings.Email.TestSID;
		else verifySid = _email.SendOtp(email);

		if (verifySid != null) _log.Info($"OTP sent to {email}");

		return verifySid;
	}

	public async Task<VerifiedUser> VerifyOtp(string email, string sid, string otp)
	{
		string otpStatus;
		if (_settings.Email.TestEmails.Contains(email) && sid.Equals(_settings.Email.TestSID) && otp.Equals(_settings.Email.TestOTP))
			otpStatus = "approved";
		else
			otpStatus = _email.VerifyOtp(otp, sid, email);

		if (!"approved".Equals(otpStatus))
			throw new ArgumentException($"OTP not verified for {email}, status: {otpStatus}");

		using var scope = new TransactionScope(TransactionScopeAsyncFlowOption.Enabled);
		var user = await _userRepo.GetOne(email);
		var deviceId = Req.DeviceId;
		if (user is null)
		{
			await _userRepo.UpdateEmail(Req.UserId, email);
			_log.Info($"User loggedIn for {email}");
		}
		else
		{
			// move plans of current user to the one with email
			await _planRepo.UpdateUserId(Req.UserId, user.Id);
			var devices = await _deviceRepo.GetMany(user.Id);
			if (devices != null && devices.Count >= 5)
			{
				await _deviceRepo.Delete(devices.Last().Id);
			}

			await _deviceRepo.UpdateUserId(deviceId, user.Id);
			await _userRepo.Delete(Req.UserId);
			_log.Info($"Merging userId:{Req.UserId} to {user.Id}");
		}


		var newUserId = user is null ? Req.UserId! : user.Id;
		string jwt = _auth.CreateToken(newUserId.ToString(), deviceId.ToString());
		scope.Complete();

		_log.Info($"OTP verified for {email}");
		return new VerifiedUser { UserId = newUserId, DeviceId = deviceId, Jwt = jwt, UserFullName = user?.Name, Email = email };
	}

	public async Task<VerifiedUser> RefreshToken()
	{
		var user = await _userRepo.GetOne(Req.UserId);
		string jwt = _auth.CreateToken(Req.UserId.ToString(), Req.DeviceId.ToString());

		return new VerifiedUser { UserId = Req.UserId, DeviceId = Req.DeviceId, Jwt = jwt, UserFullName = user?.Name, Email = user?.Email };
	}

	public async Task UpdateName(string name)
	{
		await _userRepo.UpdateName(Req.UserId, name);
	}

	public async Task Logout(Guid deviceId)
	{
		var device = await _deviceRepo.GetOne(deviceId);
		if (device is null || !device.UserId.Equals(Req.UserId))
			throw new UnauthorizedException("Invalid deviceId");
		await _deviceRepo.Delete(deviceId);
	}

	public async Task DeleteSuggestedEmail(Guid suggestedEmailId)
	{
		var suggestedEmail = await _suggestedEmailsRepo.GetOne(suggestedEmailId);
		if (suggestedEmail is null || !suggestedEmail.UserId.Equals(Req.UserId))
			throw new UnauthorizedException("Invalid suggestedEmailId");
		await _suggestedEmailsRepo.Delete(suggestedEmailId);
	}

	public async Task Delete(string sid, string otp)
	{
		var user = await _userRepo.GetOne(Req.UserId);

		string otpStatus;
		if (user.Email != null && _settings.Email.TestEmails.Contains(user.Email) && sid.Equals(_settings.Email.TestSID) && otp.Equals(_settings.Email.TestOTP))
			otpStatus = "approved";
		else
			otpStatus = _email.VerifyOtp(otp, sid, user.Email ?? string.Empty);

		if (!"approved".Equals(otpStatus))
			throw new ArgumentException($"OTP not approved for {user.Email ?? "unknown"}, status: {otpStatus}");

		using var scope = new TransactionScope(TransactionScopeAsyncFlowOption.Enabled);

		if (user.Email != null)
			await _suggestedEmailsRepo.DeleteManyByEmail(user.Email);
		await _userRepo.Delete(Req.UserId);
		scope.Complete();
	}

	public async Task<List<Device>> GetDevices()
	{
		return await _deviceRepo.GetMany(Req.UserId);
	}

	public async Task<List<SuggestedEmail>> GetSuggestedEmails()
	{
		return await _suggestedEmailsRepo.GetMany(Req.UserId);
	}
}