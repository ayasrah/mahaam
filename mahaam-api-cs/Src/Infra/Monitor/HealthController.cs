using Mahaam.Infra;
using Microsoft.AspNetCore.Mvc;

namespace Mahaam.Infra.Monitoring;

public interface IHealthController
{
	IActionResult GetStatus();
}

[ApiController]
[Route("health")]
public class HealthController : ControllerBase, IHealthController
{
	[HttpGet]
	public IActionResult GetStatus()
	{
		var result = new Health()
		{
			Id = Cache.HealthId,
			ApiName = Cache.ApiName,
			ApiVersion = Cache.ApiVersion,
			NodeIP = Cache.NodeIP,
			NodeName = Cache.NodeName,
			EnvName = Cache.EnvName
		};
		return StatusCode(Http.Ok, result);
	}
}
