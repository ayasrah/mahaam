using Mahaam.Infra;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.RateLimiting;

namespace Mahaam.Feat.Users;

public interface IUserController
{
	Task<IActionResult> Create(
		 string platform,
		 bool isPhysicalDevice,
		 string deviceFingerprint,
		 string deviceInfo);
	Task<IActionResult> SendMeOtp(string email);
	Task<IActionResult> VerifyOtp(string email, string sid, string otp);
	Task<IActionResult> RefreshToken();
	Task<IActionResult> UpdateName(string name);
	Task<IActionResult> Logout(Guid deviceId);
	Task<IActionResult> Delete(string sid, string otp);
	Task<IActionResult> GetDevices();
	Task<IActionResult> GetSuggestedEmails();
	Task<IActionResult> DeleteSuggestedEmail(Guid suggestedEmailId);
}

[ApiController]
[Route("users")]
public class UserController(IUserService userService) : ControllerBase, IUserController
{
	[HttpPost]
	[EnableRateLimiting("PerUserRateLimit")]
	[Route("send-me-otp")]
	public async Task<IActionResult> SendMeOtp([FromForm] string email)
	{
		Rule.ValidateEmail(email);
		var verificationSid = await userService.SendMeOtp(email);
		return Ok(verificationSid);

	}

	[HttpPost]
	[EnableRateLimiting("PerUserRateLimit")]
	[Route("create")]
	public async Task<IActionResult> Create(
		[FromForm] string platform,
		[FromForm] bool isPhysicalDevice,
		[FromForm] string deviceFingerprint,
		[FromForm] string deviceInfo)
	{

		Rule.Required(isPhysicalDevice, "isPhysicalDevice");
		Rule.Required(platform, "platform");
		Rule.Required(deviceFingerprint, "deviceFingerprint");
		Rule.Required(deviceInfo, "deviceInfo");
		Rule.FailIf(!isPhysicalDevice, "Device should be real not simulator");
		var device = new Device { Platform = platform, Fingerprint = deviceFingerprint, Info = deviceInfo };

		var createdUser = await userService.Create(device);
		return Ok(createdUser);
	}

	[HttpPost]
	[Route("verify-otp")]
	public async Task<IActionResult> VerifyOtp(
		[FromForm] string email,
		[FromForm] string sid,
		[FromForm] string otp)
	{

		Rule.Required(email, "email");
		Rule.Required(sid, "sid");
		Rule.Required(otp, "otp");

		var verifiedUser = await userService.VerifyOtp(email, sid, otp);

		return Ok(verifiedUser);
	}

	[HttpPost]
	[Route("refresh-token")]
	public async Task<IActionResult> RefreshToken()
	{
		var verifiedUser = await userService.RefreshToken();
		return Ok(verifiedUser);
	}

	[HttpPatch]
	[Route("name")]
	public async Task<IActionResult> UpdateName([FromForm] string name)
	{
		Rule.Required(name, "name");
		await userService.UpdateName(name);
		return Ok();
	}

	[HttpPost]
	[Route("logout")]
	public async Task<IActionResult> Logout([FromForm] Guid deviceId)
	{
		Rule.Required(deviceId, "deviceId");
		await userService.Logout(deviceId);
		return Ok();
	}

	[HttpDelete]
	[Route("")]
	public async Task<IActionResult> Delete([FromForm] string sid, [FromForm] string otp)
	{
		Rule.Required(sid, "sid");
		Rule.Required(otp, "otp");
		await userService.Delete(sid, otp);
		return NoContent();
	}

	[HttpGet]
	[Route("devices")]
	public async Task<IActionResult> GetDevices()
	{
		var devices = await userService.GetDevices();
		return Ok(devices);
	}

	[HttpGet]
	[Route("suggested-emails")]
	public async Task<IActionResult> GetSuggestedEmails()
	{
		var suggestedEmails = await userService.GetSuggestedEmails();
		return Ok(suggestedEmails);
	}

	[HttpDelete]
	[Route("suggested-emails")]
	public async Task<IActionResult> DeleteSuggestedEmail([FromForm] Guid suggestedEmailId)
	{
		Rule.Required(suggestedEmailId, "suggestedEmailId");
		await userService.DeleteSuggestedEmail(suggestedEmailId);
		return NoContent();
	}
}