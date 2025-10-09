package mahaam.feat.plan;

import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.UUID;
import java.util.stream.Collectors;

import org.jdbi.v3.core.result.RowView;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import mahaam.feat.plan.PlanModel.Plan;
import mahaam.feat.user.UserModel.User;
import mahaam.infra.DB;
import mahaam.infra.Mapper;

public interface PlanMembersRepo {
	void create(UUID planId, UUID userId);

	int delete(UUID planId, UUID userId);

	List<Plan> getOtherPlans(UUID userId);

	List<User> getUsers(UUID planId);

	long getPlansCount(UUID userId);

	long getUsersCount(UUID planId);
}

@ApplicationScoped
class DefaultPlanMembersRepo implements PlanMembersRepo {

	@Inject
	DB db;

	@Override
	public void create(UUID planId, UUID userId) {
		String query = """
				INSERT INTO plan_members( plan_id, user_id, created_at)
				VALUES(:planId, :userId, current_timestamp)""";

		db.insert(query, Mapper.of("planId", planId, "userId", userId));
	}

	@Override
	public int delete(UUID planId, UUID userId) {
		String query = "DELETE FROM plan_members WHERE plan_id = :planId AND user_id = :userId";
		return db.delete(query, Mapper.of("planId", planId, "userId", userId));
	}

	@Override
	public List<Plan> getOtherPlans(UUID userId) {
		String query = """
				SELECT c.id c_id, c.title c_title, c.starts c_starts, c.ends c_ends, c.type c_type,
					c.done_percent c_donePercent, c.sort_order c_sortOrder,
					c.created_at c_createdAt, u.id u_id, u.email u_email, u.name u_name
				FROM plan_members cm
				LEFT JOIN plans c ON cm.plan_id = c.id
				LEFT JOIN users u ON c.user_id = u.id
				WHERE cm.user_id = :userId
				ORDER BY c.created_at DESC""";
		var plans = db.getJdbi()
				.withHandle(
						handle -> {
							return handle
									.createQuery(query)
									.bind("userId", userId)
									.reduceRows(
											new LinkedHashMap<Plan, User>(),
											(Map<Plan, User> map, RowView rowView) -> {
												Plan g = rowView.getRow(Plan.class);
												User u = rowView.getRow(User.class);
												g.user = u;
												g.isShared = true;
												map.put(g, u);
												return map;
											})
									.keySet()
									.stream()
									.collect(Collectors.toList());
						});
		return plans;
	}

	@Override
	public List<User> getUsers(UUID planId) {
		String query = """
				SELECT u.id u_id, u.email u_email, u.name u_name
				FROM plan_members cm
				LEFT JOIN users u ON cm.user_id = u.id
				WHERE cm.plan_id = :planId
				ORDER BY u.created_at DESC""";
		return db.selectList(query, User.class, Mapper.of("planId", planId));
	}

	@Override
	public long getPlansCount(UUID userId) {
		String query = """
				SELECT COUNT(1) FROM plan_members cm
				LEFT JOIN plans c ON cm.plan_id = c.id
				WHERE c.user_id = :userId""";
		return db.selectCount(query, Mapper.of("userId", userId));
	}

	@Override
	public long getUsersCount(UUID planId) {
		String query = "SELECT COUNT(1) FROM plan_members WHERE plan_id = :planId";
		return db.selectCount(query, Mapper.of("planId", planId));
	}
}
