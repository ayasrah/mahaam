using System.Transactions;
using Mahaam.Feat.Plans;
using Mahaam.Infra;

namespace Mahaam.Feat.Tasks;

public interface ITaskService
{
	Task<Guid> Create(Guid planId, string title);
	Task<List<Task>> GetList(Guid planId);
	System.Threading.Tasks.Task Delete(Guid planId, Guid id);
	System.Threading.Tasks.Task UpdateDone(Guid planId, Guid id, bool done);
	System.Threading.Tasks.Task UpdateTitle(Guid id, string title);
	System.Threading.Tasks.Task ReOrder(Guid planId, int oldOrder, int newOrder);
}

public class TaskService(ITaskRepo taskRepo, IPlanRepo planRepo) : ITaskService
{
	private readonly ITaskRepo _taskRepo = taskRepo;
	private readonly IPlanRepo _planRepo = planRepo;

	public async Task<Guid> Create(Guid planId, string title)
	{
		var count = await _taskRepo.GetCount(planId);
		if (count >= 100) throw new LogicException("Max is 100", "max_is_100");

		using var scope = new TransactionScope(TransactionScopeAsyncFlowOption.Enabled);
		var id = await _taskRepo.Create(planId, title);
		await _planRepo.UpdateDonePercent(planId);
		scope.Complete();
		return id;
	}

	public async Task<List<Task>> GetList(Guid planId)
	{
		return await _taskRepo.GetAll(planId);
	}

	public async System.Threading.Tasks.Task Delete(Guid planId, Guid id)
	{

		using var scope = new TransactionScope(TransactionScopeAsyncFlowOption.Enabled);
		await _taskRepo.UpdateOrderBeforeDelete(planId, id);
		await _taskRepo.DeleteOne(id);
		await _planRepo.UpdateDonePercent(planId);
		scope.Complete();
	}

	public async System.Threading.Tasks.Task UpdateDone(Guid planId, Guid id, bool done)
	{
		using var scope = new TransactionScope(TransactionScopeAsyncFlowOption.Enabled);
		await _taskRepo.UpdateDone(id, done);
		await _planRepo.UpdateDonePercent(planId);

		var count = await _taskRepo.GetCount(planId);
		var task = await _taskRepo.GetOne(id);
		await ReOrder(planId, task.SortOrder, done ? 0 : count - 1);
		scope.Complete();
	}

	public async System.Threading.Tasks.Task UpdateTitle(Guid id, string title)
	{
		await _taskRepo.UpdateTitle(id, title);
	}

	public async System.Threading.Tasks.Task ReOrder(Guid planId, int oldOrder, int newOrder)
	{
		if (oldOrder == newOrder) return;
		var count = await _taskRepo.GetCount(planId);
		if (oldOrder > count || newOrder > count)
			throw new InputException($"oldOrder and newOrder should be less than {count}");
		await _taskRepo.UpdateOrder(planId, oldOrder, newOrder);
	}
}
