using Mahaam.Infra;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.RateLimiting;

namespace Mahaam.Feat.Users;

public interface IUserController
{
	IActionResult Create(
		 string platform,
		 bool isPhysicalDevice,
		 string deviceFingerprint,
		 string deviceInfo);
	IActionResult SendMeOtp(string email);
	IActionResult VerifyOtp(string email, string sid, string otp);
	IActionResult RefreshToken();
	IActionResult UpdateName(string name);
	IActionResult Logout(Guid deviceId);
	IActionResult Delete(string sid, string otp);
	IActionResult GetDevices();
	IActionResult GetSuggestedEmails();
	IActionResult DeleteSuggestedEmail(Guid suggestedEmailId);
}

[ApiController]
[Route("users")]
public class UserController : ControllerBase, IUserController
{

	[HttpPost]
	[EnableRateLimiting("PerUserRateLimit")]
	[Route("send-me-otp")]
	public IActionResult SendMeOtp([FromForm] string email)
	{
		Rule.ValidateEmail(email);
		var verificationSid = App.UserService.SendMeOtp(email);
		return Ok(verificationSid);

	}

	[HttpPost]
	[EnableRateLimiting("PerUserRateLimit")]
	[Route("create")]
	public IActionResult Create(
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

		var createdUser = App.UserService.Create(device);
		return Ok(createdUser);
	}

	[HttpPost]
	[Route("verify-otp")]
	public IActionResult VerifyOtp(
		[FromForm] string email,
		[FromForm] string sid,
		[FromForm] string otp)
	{

		Rule.Required(email, "email");
		Rule.Required(sid, "sid");
		Rule.Required(otp, "otp");

		var verifiedUser = App.UserService.VerifyOtp(email, sid, otp);

		return Ok(verifiedUser);
	}

	[HttpPost]
	[Route("refresh-token")]
	public IActionResult RefreshToken()
	{
		var verifiedUser = App.UserService.RefreshToken();
		return Ok(verifiedUser);
	}

	[HttpPatch]
	[Route("name")]
	public IActionResult UpdateName([FromForm] string name)
	{
		Rule.Required(name, "name");
		App.UserService.UpdateName(name);
		return Ok();
	}

	[HttpPost]
	[Route("logout")]
	public IActionResult Logout([FromForm] Guid deviceId)
	{
		Rule.Required(deviceId, "deviceId");
		App.UserService.Logout(deviceId);
		return Ok();
	}

	[HttpDelete]
	[Route("")]
	public IActionResult Delete([FromForm] string sid, [FromForm] string otp)
	{
		Rule.Required(sid, "sid");
		Rule.Required(otp, "otp");
		App.UserService.Delete(sid, otp);
		return NoContent();
	}

	[HttpGet]
	[Route("devices")]
	public IActionResult GetDevices()
	{
		var devices = App.UserService.GetDevices();
		return Ok(devices);
	}

	[HttpGet]
	[Route("suggested-emails")]
	public IActionResult GetSuggestedEmails()
	{
		var suggestedEmails = App.UserService.GetSuggestedEmails();
		return Ok(suggestedEmails);
	}

	[HttpDelete]
	[Route("suggested-emails")]
	public IActionResult DeleteSuggestedEmail([FromForm] Guid suggestedEmailId)
	{
		Rule.Required(suggestedEmailId, "suggestedEmailId");
		App.UserService.DeleteSuggestedEmail(suggestedEmailId);
		return NoContent();
	}
}