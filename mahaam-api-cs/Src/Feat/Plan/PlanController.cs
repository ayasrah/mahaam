using System.Net.Mime;
using Mahaam.Infra;
using Microsoft.AspNetCore.Mvc;

namespace Mahaam.Feat.Plans;

public interface IPlanController
{
	IActionResult Create(PlanIn plan);
	IActionResult Update(PlanIn plan);
	IActionResult Delete(Guid id);
	IActionResult Share(Guid id, string email);
	IActionResult Unshare(Guid id, string email);
	IActionResult Leave(Guid id);
	IActionResult UpdateType(Guid id, string type);
	IActionResult ReOrder(string type, int oldOrder, int newOrder);
	IActionResult GetOne(Guid id);
	IActionResult GetMany(string? type);
}


[ApiController]
[Route("plans")]
public class PlanController : ControllerBase, IPlanController
{

	[HttpPost]
	[Consumes(MediaTypeNames.Application.Json)]
	public IActionResult Create([FromBody] PlanIn plan)
	{
		Rule.OneAtLeastRequired([plan.Title, plan.Starts, plan.Ends], "title or starts or ends is required");

		var id = App.PlanService.Create(plan);
		return Created($"/plans/{id}", id);
	}

	[HttpPut]
	[Consumes(MediaTypeNames.Application.Json)]
	public IActionResult Update([FromBody] PlanIn plan)
	{
		Rule.Required(plan.Id, "Id");
		Rule.OneAtLeastRequired([plan.Title, plan.Starts, plan.Ends], "title or starts or ends is required");

		App.PlanService.Update(plan);
		return Ok();
	}

	[HttpDelete]
	[Route("{id}")]
	public IActionResult Delete(Guid id)
	{
		Rule.Required(id, "id");
		App.PlanService.Delete(id);
		return NoContent();
	}

	[HttpPatch]
	[Route("{id}/share")]
	public IActionResult Share(Guid id, [FromForm] string email)
	{
		Rule.Required(id, "id");
		Rule.Required(email, "email");

		App.PlanService.Share(id, email);
		return Ok();
	}

	[HttpPatch]
	[Route("{id}/unshare")]
	public IActionResult Unshare(Guid id, [FromForm] string email)
	{
		Rule.Required(id, "id");
		Rule.Required(email, "email");

		App.PlanService.Unshare(id, email);
		return Ok();
	}

	[HttpPatch]
	[Route("{id}/leave")]
	public IActionResult Leave(Guid id)
	{
		Rule.Required(id, "id");
		App.PlanService.Leave(id);
		return Ok();
	}

	[HttpPatch]
	[Route("{id}/type")]
	public IActionResult UpdateType(Guid id, [FromForm] string type)
	{
		Rule.Required(id, "id");
		Rule.Required(type, "type");
		Rule.In(type, PlanType.All);

		App.PlanService.UpdateType(id, type);
		return Ok();
	}

	[HttpPatch]
	[Route("reorder")]
	public IActionResult ReOrder([FromForm] string type, [FromForm] int oldOrder, [FromForm] int newOrder)
	{
		Rule.Required(type, "type");
		Rule.In(type, PlanType.All);
		Rule.Required(oldOrder, "oldOrder");
		Rule.Required(newOrder, "newOrder");
		App.PlanService.ReOrder(type, oldOrder, newOrder);
		return Ok();
	}

	[HttpGet]
	[Route("{id}")]
	public IActionResult GetOne(Guid id)
	{
		Rule.Required(id, "id");

		var plan = App.PlanService.GetOne(id);
		return Ok(plan);
	}

	[HttpGet]
	[Route("")]
	public IActionResult GetMany([FromQuery] string? type)
	{
		if (type is not null) Rule.In(type, PlanType.All);
		else type = PlanType.Main;

		var plans = App.PlanService.GetMany(type);
		return Ok(plans);
	}
}
