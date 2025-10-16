using System.Transactions;
using Mahaam.Feat.Users;
using Mahaam.Infra;

namespace Mahaam.Feat.Plans;
public interface IPlanService
{
	Task<Plan> GetOne(Guid planId);
	Task<List<Plan>> GetMany(string type);
	Task<Guid> Create(PlanIn plan);
	Task Update(PlanIn plan);
	Task Delete(Guid id);
	Task Share(Guid id, string email);
	Task Unshare(Guid id, string email);
	Task Leave(Guid id);
	Task UpdateType(Guid id, string type);
	Task ReOrder(string type, int oldOrder, int newOrder);
	Task ValidateUserOwnsThePlan(Guid planId);
}

public class PlanService(IPlanRepo planRepo, IPlanMembersRepo planMembersRepo, IUserRepo userRepo, ISuggestedEmailsRepo suggestedEmailsRepo, ILog log) : IPlanService
{
	private readonly IPlanRepo _planRepo = planRepo;
	private readonly IPlanMembersRepo _planMembersRepo = planMembersRepo;
	private readonly IUserRepo _userRepo = userRepo;
	private readonly ISuggestedEmailsRepo _suggestedEmailsRepo = suggestedEmailsRepo;
	private readonly ILog _log = log;

	public async Task<Plan> GetOne(Guid planId)
	{
		var plan = await _planRepo.GetOne(planId);
		if (plan is { IsShared: true }) plan.Members = await _planMembersRepo.GetUsers(planId);
		return plan;
	}

	public async Task<List<Plan>> GetMany(string type)
	{
		// plans of the user shared or not
		var plans = await _planRepo.GetMany(Req.UserId, type);
		if (Req.IsLoggedIn)
		{
			// plans of others shared with the user
			var sharedPlans = await _planMembersRepo.GetOtherPlans(Req.UserId);
			plans.AddRange(sharedPlans);
		}
		return plans;
	}

	public async Task<Guid> Create(PlanIn plan)
	{
		var userId = Req.UserId;
		var count = await _planRepo.GetCount(userId, "Main");
		if (count >= 100) throw new LogicException("max_is_100", "Max is 100");

		return await _planRepo.Create(plan);
	}

	public async Task Update(PlanIn plan)
	{
		await ValidateUserOwnsThePlan(plan.Id);
		await _planRepo.Update(plan);
	}

	public async Task Delete(Guid id)
	{
		await ValidateUserOwnsThePlan(id);
		using var scope = new TransactionScope(TransactionScopeAsyncFlowOption.Enabled);
		await _planRepo.RemoveFromOrder(Req.UserId, id);
		await _planRepo.Delete(id); // This will delete all related records as it is a cascade delete
		scope.Complete();
	}

	public async Task Share(Guid id, string email)
	{
		await ValidateUserLoggedIn();
		await ValidateUserOwnsThePlan(id);
		var user = await _userRepo.GetOne(email) ?? throw new ArgumentException("email_not_found", $"email:{email} was not found");

		if (user.Id.Equals(Req.UserId))
			throw new LogicException("not_allowed_to_share_with_creator", "Not allowed to share with creator");

		var limit = 20;
		var plan = await _planRepo.GetOne(id);
		if (plan is { IsShared: true })
		{
			var membersCount = await _planMembersRepo.GetUsersCount(id);
			if (membersCount >= limit) throw new LogicException("max_is_20", "Max is 20");
		}
		else
		{
			var plansCount = await _planMembersRepo.GetPlansCount(Req.UserId);
			if (plansCount >= limit) throw new LogicException("max_is_20", "Max is 20");
		}
		await _planMembersRepo.Create(id, user.Id);
		await _suggestedEmailsRepo.Create(Req.UserId, email);
		var usr = await _userRepo.GetOne(Req.UserId);
		await _suggestedEmailsRepo.Create(user.Id, usr!.Email!);
	}

	public async Task Unshare(Guid id, string email)
	{
		await ValidateUserLoggedIn();
		await ValidateUserOwnsThePlan(id);
		var user = await _userRepo.GetOne(email);
		if (user is null)
		{
			throw new ArgumentException("email_not_found", $"email:{email} was not found");
		}
		await _planMembersRepo.Delete(id, user.Id);
	}

	public async Task Leave(Guid id)
	{
		await ValidateUserLoggedIn();
		int deletedRecords = await _planMembersRepo.Delete(id, Req.UserId);
		if (deletedRecords == 1)
			_log.Info($"user {Req.UserId} left plan {id}");
		else
			throw new ArgumentException($"userId={Req.UserId} unable to leave planId={id}");

	}

	public async Task UpdateType(Guid id, string type)
	{
		await ValidateUserOwnsThePlan(id);
		var count = await _planRepo.GetCount(Req.UserId, type);
		if (count >= 100) throw new LogicException("max_is_100", "Max is 100");

		using var scope = new TransactionScope(TransactionScopeAsyncFlowOption.Enabled);
		await _planRepo.RemoveFromOrder(Req.UserId, id);
		await _planRepo.UpdateType(Req.UserId, id, type);
		scope.Complete();
	}


	public async Task ReOrder(string type, int oldOrder, int newOrder)
	{
		var count = await _planRepo.GetCount(Req.UserId, type);
		if (oldOrder > count || newOrder > count)
			throw new InputException("oldOrder and newOrder should be less than " + count);

		await _planRepo.UpdateOrder(Req.UserId, type, oldOrder, newOrder);
	}

	public async Task ValidateUserOwnsThePlan(Guid planId)
	{
		var plan = await _planRepo.GetOne(planId);
		if (plan is null) throw new ArgumentException("planId not found");
		if (!plan!.User.Id.Equals(Req.UserId))
			throw new UnauthorizedException("User does not own this plan");
	}

	private async Task ValidateUserLoggedIn()
	{
		var user = await _userRepo.GetOne(Req.UserId);
		if (user is null || user.Email == null)
		{
			throw new LogicException("you_are_not_logged_in", "You are not logged In");
		}
	}
}

