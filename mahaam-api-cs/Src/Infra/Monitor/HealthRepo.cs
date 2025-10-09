namespace Mahaam.Infra.Monitoring;

public interface IHealthRepo
{
	int Create(Health health);
	int UpdatePulse(Guid id);
	int UpdateStopped(Guid id);
}

public class HealthRepo : IHealthRepo
{
	public int Create(Health health)
	{
		var query = @"
			INSERT INTO x_health (id, api_name, api_version, env_name, node_ip, node_name, started_at) 
			VALUES(@id, @apiName, @apiVersion, @envName, @nodeIP, @nodeName, current_timestamp)";
		return DB.Insert(query, health);
	}

	public int UpdatePulse(Guid id)
	{
		var query = "UPDATE x_health SET pulsed_at = current_timestamp WHERE id = @id";
		return DB.Update(query, new { id });
	}

	public int UpdateStopped(Guid id)
	{
		var query = "UPDATE x_health SET stopped_at = current_timestamp WHERE id = @id";
		return DB.Update(query, new { id });
	}
}
