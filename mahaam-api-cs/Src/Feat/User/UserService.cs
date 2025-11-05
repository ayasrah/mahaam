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

class UserService(IUserRepo userRepo, IDeviceRepo deviceRepo, IPlanRepo planRepo, ISuggestedEmailsRepo suggestedEmailsRepo, ILog log, IAuth auth, IEmail emailService, Settings settings) : IUserService
{
	public async Task<CreatedUser> Create(Device device)
	{
		using var scope = new TransactionScope(TransactionScopeAsyncFlowOption.Enabled);
		var userId = await userRepo.Create();

		// add device
		device.UserId = userId;
		await deviceRepo.DeleteByFingerprint(device.Fingerprint);
		var deviceId = await deviceRepo.Create(device);

		string jwt = auth.CreateToken(userId!.ToString(), deviceId.ToString());
		scope.Complete();

		log.Info($"User Created with id:{userId}, deviceId:{device.Id}.");
		return new CreatedUser { Id = userId, DeviceId = deviceId, Jwt = jwt };
	}

	public async Task<string?> SendMeOtp(string email)
	{
		string? verifySid;
		if (settings.Email.TestEmails.Contains(email)) verifySid = settings.Email.TestSID;
		else verifySid = emailService.SendOtp(email);

		if (verifySid != null) log.Info($"OTP sent to {email}");

		return verifySid;
	}

	public async Task<VerifiedUser> VerifyOtp(string email, string sid, string otp)
	{
		string otpStatus;
		if (settings.Email.TestEmails.Contains(email) && sid.Equals(settings.Email.TestSID) && otp.Equals(settings.Email.TestOTP))
			otpStatus = "approved";
		else
			otpStatus = emailService.VerifyOtp(otp, sid, email);

		if (!"approved".Equals(otpStatus))
			throw new ArgumentException($"OTP not verified for {email}, status: {otpStatus}");

		using var scope = new TransactionScope(TransactionScopeAsyncFlowOption.Enabled);
		var user = await userRepo.GetOne(email);
		var deviceId = Req.DeviceId;
		if (user is null)
		{
			await userRepo.UpdateEmail(Req.UserId, email);
			log.Info($"User loggedIn for {email}");
		}
		else
		{
			// move plans of current user to the one with email
			await planRepo.UpdateUserId(Req.UserId, user.Id);
			var devices = await deviceRepo.GetMany(user.Id);
			if (devices != null && devices.Count >= 5)
			{
				await deviceRepo.Delete(devices.Last().Id);
			}

			await deviceRepo.UpdateUserId(deviceId, user.Id);
			await userRepo.Delete(Req.UserId);
			log.Info($"Merging userId:{Req.UserId} to {user.Id}");
		}


		var newUserId = user is null ? Req.UserId! : user.Id;
		string jwt = auth.CreateToken(newUserId.ToString(), deviceId.ToString());
		scope.Complete();

		log.Info($"OTP verified for {email}");
		return new VerifiedUser { UserId = newUserId, DeviceId = deviceId, Jwt = jwt, UserFullName = user?.Name, Email = email };
	}

	public async Task<VerifiedUser> RefreshToken()
	{
		var user = await userRepo.GetOne(Req.UserId);
		string jwt = auth.CreateToken(Req.UserId.ToString(), Req.DeviceId.ToString());

		return new VerifiedUser { UserId = Req.UserId, DeviceId = Req.DeviceId, Jwt = jwt, UserFullName = user?.Name, Email = user?.Email };
	}

	public async Task UpdateName(string name)
	{
		await userRepo.UpdateName(Req.UserId, name);
	}

	public async Task Logout(Guid deviceId)
	{
		var device = await deviceRepo.GetOne(deviceId);
		if (device is null || !device.UserId.Equals(Req.UserId))
			throw new UnauthorizedException("Invalid deviceId");
		await deviceRepo.Delete(deviceId);
	}

	public async Task DeleteSuggestedEmail(Guid suggestedEmailId)
	{
		var suggestedEmail = await suggestedEmailsRepo.GetOne(suggestedEmailId);
		if (suggestedEmail is null || !suggestedEmail.UserId.Equals(Req.UserId))
			throw new UnauthorizedException("Invalid suggestedEmailId");
		await suggestedEmailsRepo.Delete(suggestedEmailId);
	}

	public async Task Delete(string sid, string otp)
	{
		var user = await userRepo.GetOne(Req.UserId);

		string otpStatus;
		if (user.Email != null && settings.Email.TestEmails.Contains(user.Email) && sid.Equals(settings.Email.TestSID) && otp.Equals(settings.Email.TestOTP))
			otpStatus = "approved";
		else
			otpStatus = emailService.VerifyOtp(otp, sid, user.Email ?? string.Empty);

		if (!"approved".Equals(otpStatus))
			throw new ArgumentException($"OTP not approved for {user.Email ?? "unknown"}, status: {otpStatus}");

		using var scope = new TransactionScope(TransactionScopeAsyncFlowOption.Enabled);

		if (user.Email != null)
			await suggestedEmailsRepo.DeleteManyByEmail(user.Email);
		await userRepo.Delete(Req.UserId);
		scope.Complete();
	}

	public async Task<List<Device>> GetDevices()
	{
		return await deviceRepo.GetMany(Req.UserId);
	}

	public async Task<List<SuggestedEmail>> GetSuggestedEmails()
	{
		return await suggestedEmailsRepo.GetMany(Req.UserId);
	}
}