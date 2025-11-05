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

	public async Task<Guid> Create(Guid planId, string title)
	{
		var count = await taskRepo.GetCount(planId);
		if (count >= 100) throw new LogicException("Max is 100", "max_is_100");

		using var scope = new TransactionScope(TransactionScopeAsyncFlowOption.Enabled);
		var id = await taskRepo.Create(planId, title);
		await planRepo.UpdateDonePercent(planId);
		scope.Complete();
		return id;
	}

	public async Task<List<Task>> GetList(Guid planId)
	{
		return await taskRepo.GetAll(planId);
	}

	public async System.Threading.Tasks.Task Delete(Guid planId, Guid id)
	{

		using var scope = new TransactionScope(TransactionScopeAsyncFlowOption.Enabled);
		await taskRepo.UpdateOrderBeforeDelete(planId, id);
		await taskRepo.DeleteOne(id);
		await planRepo.UpdateDonePercent(planId);
		scope.Complete();
	}

	public async System.Threading.Tasks.Task UpdateDone(Guid planId, Guid id, bool done)
	{
		using var scope = new TransactionScope(TransactionScopeAsyncFlowOption.Enabled);
		await taskRepo.UpdateDone(id, done);
		await planRepo.UpdateDonePercent(planId);

		var count = await taskRepo.GetCount(planId);
		var task = await taskRepo.GetOne(id);
		await ReOrder(planId, task.SortOrder, done ? 0 : count - 1);
		scope.Complete();
	}

	public async System.Threading.Tasks.Task UpdateTitle(Guid id, string title)
	{
		await taskRepo.UpdateTitle(id, title);
	}

	public async System.Threading.Tasks.Task ReOrder(Guid planId, int oldOrder, int newOrder)
	{
		if (oldOrder == newOrder) return;
		var count = await taskRepo.GetCount(planId);
		if (oldOrder > count || newOrder > count)
			throw new InputException($"oldOrder and newOrder should be less than {count}");
		await taskRepo.UpdateOrder(planId, oldOrder, newOrder);
	}
}
