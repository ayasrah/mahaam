package mahaam.infra.monitor;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import jakarta.ws.rs.FormParam;
import jakarta.ws.rs.POST;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.core.Response;
import mahaam.infra.Config;
import mahaam.infra.Log;

public interface AuditController {
	Response error(String error);

	Response trace(String info);
}

@ApplicationScoped
@Path("/audit")
class DefaultAuditController implements AuditController {

	@Inject
	Config config;

	@POST
	@Path("/error")
	@Override
	public Response error(@FormParam("error") String error) {
		Log.error("mahaam-mb: " + error);
		return Response.status(Response.Status.CREATED).build();
	}

	@POST
	@Path("/info")
	@Override
	public Response trace(@FormParam("info") String info) {
		Log.info("mahaam-mb: " + info);
		return Response.status(Response.Status.CREATED).build();
	}
}
