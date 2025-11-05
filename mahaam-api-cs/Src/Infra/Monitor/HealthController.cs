using Mahaam.Infra;
using Microsoft.AspNetCore.Mvc;

namespace Mahaam.Infra.Monitoring;

public interface IHealthController
{
	IActionResult GetStatus();
}

[ApiController]
[Route("health")]
public class HealthController(Settings settings, ICache cache) : ControllerBase, IHealthController
{
	[HttpGet]
	public IActionResult GetStatus()
	{
		var result = new Health()
		{
			Id = cache.HealthId(),
			ApiName = settings.Api.Name,
			ApiVersion = settings.Api.Version,
			NodeIP = cache.NodeIP(),
			NodeName = cache.NodeName(),
			EnvName = settings.Api.EnvName
		};
		return Ok(result);
	}
}
