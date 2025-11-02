package mahaam.infra.monitor;

import java.util.UUID;
import java.util.concurrent.CompletableFuture;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import mahaam.infra.Cache;
import mahaam.infra.DB;
import mahaam.infra.Log;
import mahaam.infra.Mapper;

public interface LogRepo {
	void create(UUID trafficId, String type, String message);
}

@ApplicationScoped
class DefaultLogRepo implements LogRepo {

	@Inject
	DB db;

	@Inject
	Cache cache;

	@Override
	public void create(UUID trafficId, String type, String message) {
		CompletableFuture.runAsync(
				() -> {
					try {
						String query = """
								INSERT INTO monitor.logs (traffic_id, type, message, node_ip, created_at)
								VALUES (:trafficId, :type, :message, :node_ip, current_timestamp)
								""";

						var params = Mapper.of(
								"trafficId", trafficId,
								"type", type,
								"message", message,
								"node_ip", cache.nodeIP());

						db.insert(query, params);
					} catch (Exception ex) {
						Log.error(
								"Unable to create log record:("
										+ type
										+ ": "
										+ message
										+ "). Cause: "
										+ ex.toString());
					}
				});
	}
}
