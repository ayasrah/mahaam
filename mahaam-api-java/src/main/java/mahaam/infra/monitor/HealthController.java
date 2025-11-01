package mahaam.infra.monitor;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.ws.rs.GET;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.core.MediaType;
import jakarta.ws.rs.core.Response;
import mahaam.infra.Cache;
import mahaam.infra.monitor.MonitorModel.Health;

public interface HealthController {
	Response getInfo();
}

@ApplicationScoped
@Path("/health")
class DefaultHealthController implements HealthController {

	@GET
	@Produces(MediaType.APPLICATION_JSON)
	public Response getInfo() {
		Health health = new Health();
		health.apiName = Cache.getApiName();
		health.apiVersion = Cache.getApiVersion();
		health.envName = Cache.getEnvName();
		health.id = Cache.getHealthId();
		health.nodeIP = Cache.getNodeIP();
		health.nodeName = Cache.getNodeName();
		return Response.ok(health).build();
	}
}
