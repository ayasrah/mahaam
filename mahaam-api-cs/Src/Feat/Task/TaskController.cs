
using System.Net;
using Mahaam.Infra;
using Microsoft.AspNetCore.Mvc;

namespace Mahaam.Feat.Tasks;

public interface ITaskController
{
	IActionResult Create(Guid planId, string title);
	IActionResult Delete(Guid planId, Guid id);
	IActionResult UpdateDone(Guid planId, Guid id, bool done);
	IActionResult UpdateTitle(Guid id, string title);
	IActionResult ReOrder(Guid planId, int oldOrder, int newOrder);
	IActionResult GetMany(Guid planId);
}

[ApiController]
[Route("plans/{planId}/tasks")]
public class TaskController(ITaskService taskService) : ControllerBase, ITaskController
{
	private readonly ITaskService _taskService = taskService;
	[HttpPost]
	public IActionResult Create(Guid planId, [FromForm] string title)
	{
		Rule.Required(planId, "planId");
		Rule.Required(title, "title");
		var id = _taskService.Create(planId, title);
		return Created($"/plans/{planId}/tasks/{id}", id);
	}

	[HttpDelete]
	[Route("{id}")]
	public IActionResult Delete(Guid planId, Guid id)
	{
		Rule.Required(planId, "planId");
		Rule.Required(id, "id");
		_taskService.Delete(planId, id);
		return NoContent();
	}

	[HttpPatch]
	[Route("{id}/done")]
	public IActionResult UpdateDone(Guid planId, Guid id, [FromForm] bool done)
	{
		Rule.Required(planId, "planId");
		Rule.Required(id, "id");
		Rule.Required(done, "done");
		_taskService.UpdateDone(planId, id, done);
		return Ok();
	}

	[HttpPatch]
	[Route("{id}/title")]
	public IActionResult UpdateTitle(Guid id, [FromForm] string title)
	{
		Rule.Required(id, "id");
		Rule.Required(title, "title");
		_taskService.UpdateTitle(id, title);
		return Ok();
	}

	[HttpPatch]
	[Route("reorder")]
	public IActionResult ReOrder(Guid planId, [FromForm] int oldOrder, [FromForm] int newOrder)
	{
		Rule.Required(planId, "planId");
		Rule.Required(oldOrder, "oldOrder");
		Rule.Required(newOrder, "newOrder");
		_taskService.ReOrder(planId, oldOrder, newOrder);
		return Ok();
	}

	[HttpGet]
	[Route("")]
	public IActionResult GetMany(Guid planId)
	{
		Rule.Required(planId, "planId");
		var result = _taskService.GetList(planId);
		return Ok(result);
	}

}
