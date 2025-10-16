using Mahaam.Infra;

namespace Mahaam.Feat.Tasks;

public interface ITaskRepo
{
	List<Task> GetAll(Guid planId);
	Task GetOne(Guid id);
	Guid Create(Guid planId, string title);
	void DeleteOne(Guid id);
	void DeleteAll(Guid planId);
	void UpdateDone(Guid id, bool done);
	void UpdateTitle(Guid id, string title);
	void UpdateOrder(Guid planId, int oldOrder, int newOrder);
	void UpdateOrderBeforeDelete(Guid planId, Guid id);
	int GetCount(Guid planId);
}

public class TaskRepo(IDB db) : ITaskRepo
{
	private readonly IDB _db = db;
	public List<Task> GetAll(Guid planId)
	{
		var query = @"SELECT id, plan_id, title, done, sort_order, created_at, updated_at 
			FROM tasks WHERE plan_id = @planId order by sort_order desc;";
		return _db.SelectMany<Task>(query, new { planId });
	}

	public Task GetOne(Guid id)
	{
		var query = @"SELECT id, plan_id, title, done, sort_order, created_at, updated_at FROM tasks WHERE id = @id";
		return _db.SelectOne<Task>(query, new { id });
	}

	public Guid Create(Guid planId, string title)
	{
		var query = @"INSERT INTO tasks (id, plan_id, title, done, sort_order, created_at) 
			VALUES (@id, @planId, @title, @done, (SELECT COUNT(1) FROM tasks WHERE plan_id = @planId), current_timestamp)";

		var id = Guid.NewGuid();
		_db.Insert(query, new { id, planId, title, done = false });
		return id;
	}

	public void DeleteOne(Guid id)
	{
		var query = "DELETE FROM tasks WHERE id = @id";
		var deletedRows = _db.Delete(query, new { id });
		if (deletedRows == 0) throw new NotFoundException($"task id={id} not found");
	}

	public void DeleteAll(Guid planId)
	{
		var query = "DELETE FROM tasks WHERE plan_id = @planId";
		_db.Delete(query, new { planId });
	}

	public void UpdateDone(Guid id, bool done)
	{
		var query = "UPDATE tasks SET done = @done, updated_at = current_timestamp WHERE id = @id";
		var updatedRows = _db.Update(query, new { id, done });
		if (updatedRows == 0) throw new NotFoundException($"task id={id} not found");

	}

	public void UpdateTitle(Guid id, string title)
	{
		var query = "UPDATE tasks SET title = @title, updated_at = current_timestamp WHERE id = @id";
		var updatedRows = _db.Update(query, new { id, title });
		if (updatedRows == 0) throw new NotFoundException($"task id={id} not found");

	}

	public void UpdateOrderBeforeDelete(Guid planId, Guid id)
	{
		var query = @"UPDATE tasks SET sort_order = sort_order - 1
			WHERE plan_id = @planId AND sort_order > (SELECT sort_order FROM tasks WHERE id =@id)";
		_db.Update(query, new { planId, id });
	}

	public void UpdateOrder(Guid planId, int oldOrder, int newOrder)
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
		_db.Update(query, new { planId, oldOrder, newOrder });
	}

	public int GetCount(Guid planId)
	{
		var query = "SELECT COUNT(1) FROM tasks WHERE plan_id = @planId";
		return _db.SelectOne<int>(query, new { planId });
	}
}
