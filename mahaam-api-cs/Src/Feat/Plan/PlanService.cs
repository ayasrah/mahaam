using System.Transactions;
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

public class PlanService : IPlanService
{

	public Plan GetOne(Guid planId)
	{
		var plan = App.PlanRepo.GetOne(planId);
		if (plan is { IsShared: true }) plan.Members = App.PlanMembersRepo.GetUsers(planId);
		return plan;
	}

	public List<Plan> GetMany(string type)
	{
		// plans of the user shared or not
		var plans = App.PlanRepo.GetMany(Req.UserId, type);
		if (Req.IsLoggedIn)
		{
			// plans of others shared with the user
			var sharedPlans = App.PlanMembersRepo.GetOtherPlans(Req.UserId);
			plans.AddRange(sharedPlans);
		}
		return plans;
	}

	public Guid Create(PlanIn plan)
	{
		var userId = Req.UserId;
		var count = App.PlanRepo.GetCount(userId, "Main");
		if (count >= 100) throw new LogicException("max_is_100", "Max is 100");

		return App.PlanRepo.Create(plan);
	}

	public void Update(PlanIn plan)
	{
		ValidateUserOwnsThePlan(plan.Id);
		App.PlanRepo.Update(plan);
	}

	public void Delete(Guid id)
	{
		ValidateUserOwnsThePlan(id);
		using var scope = new TransactionScope();
		App.PlanRepo.RemoveFromOrder(Req.UserId, id);
		App.PlanRepo.Delete(id); // This will delete all related records as it is a cascade delete
		scope.Complete();
	}

	public void Share(Guid id, string email)
	{
		ValidateUserLoggedIn();
		ValidateUserOwnsThePlan(id);
		var user = App.UserRepo.GetOne(email) ?? throw new ArgumentException("email_not_found", $"email:{email} was not found");

		if (user.Id.Equals(Req.UserId))
			throw new LogicException("not_allowed_to_share_with_creator", "Not allowed to share with creator");

		var limit = 20;
		var plan = App.PlanRepo.GetOne(id);
		if (plan is { IsShared: true })
		{
			var membersCount = App.PlanMembersRepo.GetUsersCount(id);
			if (membersCount >= limit) throw new LogicException("max_is_20", "Max is 20");
		}
		else
		{
			var plansCount = App.PlanMembersRepo.GetPlansCount(Req.UserId);
			if (plansCount >= limit) throw new LogicException("max_is_20", "Max is 20");
		}
		App.PlanMembersRepo.Create(id, user.Id);
		App.SuggestedEmailsRepo.Create(Req.UserId, email);
		var usr = App.UserRepo.GetOne(Req.UserId);
		App.SuggestedEmailsRepo.Create(user.Id, usr.Email!);
	}

	public void Unshare(Guid id, string email)
	{
		ValidateUserLoggedIn();
		ValidateUserOwnsThePlan(id);
		var user = App.UserRepo.GetOne(email);
		if (user is null)
		{
			throw new ArgumentException("email_not_found", $"email:{email} was not found");
		}
		App.PlanMembersRepo.Delete(id, user.Id);
	}

	public void Leave(Guid id)
	{
		ValidateUserLoggedIn();
		int deletedRecords = App.PlanMembersRepo.Delete(id, Req.UserId);
		if (deletedRecords == 1)
			Log.Info($"user {Req.UserId} left plan {id}");
		else
			throw new ArgumentException($"userId={Req.UserId} unable to leave planId={id}");

	}

	public void UpdateType(Guid id, string type)
	{
		ValidateUserOwnsThePlan(id);
		var count = App.PlanRepo.GetCount(Req.UserId, type);
		if (count >= 100) throw new LogicException("max_is_100", "Max is 100");

		using var scope = new TransactionScope();
		App.PlanRepo.RemoveFromOrder(Req.UserId, id);
		App.PlanRepo.UpdateType(Req.UserId, id, type);
		scope.Complete();
	}


	public void ReOrder(string type, int oldOrder, int newOrder)
	{
		var count = App.PlanRepo.GetCount(Req.UserId, type);
		if (oldOrder > count || newOrder > count)
			throw new InputException("oldOrder and newOrder should be less than " + count);

		App.PlanRepo.UpdateOrder(Req.UserId, type, oldOrder, newOrder);
	}

	public void ValidateUserOwnsThePlan(Guid planId)
	{
		var plan = App.PlanRepo.GetOne(planId);
		if (plan is null) throw new ArgumentException("planId not found");
		if (!plan.User.Id.Equals(Req.UserId))
			throw new UnauthorizedException("User does not own this plan");
	}

	private static void ValidateUserLoggedIn()
	{
		var user = App.UserRepo.GetOne(Req.UserId);
		if (user.Email == null)
		{
			throw new LogicException("you_are_not_logged_in", "You are not logged In");
		}
	}
}

