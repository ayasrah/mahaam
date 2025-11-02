namespace Mahaam.Infra.Monitoring;

public interface IHealthRepo
{
	Task<int> Create(Health health);
	Task<int> UpdatePulse(Guid id);
	Task<int> UpdateStopped(Guid id);
}

public class HealthRepo(IDB db) : IHealthRepo
{
	private readonly IDB _db = db;
	public async Task<int> Create(Health health)
	{
		var query = @"
			INSERT INTO monitor.health (id, api_name, api_version, env_name, node_ip, node_name, started_at) 
			VALUES(@id, @apiName, @apiVersion, @envName, @nodeIP, @nodeName, current_timestamp)";
		return await _db.Insert(query, health);
	}

	public async Task<int> UpdatePulse(Guid id)
	{
		var query = "UPDATE monitor.health SET pulsed_at = current_timestamp WHERE id = @id";
		return await _db.Update(query, new { id });
	}

	public async Task<int> UpdateStopped(Guid id)
	{
		var query = "UPDATE monitor.health SET stopped_at = current_timestamp WHERE id = @id";
		return await _db.Update(query, new { id });
	}
}
