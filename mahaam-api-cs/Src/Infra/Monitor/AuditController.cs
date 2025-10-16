using System.Net;
using Microsoft.AspNetCore.Mvc;

namespace Mahaam.Infra.Monitoring;

public interface IAuditController
{
	IActionResult Error(string error);
	IActionResult Info(string info);
}

[Route("audit")]
[ApiController]
public class AuditController(ILog log) : ControllerBase, IAuditController
{
	private readonly ILog _log = log;
	[HttpPost]
	[Route("error")]
	public IActionResult Error([FromForm] string error)
	{

		_log.Error(error);
		return StatusCode((int)HttpStatusCode.Created);
	}

	[HttpPost]
	[Route("info")]
	public IActionResult Info([FromForm] string info)
	{
		_log.Info("mahaam-mb:" + info);
		return StatusCode((int)HttpStatusCode.Created);
	}
}
