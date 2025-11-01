package mahaam.infra.monitor;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import jakarta.ws.rs.GET;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.core.MediaType;
import jakarta.ws.rs.core.Response;
import mahaam.infra.Cache;
import mahaam.infra.Config;
import mahaam.infra.monitor.MonitorModel.Health;

public interface HealthController {
	Response getInfo();
}

@ApplicationScoped
@Path("/health")
class DefaultHealthController implements HealthController {

	@Inject
	Config config;

	@Inject
	Cache cache;

	@GET
	@Produces(MediaType.APPLICATION_JSON)
	public Response getInfo() {
		Health health = new Health();
		health.apiName = config.apiName();
		health.apiVersion = config.apiVersion();
		health.envName = config.envName();
		health.id = cache.healthId();
		health.nodeIP = cache.nodeIP();
		health.nodeName = cache.nodeName();
		return Response.ok(health).build();
	}
}
