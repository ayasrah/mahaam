package mahaam.feat.plan;

import java.util.List;
import java.util.UUID;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import jakarta.transaction.Transactional;
import mahaam.feat.plan.PlanModel.Plan;
import mahaam.feat.plan.PlanModel.PlanIn;
import mahaam.feat.user.SuggestedEmailsRepo;
import mahaam.feat.user.UserModel.User;
import mahaam.feat.user.UserRepo;
import mahaam.infra.Exceptions.InputException;
import mahaam.infra.Exceptions.LogicException;
import mahaam.infra.Exceptions.UnauthorizedException;
import mahaam.infra.Log;
import mahaam.infra.Req;

public interface PlanService {
	Plan getOne(UUID planId);

	List<Plan> getMany(String type);

	UUID create(PlanIn plan);

	void update(PlanIn plan);

	void delete(UUID id);

	void share(UUID id, String email);

	void unshare(UUID id, String email);

	void leave(UUID id);

	void updateType(UUID id, String type);

	void reOrder(String type, int oldOrder, int newOrder);

	void validateUserOwnsThePlan(UUID planId);
}

@ApplicationScoped
class DefaultPlanService implements PlanService {

	@Inject
	PlanRepo planRepo;

	@Inject
	PlanMembersRepo planMembersRepo;

	@Inject
	UserRepo userRepo;

	@Inject
	SuggestedEmailsRepo suggestedEmailsRepo;

	@Override
	public Plan getOne(UUID planId) {
		Plan plan = planRepo.getOne(planId);
		if (plan != null && plan.isShared) {
			List<User> members = planMembersRepo.getUsers(planId);
			// Since Plan is now a class (mutable), we can modify fields directly
			plan.members = members;
		}
		return plan;
	}

	@Override
	public List<Plan> getMany(String type) {
		List<Plan> plans = planRepo.getMany(Req.getUserId(), type); // plans of the user shared
																	// or not
		if (Req.isLoggedIn()) {
			// plans of others shared with the user
			List<Plan> sharedPlans = planMembersRepo.getOtherPlans(Req.getUserId());
			plans.addAll(sharedPlans);
		}
		return plans;
	}

	@Override
	public UUID create(PlanIn plan) {
		var count = planRepo.getCount(Req.getUserId(), "Main");
		if (count >= PlanConstants.MAX_GROUPS_PER_USER) {
			throw new LogicException("max_is_100", "Max is 100");
		}

		return planRepo.create(plan);
	}

	@Override
	public void update(PlanIn plan) {
		validateUserOwnsThePlan(plan.id);
		planRepo.update(plan);
	}

	@Override
	@Transactional
	public void delete(UUID id) {
		validateUserOwnsThePlan(id);

		planRepo.removeFromOrder(Req.getUserId(), id);
		planRepo.delete(id); // This will delete all related records as it is a cascade delete
	}

	@Override
	public void share(UUID id, String email) {
		validateUserLoggedIn();
		validateUserOwnsThePlan(id);
		User user = userRepo.getOne(email);
		if (user == null) {
			throw new InputException("email_not_found: email:" + email + " was not found");
		}

		if (user.id.equals(Req.getUserId())) {
			throw new LogicException("not_allowed_to_share_with_creator", "Not allowed to share with creator");
		}

		Plan plan = planRepo.getOne(id);
		if (plan != null && plan.isShared) {
			var membersCount = planMembersRepo.getUsersCount(id);
			if (membersCount >= PlanConstants.MAX_SHARED_USERS) {
				throw new LogicException("max_is_20", "Max is 20");
			}
		} else {
			var plansCount = planMembersRepo.getPlansCount(Req.getUserId());
			if (plansCount >= PlanConstants.MAX_SHARED_GROUPS) {
				throw new LogicException("max_is_20", "Max is 20");
			}
		}

		planMembersRepo.create(id, user.id);
		suggestedEmailsRepo.create(Req.getUserId(), email);
		User usr = userRepo.getOne(Req.getUserId());
		suggestedEmailsRepo.create(user.id, usr.email);
	}

	@Override
	public void unshare(UUID id, String email) {
		validateUserLoggedIn();
		validateUserOwnsThePlan(id);
		User user = userRepo.getOne(email);
		if (user == null) {
			throw new InputException("email_not_found: email:" + email + " was not found");
		}

		planMembersRepo.delete(id, user.id);
	}

	@Override
	public void leave(UUID id) {
		validateUserLoggedIn();
		int deletedRecords = planMembersRepo.delete(id, Req.getUserId());
		if (deletedRecords == 1) {
			Log.info("user " + Req.getUserId() + " left plan " + id);
		} else {
			throw new InputException("userId=" + Req.getUserId() + " cannot leave planId=" + id);
		}
	}

	@Override
	@Transactional
	public void updateType(UUID id, String type) {
		validateUserOwnsThePlan(id);
		var count = planRepo.getCount(Req.getUserId(), type);
		if (count >= 100) {
			throw new LogicException("max_is_100", "Max is 100");
		}

		planRepo.removeFromOrder(Req.getUserId(), id);
		planRepo.updateType(Req.getUserId(), id, type);
	}

	@Override
	public void reOrder(String type, int oldOrder, int newOrder) {
		var count = planRepo.getCount(Req.getUserId(), type);
		if (oldOrder > count || newOrder > count) {
			throw new InputException("oldOrder and newOrder should be less than " + count);
		}
		planRepo.updateOrder(Req.getUserId(), type, oldOrder, newOrder);
	}

	@Override
	public void validateUserOwnsThePlan(UUID planId) {
		Plan plan = planRepo.getOne(planId);
		if (plan == null) {
			throw new InputException("planId not found");
		}
		if (!plan.user.id.equals(Req.getUserId())) {
			throw new UnauthorizedException("User does not own this plan");
		}
	}

	private void validateUserLoggedIn() {
		if (!Req.isLoggedIn()) {
			throw new LogicException("you_are_not_logged_in", "You are not logged In");
		}
	}
}

class PlanConstants {
	public static final int MAX_GROUPS_PER_USER = 100;
	public static final int MAX_SHARED_USERS = 20;
	public static final int MAX_SHARED_GROUPS = 20;
}