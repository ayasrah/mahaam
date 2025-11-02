package mahaam.infra.monitor;

import java.util.concurrent.CompletableFuture;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import mahaam.infra.DB;
import mahaam.infra.Log;
import mahaam.infra.Mapper;
import mahaam.infra.monitor.MonitorModel.Traffic;

public interface TrafficRepo {
	void create(Traffic traffic);
}

@ApplicationScoped
class DefaultTrafficRepo implements TrafficRepo {

	@Inject
	DB db;

	@Override
	public void create(Traffic traffic) {
		CompletableFuture.runAsync(
				() -> {
					try {
						String query = """
								INSERT INTO monitor.traffic (id, health_id, method, path, code, elapsed, headers, request, response, created_at)
								VALUES (:id, :healthId, :method, :path, :code, :elapsed, :headers, :request, :response, current_timestamp)
								""";

						var params = Mapper.of(
								"id", traffic.id,
								"healthId", traffic.healthId,
								"method", traffic.method,
								"path", traffic.path,
								"code", traffic.code,
								"elapsed", traffic.elapsed,
								"headers", traffic.headers,
								"request", traffic.request,
								"response", traffic.response);

						db.insert(query, params);
					} catch (Exception e) {
						Log.error("error creating traffic record: " + e);
					}
				});
	}
}
