package mahaam.infra.monitor;

import java.util.UUID;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import mahaam.infra.DB;
import mahaam.infra.Mapper;
import mahaam.infra.monitor.MonitorModel.Health;

public interface HealthRepo {
	int create(Health health);

	int updatePulse(UUID id);

	int updateStopped(UUID id);
}

@ApplicationScoped
class DefaultHealthRepo implements HealthRepo {

	@Inject
	DB db;

	@Override
	public int create(Health health) {
		String query = """
				INSERT INTO x_health (id, api_name, api_version, env_name, node_ip, node_name, started_at)
				VALUES (:id, :apiName, :apiVersion, :envName, :nodeIP, :nodeName, current_timestamp)
				""";
		var params = Mapper.of(
				"id", health.id,
				"apiName", health.apiName,
				"apiVersion", health.apiVersion,
				"envName", health.envName,
				"nodeIP", health.nodeIP,
				"nodeName", health.nodeName);

		return db.insert(query, params);
	}

	@Override
	public int updatePulse(UUID id) {
		String query = "UPDATE x_health SET pulsed_at = current_timestamp WHERE id = :id";
		return db.update(query, Mapper.of("id", id));
	}

	@Override
	public int updateStopped(UUID id) {
		String query = "UPDATE x_health SET stopped_at = current_timestamp WHERE id = :id";
		return db.update(query, Mapper.of("id", id));
	}
}
