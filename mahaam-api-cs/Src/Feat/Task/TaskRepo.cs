using Mahaam.Infra;

namespace Mahaam.Feat.Tasks;


public interface ITaskRepo
{
	Task<List<Task>> GetAll(Guid planId);
	Task<Task> GetOne(Guid id);
	Task<Guid> Create(Guid planId, string title);
	System.Threading.Tasks.Task DeleteOne(Guid id);
	System.Threading.Tasks.Task DeleteAll(Guid planId);
	System.Threading.Tasks.Task UpdateDone(Guid id, bool done);
	System.Threading.Tasks.Task UpdateTitle(Guid id, string title);
	System.Threading.Tasks.Task UpdateOrder(Guid planId, int oldOrder, int newOrder);
	System.Threading.Tasks.Task UpdateOrderBeforeDelete(Guid planId, Guid id);
	Task<int> GetCount(Guid planId);
}

public class TaskRepo(IDB db) : ITaskRepo
{
	private readonly IDB _db = db;
	public async Task<List<Task>> GetAll(Guid planId)
	{
		var query = @"SELECT id, plan_id, title, done, sort_order, created_at, updated_at 
			FROM tasks WHERE plan_id = @planId order by sort_order desc;";
		return await _db.SelectMany<Tasks.Task>(query, new { planId });
	}

	public async Task<Task> GetOne(Guid id)
	{
		var query = @"SELECT id, plan_id, title, done, sort_order, created_at, updated_at FROM tasks WHERE id = @id";
		return await _db.SelectOne<Tasks.Task>(query, new { id });
	}

	public async Task<Guid> Create(Guid planId, string title)
	{
		var query = @"INSERT INTO tasks (id, plan_id, title, done, sort_order, created_at) 
			VALUES (@id, @planId, @title, @done, (SELECT COUNT(1) FROM tasks WHERE plan_id = @planId), current_timestamp)";

		var id = Guid.NewGuid();
		await _db.Insert(query, new { id, planId, title, done = false });
		return id;
	}

	public async System.Threading.Tasks.Task DeleteOne(Guid id)
	{
		var query = "DELETE FROM tasks WHERE id = @id";
		var deletedRows = await _db.Delete(query, new { id });
		if (deletedRows == 0) throw new NotFoundException($"task id={id} not found");
	}

	public async System.Threading.Tasks.Task DeleteAll(Guid planId)
	{
		var query = "DELETE FROM tasks WHERE plan_id = @planId";
		await _db.Delete(query, new { planId });
	}

	public async System.Threading.Tasks.Task UpdateDone(Guid id, bool done)
	{
		var query = "UPDATE tasks SET done = @done, updated_at = current_timestamp WHERE id = @id";
		var updatedRows = await _db.Update(query, new { id, done });
		if (updatedRows == 0) throw new NotFoundException($"task id={id} not found");

	}

	public async System.Threading.Tasks.Task UpdateTitle(Guid id, string title)
	{
		var query = "UPDATE tasks SET title = @title, updated_at = current_timestamp WHERE id = @id";
		var updatedRows = await _db.Update(query, new { id, title });
		if (updatedRows == 0) throw new NotFoundException($"task id={id} not found");

	}

	public async System.Threading.Tasks.Task UpdateOrderBeforeDelete(Guid planId, Guid id)
	{
		var query = @"UPDATE tasks SET sort_order = sort_order - 1
			WHERE plan_id = @planId AND sort_order > (SELECT sort_order FROM tasks WHERE id =@id)";
		await _db.Update(query, new { planId, id });
	}

	public async System.Threading.Tasks.Task UpdateOrder(Guid planId, int oldOrder, int newOrder)
	{
		var query = @"
			UPDATE tasks SET sort_order = 
				CASE 
					WHEN sort_order = @oldOrder THEN @newOrder
					WHEN sort_order > @oldOrder AND sort_order <= @newOrder THEN sort_order - 1
					WHEN sort_order >= @newOrder AND sort_order < @oldOrder THEN sort_order + 1
					ELSE sort_order
				END
			WHERE plan_id = @planId;";
		await _db.Update(query, new { planId, oldOrder, newOrder });
	}

	public async Task<int> GetCount(Guid planId)
	{
		var query = "SELECT COUNT(1) FROM tasks WHERE plan_id = @planId";
		return await _db.SelectOne<int>(query, new { planId });
	}
}
