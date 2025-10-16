using Mahaam.Feat.Users;
using Mahaam.Infra;

namespace Mahaam.Feat.Plans;

public interface IPlanMembersRepo
{
	void Create(Guid planId, Guid userId);
	int Delete(Guid planId, Guid userId);
	List<Plan> GetOtherPlans(Guid userId);
	List<User> GetUsers(Guid planId);
	public int GetPlansCount(Guid userId);
	public int GetUsersCount(Guid planId);
}

public class PlanMembersRepo(IDB db) : IPlanMembersRepo
{
	private readonly IDB _db = db;
	public void Create(Guid planId, Guid userId)
	{
		var query = @"INSERT INTO plan_members(plan_id, user_id, created_at) 
			VALUES(@planId, @userId, current_timestamp)";
		_db.Insert(query, new { planId, userId });
	}

	public int Delete(Guid planId, Guid userId)
	{
		var query = @"DELETE FROM plan_members WHERE plan_id = @planId AND user_id = @userId";
		return _db.Delete(query, new { planId, userId });
	}


	/// <summary>
	/// Get other plans that are shared with userId
	/// </summary>
	public List<Plan> GetOtherPlans(Guid userId)
	{
		var query = @"
			SELECT c.id, c.title, c.starts, c.ends, c.type, c.done_percent, c.sort_order, 
			c.created_at, u.id, u.email, u.name
			FROM plan_members cm
			LEFT JOIN plans c ON cm.plan_id = c.id
			LEFT JOIN users u ON c.user_id = u.id
			WHERE cm.user_id = @userId
			ORDER BY c.created_at DESC";

		return _db.SelectMany<Plan, User, Plan>(query, (plan, user) =>
		{
			plan.User = user;
			plan.IsShared = true;
			return plan;
		}, new { userId }, splitOn: "id");
	}

	/// <summary>
	/// Get users that are sharing the plan
	/// </summary>
	public List<User> GetUsers(Guid planId)
	{
		var query = @"SELECT u.id, u.email,u.name
			FROM plan_members cm
			LEFT JOIN users u ON cm.user_id = u.id
			WHERE cm.plan_id = @planId
			ORDER BY u.created_at DESC";
		return _db.SelectMany<User>(query, new { planId });
	}

	/// <summary>
	/// Get counts of plans owned by the user that are shared
	/// </summary>
	public int GetPlansCount(Guid userId)
	{
		var query = @"SELECT COUNT(1) FROM plan_members cm
			LEFT JOIN plans c ON cm.plan_id = c.id
			WHERE c.user_id = @userId";
		return _db.SelectOne<int>(query, new { userId });
	}

	/// <summary>
	/// Get counts of users that are shared with the plan
	/// </summary>
	public int GetUsersCount(Guid planId)
	{
		var query = "SELECT COUNT(1) FROM plan_members WHERE plan_id = @planId";
		return _db.SelectOne<int>(query, new { planId });
	}
}