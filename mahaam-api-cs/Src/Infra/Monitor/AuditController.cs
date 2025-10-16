using Microsoft.AspNetCore.Mvc;

namespace Mahaam.Infra.Monitoring;

public interface IAuditController
{
	IActionResult Error(string error);
	IActionResult Info(string info);
}

[Route("audit")]
[ApiController]
public class AuditController : ControllerBase, IAuditController
{

	[HttpPost]
	[Route("error")]
	public IActionResult Error([FromForm] string error)
	{

		Log.Error(error);
		return Created();
	}

	[HttpPost]
	[Route("info")]
	public IActionResult Info([FromForm] string info)
	{
		Log.Info("mahaam-mb:" + info);
		return Created();
	}
}
