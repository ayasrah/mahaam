package mahaam.feat.task;

import java.util.List;
import java.util.UUID;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import jakarta.ws.rs.Consumes;
import jakarta.ws.rs.DELETE;
import jakarta.ws.rs.FormParam;
import jakarta.ws.rs.GET;
import jakarta.ws.rs.PATCH;
import jakarta.ws.rs.POST;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.PathParam;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.core.MediaType;
import jakarta.ws.rs.core.Response;
import mahaam.infra.Json;
import mahaam.infra.Rule;

public interface TaskController {
	Response create(UUID planId, String title);

	Response delete(UUID planId, UUID id);

	Response updateDone(UUID planId, UUID id, Boolean done);

	Response updateTitle(UUID id, String title);

	Response reOrder(UUID planId, int oldOrder, int newOrder);

	Response getMany(UUID planId);
}

@ApplicationScoped
@Path("/plans/{planId}/tasks")
@Consumes(MediaType.APPLICATION_FORM_URLENCODED)
@Produces(MediaType.APPLICATION_JSON)
class DefaultTaskController implements TaskController {

	@Inject
	TaskService taskService;

	@POST
	public Response create(@PathParam("planId") UUID planId, @FormParam("title") String title) {
		Rule.required(planId, "planId");
		Rule.required(title, "title");

		UUID id = taskService.create(planId, title);
		return Response.status(Response.Status.CREATED).entity(Json.toString(id)).build();
	}

	@DELETE
	@Path("/{id}")
	public Response delete(@PathParam("planId") UUID planId, @PathParam("id") UUID id) {
		Rule.required(planId, "planId");
		Rule.required(id, "id");

		taskService.delete(planId, id);
		return Response.noContent().build();
	}

	@PATCH
	@Path("/{id}/done")
	public Response updateDone(
			@PathParam("planId") UUID planId,
			@PathParam("id") UUID id,
			@FormParam("done") Boolean done) {
		Rule.required(planId, "planId");
		Rule.required(id, "id");
		Rule.required(done, "done");

		taskService.updateDone(planId, id, done);
		return Response.ok().build();
	}

	@PATCH
	@Path("/{id}/title")
	public Response updateTitle(@PathParam("id") UUID id, @FormParam("title") String title) {
		Rule.required(id, "id");
		Rule.required(title, "title");

		taskService.updateTitle(id, title);
		return Response.ok().build();
	}

	@PATCH
	@Path("/reorder")
	public Response reOrder(
			@PathParam("planId") UUID planId,
			@FormParam("oldOrder") int oldOrder,
			@FormParam("newOrder") int newOrder) {
		Rule.required(planId, "planId");
		Rule.required(oldOrder, "oldOrder");
		Rule.required(newOrder, "newOrder");

		taskService.reOrder(planId, oldOrder, newOrder);
		return Response.ok().build();
	}

	@GET
	public Response getMany(@PathParam("planId") UUID planId) {
		Rule.required(planId, "planId");

		List<Task> result = taskService.getList(planId);
		return Response.ok(Json.toString(result)).build();
	}
}
