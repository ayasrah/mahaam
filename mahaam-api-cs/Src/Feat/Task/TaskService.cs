using System.Transactions;
using Mahaam.Feat.Plans;
using Mahaam.Infra;

namespace Mahaam.Feat.Tasks;

public interface ITaskService
{
	Guid Create(Guid planId, string title);
	List<Task> GetList(Guid planId);
	void Delete(Guid planId, Guid id);
	void UpdateDone(Guid planId, Guid id, bool done);
	void UpdateTitle(Guid id, string title);
	void ReOrder(Guid planId, int oldOrder, int newOrder);
}

public class TaskService(ITaskRepo taskRepo, IPlanRepo planRepo) : ITaskService
{
	private readonly ITaskRepo _taskRepo = taskRepo;
	private readonly IPlanRepo _planRepo = planRepo;

	public Guid Create(Guid planId, string title)
	{
		var count = _taskRepo.GetCount(planId);
		if (count >= 100) throw new LogicException("Max is 100", "max_is_100");

		using var scope = new TransactionScope();
		var id = _taskRepo.Create(planId, title);
		_planRepo.UpdateDonePercent(planId);
		scope.Complete();
		return id;
	}

	public List<Task> GetList(Guid planId)
	{
		return _taskRepo.GetAll(planId);
	}

	public void Delete(Guid planId, Guid id)
	{

		using var scope = new TransactionScope();
		_taskRepo.UpdateOrderBeforeDelete(planId, id);
		_taskRepo.DeleteOne(id);
		_planRepo.UpdateDonePercent(planId);
		scope.Complete();
	}

	public void UpdateDone(Guid planId, Guid id, bool done)
	{
		using var scope = new TransactionScope();
		_taskRepo.UpdateDone(id, done);
		_planRepo.UpdateDonePercent(planId);

		var count = _taskRepo.GetCount(planId);
		var task = _taskRepo.GetOne(id);
		ReOrder(planId, task.SortOrder, done ? 0 : count - 1);
		scope.Complete();
	}

	public void UpdateTitle(Guid id, string title)
	{
		_taskRepo.UpdateTitle(id, title);
	}

	public void ReOrder(Guid planId, int oldOrder, int newOrder)
	{
		if (oldOrder == newOrder) return;
		var count = _taskRepo.GetCount(planId);
		if (oldOrder > count || newOrder > count)
			throw new InputException($"oldOrder and newOrder should be less than {count}");
		_taskRepo.UpdateOrder(planId, oldOrder, newOrder);
	}
}
