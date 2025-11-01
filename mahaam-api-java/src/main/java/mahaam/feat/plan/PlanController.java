package mahaam.feat.plan;

import java.util.Arrays;
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
import jakarta.ws.rs.PUT;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.PathParam;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.QueryParam;
import jakarta.ws.rs.core.MediaType;
import jakarta.ws.rs.core.Response;
import mahaam.feat.plan.PlanModel.Plan;
import mahaam.feat.plan.PlanModel.PlanIn;
import mahaam.feat.plan.PlanModel.PlanType;
import mahaam.infra.Json;
import mahaam.infra.Rule;

public interface PlanController {
	Response create(PlanIn plan);

	Response update(PlanIn plan);

	Response delete(UUID id);

	Response share(UUID id, String email);

	Response unshare(UUID id, String email);

	Response leave(UUID id);

	Response updateType(UUID id, String type);

	Response reOrder(String type, int oldOrder, int newOrder);

	Response getOne(UUID planId);

	Response getMany(String type);
}

@ApplicationScoped
@Path("/plans")
@Consumes(MediaType.APPLICATION_JSON)
@Produces(MediaType.APPLICATION_JSON)
class DefaultPlanController implements PlanController {

	@Inject
	PlanService planService;

	@POST
	public Response create(PlanIn plan) {
		Rule.required(plan, "plan");
		Rule.oneAtLeastRequired(Arrays.asList(plan.title, plan.starts, plan.ends),
				"title or starts or ends is required");

		UUID id = planService.create(plan);
		return Response.status(Response.Status.CREATED).entity(Json.toString(id)).build();
	}

	@PUT
	public Response update(PlanIn plan) {
		Rule.required(plan, "plan");
		Rule.required(plan.id, "Id");
		Rule.oneAtLeastRequired(Arrays.asList(plan.title, plan.starts, plan.ends),
				"title or starts or ends is required");

		planService.update(plan);
		return Response.ok().build();
	}

	@DELETE
	@Path("/{id}")
	public Response delete(@PathParam("id") UUID id) {
		Rule.required(id, "id");
		planService.delete(id);
		return Response.noContent().build();
	}

	@PATCH
	@Path("/{id}/share")
	@Consumes(MediaType.APPLICATION_FORM_URLENCODED)
	public Response share(@PathParam("id") UUID id, @FormParam("email") String email) {
		Rule.required(id, "id");
		Rule.required(email, "email");

		planService.share(id, email);
		return Response.ok().build();
	}

	@PATCH
	@Path("/{id}/unshare")
	@Consumes(MediaType.APPLICATION_FORM_URLENCODED)
	public Response unshare(@PathParam("id") UUID id, @FormParam("email") String email) {
		Rule.required(id, "id");
		Rule.required(email, "email");

		planService.unshare(id, email);
		return Response.ok().build();
	}

	@PATCH
	@Path("/{id}/leave")
	public Response leave(@PathParam("id") UUID id) {
		Rule.required(id, "id");
		planService.leave(id);
		return Response.ok().build();
	}

	@PATCH
	@Path("/{id}/type")
	@Consumes(MediaType.APPLICATION_FORM_URLENCODED)
	public Response updateType(@PathParam("id") UUID id, @FormParam("type") String type) {
		Rule.required(id, "id");
		Rule.required(type, "type");
		Rule.in(type, PlanType.ALL);

		planService.updateType(id, type);
		return Response.ok().build();
	}

	@PATCH
	@Path("/reorder")
	@Consumes(MediaType.APPLICATION_FORM_URLENCODED)
	public Response reOrder(@FormParam("type") String type, @FormParam("oldOrder") int oldOrder,
			@FormParam("newOrder") int newOrder) {
		Rule.required(type, "type");
		Rule.in(type, PlanType.ALL);
		Rule.required(oldOrder, "oldOrder");
		Rule.required(newOrder, "newOrder");
		planService.reOrder(type, oldOrder, newOrder);
		return Response.ok().build();
	}

	@GET
	@Path("/{planId}")
	public Response getOne(@PathParam("planId") UUID planId) {
		Rule.required(planId, "planId");

		Plan plan = planService.getOne(planId);
		return Response.ok().entity(Json.toString(plan)).build();
	}

	@GET
	public Response getMany(@QueryParam("type") String type) {
		if (type != null) {
			Rule.in(type, PlanType.ALL);
		} else {
			type = PlanType.MAIN;
		}

		List<Plan> plans = planService.getMany(type);
		return Response.ok().entity(Json.toString(plans)).build();
	}
}
