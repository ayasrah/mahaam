using System.Net.Mime;
using Mahaam.Infra;
using Microsoft.AspNetCore.Mvc;

namespace Mahaam.Feat.Plans;

public interface IPlanController
{
	Task<IActionResult> Create(PlanIn plan);
	Task<IActionResult> Update(PlanIn plan);
	Task<IActionResult> Delete(Guid id);
	Task<IActionResult> Share(Guid id, string email);
	Task<IActionResult> Unshare(Guid id, string email);
	Task<IActionResult> Leave(Guid id);
	Task<IActionResult> UpdateType(Guid id, string type);
	Task<IActionResult> ReOrder(string type, int oldOrder, int newOrder);
	Task<IActionResult> GetOne(Guid id);
	Task<IActionResult> GetMany(string? type);
}


[ApiController]
[Route("plans")]
public class PlanController(IPlanService planService) : ControllerBase, IPlanController
{
	private readonly IPlanService _planService = planService;

	[HttpPost]
	[Consumes(MediaTypeNames.Application.Json)]
	public async Task<IActionResult> Create([FromBody] PlanIn plan)
	{
		Rule.OneAtLeastRequired([plan.Title, plan.Starts, plan.Ends], "title or starts or ends is required");

		var id = await _planService.Create(plan);
		return Created($"/plans/{id}", id);
	}

	[HttpPut]
	[Consumes(MediaTypeNames.Application.Json)]
	public async Task<IActionResult> Update([FromBody] PlanIn plan)
	{
		Rule.Required(plan.Id, "Id");
		Rule.OneAtLeastRequired([plan.Title, plan.Starts, plan.Ends], "title or starts or ends is required");

		await _planService.Update(plan);
		return Ok();
	}

	[HttpDelete]
	[Route("{id}")]
	public async Task<IActionResult> Delete(Guid id)
	{
		Rule.Required(id, "id");
		await _planService.Delete(id);
		return NoContent();
	}

	[HttpPatch]
	[Route("{id}/share")]
	public async Task<IActionResult> Share(Guid id, [FromForm] string email)
	{
		Rule.Required(id, "id");
		Rule.Required(email, "email");

		await _planService.Share(id, email);
		return Ok();
	}

	[HttpPatch]
	[Route("{id}/unshare")]
	public async Task<IActionResult> Unshare(Guid id, [FromForm] string email)
	{
		Rule.Required(id, "id");
		Rule.Required(email, "email");

		await _planService.Unshare(id, email);
		return Ok();
	}

	[HttpPatch]
	[Route("{id}/leave")]
	public async Task<IActionResult> Leave(Guid id)
	{
		Rule.Required(id, "id");
		await _planService.Leave(id);
		return Ok();
	}

	[HttpPatch]
	[Route("{id}/type")]
	public async Task<IActionResult> UpdateType(Guid id, [FromForm] string type)
	{
		Rule.Required(id, "id");
		Rule.Required(type, "type");
		Rule.In(type, PlanType.All);

		await _planService.UpdateType(id, type);
		return Ok();
	}

	[HttpPatch]
	[Route("reorder")]
	public async Task<IActionResult> ReOrder([FromForm] string type, [FromForm] int oldOrder, [FromForm] int newOrder)
	{
		Rule.Required(type, "type");
		Rule.In(type, PlanType.All);
		Rule.Required(oldOrder, "oldOrder");
		Rule.Required(newOrder, "newOrder");
		await _planService.ReOrder(type, oldOrder, newOrder);
		return Ok();
	}

	[HttpGet]
	[Route("{id}")]
	public async Task<IActionResult> GetOne(Guid id)
	{
		Rule.Required(id, "id");

		var plan = await _planService.GetOne(id);
		return Ok(plan);
	}

	[HttpGet]
	[Route("")]
	public async Task<IActionResult> GetMany([FromQuery] string? type)
	{
		if (type is not null) Rule.In(type, PlanType.All);
		else type = PlanType.Main;

		var plans = await _planService.GetMany(type);
		return Ok(plans);
	}
}
