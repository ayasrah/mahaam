using System.Transactions;
using Mahaam.Feat.Users;
using Mahaam.Infra;

namespace Mahaam.Feat.Plans;
public interface IPlanService
{
	Plan GetOne(Guid planId);
	List<Plan> GetMany(string type);
	Guid Create(PlanIn plan);
	void Update(PlanIn plan);
	void Delete(Guid id);
	void Share(Guid id, string email);
	void Unshare(Guid id, string email);
	void Leave(Guid id);
	void UpdateType(Guid id, string type);
	void ReOrder(string type, int oldOrder, int newOrder);
	void ValidateUserOwnsThePlan(Guid planId);
}

public class PlanService(IPlanRepo planRepo, IPlanMembersRepo planMembersRepo, IUserRepo userRepo, ISuggestedEmailsRepo suggestedEmailsRepo, ILog log) : IPlanService
{
	private readonly IPlanRepo _planRepo = planRepo;
	private readonly IPlanMembersRepo _planMembersRepo = planMembersRepo;
	private readonly IUserRepo _userRepo = userRepo;
	private readonly ISuggestedEmailsRepo _suggestedEmailsRepo = suggestedEmailsRepo;
	private readonly ILog _log = log;

	public Plan GetOne(Guid planId)
	{
		var plan = _planRepo.GetOne(planId);
		if (plan is { IsShared: true }) plan.Members = _planMembersRepo.GetUsers(planId);
		return plan;
	}

	public List<Plan> GetMany(string type)
	{
		// plans of the user shared or not
		var plans = _planRepo.GetMany(Req.UserId, type);
		if (Req.IsLoggedIn)
		{
			// plans of others shared with the user
			var sharedPlans = _planMembersRepo.GetOtherPlans(Req.UserId);
			plans.AddRange(sharedPlans);
		}
		return plans;
	}

	public Guid Create(PlanIn plan)
	{
		var userId = Req.UserId;
		var count = _planRepo.GetCount(userId, "Main");
		if (count >= 100) throw new LogicException("max_is_100", "Max is 100");

		return _planRepo.Create(plan);
	}

	public void Update(PlanIn plan)
	{
		ValidateUserOwnsThePlan(plan.Id);
		_planRepo.Update(plan);
	}

	public void Delete(Guid id)
	{
		ValidateUserOwnsThePlan(id);
		using var scope = new TransactionScope();
		_planRepo.RemoveFromOrder(Req.UserId, id);
		_planRepo.Delete(id); // This will delete all related records as it is a cascade delete
		scope.Complete();
	}

	public void Share(Guid id, string email)
	{
		ValidateUserLoggedIn();
		ValidateUserOwnsThePlan(id);
		var user = _userRepo.GetOne(email) ?? throw new ArgumentException("email_not_found", $"email:{email} was not found");

		if (user.Id.Equals(Req.UserId))
			throw new LogicException("not_allowed_to_share_with_creator", "Not allowed to share with creator");

		var limit = 20;
		var plan = _planRepo.GetOne(id);
		if (plan is { IsShared: true })
		{
			var membersCount = _planMembersRepo.GetUsersCount(id);
			if (membersCount >= limit) throw new LogicException("max_is_20", "Max is 20");
		}
		else
		{
			var plansCount = _planMembersRepo.GetPlansCount(Req.UserId);
			if (plansCount >= limit) throw new LogicException("max_is_20", "Max is 20");
		}
		_planMembersRepo.Create(id, user.Id);
		_suggestedEmailsRepo.Create(Req.UserId, email);
		var usr = _userRepo.GetOne(Req.UserId);
		_suggestedEmailsRepo.Create(user.Id, usr.Email!);
	}

	public void Unshare(Guid id, string email)
	{
		ValidateUserLoggedIn();
		ValidateUserOwnsThePlan(id);
		var user = _userRepo.GetOne(email);
		if (user is null)
		{
			throw new ArgumentException("email_not_found", $"email:{email} was not found");
		}
		_planMembersRepo.Delete(id, user.Id);
	}

	public void Leave(Guid id)
	{
		ValidateUserLoggedIn();
		int deletedRecords = _planMembersRepo.Delete(id, Req.UserId);
		if (deletedRecords == 1)
			_log.Info($"user {Req.UserId} left plan {id}");
		else
			throw new ArgumentException($"userId={Req.UserId} unable to leave planId={id}");

	}

	public void UpdateType(Guid id, string type)
	{
		ValidateUserOwnsThePlan(id);
		var count = _planRepo.GetCount(Req.UserId, type);
		if (count >= 100) throw new LogicException("max_is_100", "Max is 100");

		using var scope = new TransactionScope();
		_planRepo.RemoveFromOrder(Req.UserId, id);
		_planRepo.UpdateType(Req.UserId, id, type);
		scope.Complete();
	}


	public void ReOrder(string type, int oldOrder, int newOrder)
	{
		var count = _planRepo.GetCount(Req.UserId, type);
		if (oldOrder > count || newOrder > count)
			throw new InputException("oldOrder and newOrder should be less than " + count);

		_planRepo.UpdateOrder(Req.UserId, type, oldOrder, newOrder);
	}

	public void ValidateUserOwnsThePlan(Guid planId)
	{
		var plan = _planRepo.GetOne(planId);
		if (plan is null) throw new ArgumentException("planId not found");
		if (!plan.User.Id.Equals(Req.UserId))
			throw new UnauthorizedException("User does not own this plan");
	}

	private void ValidateUserLoggedIn()
	{
		var user = _userRepo.GetOne(Req.UserId);
		if (user.Email == null)
		{
			throw new LogicException("you_are_not_logged_in", "You are not logged In");
		}
	}
}

