using System.Net;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Options;

namespace Mahaam.Infra.Monitoring;

public interface IAuditController
{
	IActionResult Error(string error);
	IActionResult Info(string info);
}

[Route("audit")]
[ApiController]
public class AuditController(ILog log, Settings settings) : ControllerBase, IAuditController
{

	[HttpPost]
	[Route("error")]
	public IActionResult Error([FromForm] string error)
	{

		log.Error(error);
		return StatusCode((int)HttpStatusCode.Created);
	}

	[HttpPost]
	[Route("info")]
	public IActionResult Info([FromForm] string info)
	{
		log.Info($"{settings.Api.Name}-v{settings.Api.Version}: {info}");

		return StatusCode((int)HttpStatusCode.Created);
	}
}
