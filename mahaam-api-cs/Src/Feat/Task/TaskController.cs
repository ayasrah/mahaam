
using System.Net;
using Mahaam.Infra;
using Microsoft.AspNetCore.Mvc;

namespace Mahaam.Feat.Tasks;

public interface ITaskController
{
	Task<IActionResult> Create(Guid planId, string title);
	Task<IActionResult> Delete(Guid planId, Guid id);
	Task<IActionResult> UpdateDone(Guid planId, Guid id, bool done);
	Task<IActionResult> UpdateTitle(Guid id, string title);
	Task<IActionResult> ReOrder(Guid planId, int oldOrder, int newOrder);
	Task<IActionResult> GetMany(Guid planId);
}

[ApiController]
[Route("plans/{planId}/tasks")]
public class TaskController(ITaskService taskService) : ControllerBase, ITaskController
{
	private readonly ITaskService _taskService = taskService;
	[HttpPost]
	public async Task<IActionResult> Create(Guid planId, [FromForm] string title)
	{
		Rule.Required(planId, "planId");
		Rule.Required(title, "title");
		var id = await _taskService.Create(planId, title);
		return Created($"/plans/{planId}/tasks/{id}", id);
	}

	[HttpDelete]
	[Route("{id}")]
	public async Task<IActionResult> Delete(Guid planId, Guid id)
	{
		Rule.Required(planId, "planId");
		Rule.Required(id, "id");
		await _taskService.Delete(planId, id);
		return NoContent();
	}

	[HttpPatch]
	[Route("{id}/done")]
	public async Task<IActionResult> UpdateDone(Guid planId, Guid id, [FromForm] bool done)
	{
		Rule.Required(planId, "planId");
		Rule.Required(id, "id");
		Rule.Required(done, "done");
		await _taskService.UpdateDone(planId, id, done);
		return Ok();
	}

	[HttpPatch]
	[Route("{id}/title")]
	public async Task<IActionResult> UpdateTitle(Guid id, [FromForm] string title)
	{
		Rule.Required(id, "id");
		Rule.Required(title, "title");
		await _taskService.UpdateTitle(id, title);
		return Ok();
	}

	[HttpPatch]
	[Route("reorder")]
	public async Task<IActionResult> ReOrder(Guid planId, [FromForm] int oldOrder, [FromForm] int newOrder)
	{
		Rule.Required(planId, "planId");
		Rule.Required(oldOrder, "oldOrder");
		Rule.Required(newOrder, "newOrder");
		await _taskService.ReOrder(planId, oldOrder, newOrder);
		return Ok();
	}

	[HttpGet]
	[Route("")]
	public async Task<IActionResult> GetMany(Guid planId)
	{
		Rule.Required(planId, "planId");
		var result = await _taskService.GetList(planId);
		return Ok(result);
	}

}
