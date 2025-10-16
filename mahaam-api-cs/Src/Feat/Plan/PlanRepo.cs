using Mahaam.Feat.Users;
using Mahaam.Infra;

namespace Mahaam.Feat.Plans;

public interface IPlanRepo
{
	Plan GetOne(Guid id);
	List<Plan> GetMany(Guid userId, string type);
	Guid Create(PlanIn plan);
	void Update(PlanIn plan);
	void Delete(Guid id);
	void UpdateDonePercent(Guid id);
	void RemoveFromOrder(Guid userId, Guid id);
	void UpdateOrder(Guid userId, string type, int oldOrder, int newOrder);
	void UpdateType(Guid userId, Guid id, string type);
	int GetCount(Guid userId, string type);
	int UpdateUserId(Guid oldUserId, Guid newUserId);
}

public class PlanRepo(ILog log) : IPlanRepo
{
	private readonly ILog _log = log;
	public Guid Create(PlanIn plan)
	{
		var query = @"INSERT INTO plans (id, user_id, title, starts, ends, type, status, done_percent, sort_order, created_at)
			VALUES (@Id, @UserId, @Title, @Starts, @Ends, @type, @status, '0/0', 
			(SELECT COUNT(1) FROM plans WHERE user_id = @UserId AND type = @type), current_timestamp)";
		var id = Guid.NewGuid();
		DB.Insert(query, new
		{
			id,
			plan.Title,
			plan.Starts,
			plan.Ends,
			UserId = Req.UserId,
			type = PlanType.Main,
			status = PlanStatus.Open
		});
		return id;
	}

	public void Update(PlanIn plan)
	{
		var query = "UPDATE plans SET title = @title, starts = @starts, ends = @ends, updated_at = current_timestamp WHERE id = @id";
		DB.Update(query, new { id = plan.Id, title = plan.Title, starts = plan.Starts, ends = plan.Ends });
	}

	public Plan GetOne(Guid id)
	{
		var query = @"
			SELECT c.id, c.title, c.starts, c.ends, c.type, c.done_percent, c.sort_order, c.user_id,
				EXISTS(SELECT 1 FROM plan_members cm WHERE cm.plan_id = c.id) AS IsShared,
				u.id, u.email, u.name
			FROM plans c
			LEFT JOIN users u ON c.user_id = u.id
			WHERE c.id = @id";

		return DB.SelectOne<Plan, User, Plan>(
			query,
			(plan, user) =>
			{
				plan.User = user;
				return plan;
			},
			new { id }
		);
	}


	/// <summary>
	/// Get all plans that created by userId, for a given type, wether they are shared or not
	/// </summary>
	public List<Plan> GetMany(Guid userId, string type)
	{
		var query = @"
			SELECT c.id, c.title, c.starts, c.ends, c.type, c.done_percent, c.sort_order, c.user_id,
				EXISTS(SELECT 1 FROM plan_members cm WHERE cm.plan_id = c.id) AS IsShared,
				u.id, u.email, u.name
			FROM plans c
			LEFT JOIN users u ON c.user_id = u.id
			WHERE c.user_id = @userId AND c.type = @type
			ORDER BY c.sort_order DESC;";

		return DB.SelectMany<Plan, User, Plan>(
			query,
			(plan, user) =>
			{
				plan.User = user;
				return plan;
			},
			new { userId, type }
		);
	}

	public void Delete(Guid id)
	{
		var count = DB.Delete("DELETE FROM plans WHERE id = @id", new { id });
		if (count > 0) _log.Info($"Plan {id} deleted");
	}

	public void UpdateDonePercent(Guid id)
	{
		var query = "SELECT * FROM tasks WHERE plan_id = @id";
		var tasks = DB.SelectMany<Tasks.Task>(query, new { id });

		var done = tasks.Where(task => task.Done).Count();
		var notDone = tasks.Count;
		var donePercent = $"{done}/{notDone}";
		var updatequery = "UPDATE plans SET done_percent = @donePercent WHERE id = @id";
		DB.Update(updatequery, new { donePercent, id });
	}

	/// <summary>
	/// Update sort_order of Plans per userId and Plan Type after removing PlanId (id)
	/// </summary>
	/// <param name="userId"></param>
	/// <param name="id"></param>
	public void RemoveFromOrder(Guid userId, Guid id)
	{
		var query = @"
			UPDATE plans SET sort_order = sort_order - 1 
			WHERE user_id = @userId AND type = (SELECT type FROM Plans WHERE id =@id) 
				AND sort_order > (SELECT sort_order FROM plans WHERE id =@id)";
		DB.Update(query, new { userId, id });
	}

	public void UpdateOrder(Guid userId, string type, int oldOrder, int newOrder)
	{
		var query = @"
			UPDATE plans SET sort_order = 
				CASE 
					WHEN sort_order = @oldOrder THEN @newOrder
					WHEN sort_order > @oldOrder AND sort_order <= @newOrder THEN sort_order - 1
					WHEN sort_order >= @newOrder AND sort_order < @oldOrder THEN sort_order + 1
					ELSE sort_order
				END
			WHERE 
				user_id = @userId AND 
				type = @type";
		var updated = DB.Update(query, new { userId, type, oldOrder, newOrder });
		Console.WriteLine(updated);
	}

	public void UpdateType(Guid userId, Guid id, string type)
	{
		var query = @"UPDATE plans SET type = @type, 
			sort_order = (SELECT COUNT(1) FROM plans WHERE user_id = @userId AND type = @type), 
			updated_at = current_timestamp WHERE id = @id";
		DB.Update(query, new { userId, id, type });
	}

	public int GetCount(Guid userId, string type)
	{
		var queryCount = "SELECT COUNT(*) FROM plans WHERE user_id = @userId and type = @type";
		var count = DB.SelectOne<int>(queryCount, new { userId, type });
		return count;
	}

	public int UpdateUserId(Guid oldUserId, Guid newUserId)
	{
		var query = @"
		 	UPDATE plans SET user_id = @newUserId,
			sort_order = (sort_order + (Select count(1) from plans where user_id=@newUserId)),
			updated_at = current_timestamp 
			WHERE user_id = @oldUserId";
		return DB.Update(query, new { oldUserId, newUserId });
	}
}
