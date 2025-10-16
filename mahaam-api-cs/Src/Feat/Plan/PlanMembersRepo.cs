using Mahaam.Feat.Users;
using Mahaam.Infra;

namespace Mahaam.Feat.Plans;

public interface IPlanMembersRepo
{
	Task Create(Guid planId, Guid userId);
	Task<int> Delete(Guid planId, Guid userId);
	Task<List<Plan>> GetOtherPlans(Guid userId);
	Task<List<User>> GetUsers(Guid planId);
	Task<int> GetPlansCount(Guid userId);
	Task<int> GetUsersCount(Guid planId);
}

public class PlanMembersRepo(IDB db) : IPlanMembersRepo
{
	private readonly IDB _db = db;
	public async Task Create(Guid planId, Guid userId)
	{
		var query = @"INSERT INTO plan_members(plan_id, user_id, created_at) 
			VALUES(@planId, @userId, current_timestamp)";
		await _db.Insert(query, new { planId, userId });
	}

	public async Task<int> Delete(Guid planId, Guid userId)
	{
		var query = @"DELETE FROM plan_members WHERE plan_id = @planId AND user_id = @userId";
		return await _db.Delete(query, new { planId, userId });
	}


	/// <summary>
	/// Get other plans that are shared with userId
	/// </summary>
	public async Task<List<Plan>> GetOtherPlans(Guid userId)
	{
		var query = @"
			SELECT c.id, c.title, c.starts, c.ends, c.type, c.done_percent, c.sort_order, 
			c.created_at, u.id, u.email, u.name
			FROM plan_members cm
			LEFT JOIN plans c ON cm.plan_id = c.id
			LEFT JOIN users u ON c.user_id = u.id
			WHERE cm.user_id = @userId
			ORDER BY c.created_at DESC";

		return await _db.SelectMany<Plan, User, Plan>(query, (plan, user) =>
		{
			plan.User = user;
			plan.IsShared = true;
			return plan;
		}, new { userId }, splitOn: "id");
	}

	/// <summary>
	/// Get users that are sharing the plan
	/// </summary>
	public async Task<List<User>> GetUsers(Guid planId)
	{
		var query = @"SELECT u.id, u.email,u.name
			FROM plan_members cm
			LEFT JOIN users u ON cm.user_id = u.id
			WHERE cm.plan_id = @planId
			ORDER BY u.created_at DESC";
		return await _db.SelectMany<User>(query, new { planId });
	}

	/// <summary>
	/// Get counts of plans owned by the user that are shared
	/// </summary>
	public async Task<int> GetPlansCount(Guid userId)
	{
		var query = @"SELECT COUNT(1) FROM plan_members cm
			LEFT JOIN plans c ON cm.plan_id = c.id
			WHERE c.user_id = @userId";
		return await _db.SelectOne<int>(query, new { userId });
	}

	/// <summary>
	/// Get counts of users that are shared with the plan
	/// </summary>
	public async Task<int> GetUsersCount(Guid planId)
	{
		var query = "SELECT COUNT(1) FROM plan_members WHERE plan_id = @planId";
		return await _db.SelectOne<int>(query, new { planId });
	}
}