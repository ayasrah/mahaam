using System.Transactions;
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

public class TaskService : ITaskService
{
	public Guid Create(Guid planId, string title)
	{
		var count = App.TaskRepo.GetCount(planId);
		if (count >= 100) throw new LogicException("Max is 100", "max_is_100");

		using var scope = new TransactionScope();
		var id = App.TaskRepo.Create(planId, title);
		App.PlanRepo.UpdateDonePercent(planId);
		scope.Complete();
		return id;
	}

	public List<Task> GetList(Guid planId)
	{
		return App.TaskRepo.GetAll(planId);
	}

	public void Delete(Guid planId, Guid id)
	{

		using var scope = new TransactionScope();
		App.TaskRepo.UpdateOrderBeforeDelete(planId, id);
		App.TaskRepo.DeleteOne(id);
		App.PlanRepo.UpdateDonePercent(planId);
		scope.Complete();
	}

	public void UpdateDone(Guid planId, Guid id, bool done)
	{
		using var scope = new TransactionScope();
		App.TaskRepo.UpdateDone(id, done);
		App.PlanRepo.UpdateDonePercent(planId);

		var count = App.TaskRepo.GetCount(planId);
		var task = App.TaskRepo.GetOne(id);
		ReOrder(planId, task.SortOrder, done ? 0 : count - 1);
		scope.Complete();
	}

	public void UpdateTitle(Guid id, string title)
	{
		App.TaskRepo.UpdateTitle(id, title);
	}

	public void ReOrder(Guid planId, int oldOrder, int newOrder)
	{
		if (oldOrder == newOrder) return;
		var count = App.TaskRepo.GetCount(planId);
		if (oldOrder > count || newOrder > count)
			throw new InputException($"oldOrder and newOrder should be less than {count}");
		App.TaskRepo.UpdateOrder(planId, oldOrder, newOrder);
	}
}
