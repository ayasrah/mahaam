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
	private readonly Settings _settings = settings;
	private readonly ICache _cache = cache;
	[HttpGet]
	public IActionResult GetStatus()
	{
		var result = new Health()
		{
			Id = _cache.HealthId(),
			ApiName = _settings.Api.Name,
			ApiVersion = _settings.Api.Version,
			NodeIP = _cache.NodeIP(),
			NodeName = _cache.NodeName(),
			EnvName = _settings.Api.EnvName
		};
		return Ok(result);
	}
}
