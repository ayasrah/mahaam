package mahaam.feat.plan;

import java.sql.Timestamp;
import java.time.LocalDate;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.UUID;
import java.util.stream.Collectors;

import org.jdbi.v3.core.result.RowView;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import mahaam.feat.plan.PlanModel.Plan;
import mahaam.feat.plan.PlanModel.PlanIn;
import mahaam.feat.plan.PlanModel.PlanType;
import mahaam.feat.task.Task;
import mahaam.feat.user.UserModel.User;
import mahaam.infra.DB;
import mahaam.infra.Log;
import mahaam.infra.Mapper;
import mahaam.infra.Req;

public interface PlanRepo {
	Plan getOne(UUID id);

	List<Plan> getMany(UUID userId, String type);

	UUID create(PlanIn plan);

	void update(PlanIn plan);

	void delete(UUID id);

	void updateDonePercent(UUID id);

	void removeFromOrder(UUID userId, UUID id);

	void updateOrder(UUID userId, String type, int oldOrder, int newOrder);

	void updateType(UUID userId, UUID id, String type);

	long getCount(UUID userId, String type);

	int updateUserId(UUID oldUserId, UUID newUserId);
}

@ApplicationScoped
class DefaultPlanRepo implements PlanRepo {

	@Inject
	DB db;

	@Override
	public UUID create(PlanIn plan) {
		String query = """
				INSERT INTO plans (id, user_id, title, starts, ends, type, status, done_percent, sort_order, created_at)
				VALUES (:id, :userId, :title, :starts, :ends, :type, :status, '0/0',
				(SELECT COUNT(1) FROM plans WHERE user_id = :userId AND type = :type), current_timestamp) """;

		UUID id = UUID.randomUUID();
		db.insert(
				query,
				Mapper.of(
						"id",
						id,
						"title",
						plan.title,
						"starts",
						toTimestamp(plan.starts),
						"ends",
						toTimestamp(plan.ends),
						"userId",
						Req.getUserId(),
						"type",
						PlanType.MAIN,
						"status",
						"Open"));

		return id;
	}

	@Override
	public void update(PlanIn plan) {
		String query = """
				UPDATE plans
				SET title = :title, starts = :starts, ends = :ends, updated_at = current_timestamp
				WHERE id = :id
				""";

		db.update(
				query,
				Mapper.of(
						"id", plan.id,
						"title", plan.title,
						"starts", toTimestamp(plan.starts),
						"ends", toTimestamp(plan.ends)));
	}

	@Override
	public Plan getOne(UUID id) {
		String query = """
				SELECT c.id c_id, c.title c_title, c.starts c_starts, c.ends c_ends, c.type c_type,
				    c.done_percent c_donePercent, c.sort_order c_sortOrder, c.user_id c_userId,
					EXISTS(SELECT 1 FROM plan_members cm WHERE cm.plan_id = c.id) AS c_isShared,
					u.id u_id, u.email u_email, u.name u_name
				FROM plans c
				LEFT JOIN users u ON c.user_id = u.id
				WHERE c.id = :id""";

		var plan = db.getJdbi().withHandle(
				handle -> {
					return handle
							.createQuery(query)
							.bind("id", id)
							.reduceRows(
									new LinkedHashMap<Plan, User>(),
									(Map<Plan, User> map, RowView rowView) -> {
										Plan c = rowView.getRow(Plan.class);
										User u = rowView.getRow(User.class);
										c.user = u;
										map.put(c, u);
										return map;
									})
							.entrySet()
							.stream()
							.map(entry -> entry.getKey())
							.findFirst()
							.orElse(null);
				});

		return plan;
	}

	@Override
	public List<Plan> getMany(UUID userId, String type) {
		try {
			String query = """
					SELECT c.id c_id, c.title c_title, c.starts c_starts, c.ends c_ends, c.type c_type,
					    c.done_percent c_donePercent, c.sort_order c_sortOrder, c.user_id c_userId,
					    EXISTS(SELECT 1 FROM plan_members cm WHERE cm.plan_id = c.id) AS c_isShared,
					    u.id u_id, u.email u_email, u.name u_name
					FROM plans c
					LEFT JOIN users u ON c.user_id = u.id
					WHERE c.user_id = :user_id AND c.type = :type
					ORDER BY c.sort_order DESC;""";

			var plans = db.getJdbi()
					.withHandle(
							handle -> {
								return handle
										.createQuery(query)
										.bind("user_id", userId)
										.bind("type", type)
										.reduceRows(
												new LinkedHashMap<Plan, User>(),
												(Map<Plan, User> map, RowView rowView) -> {
													Plan c = rowView.getRow(Plan.class);
													User u = rowView.getRow(User.class);
													c.user = u;
													map.put(c, u);
													return map;
												})
										.keySet()
										.stream()
										.collect(Collectors.toList());
							});
			return plans;
		} catch (Exception e) {
			Log.error(e.toString());
			return new ArrayList<>();
		}
	}

	@Override
	public void delete(UUID id) {
		int count = db.delete("DELETE FROM plans WHERE id = :id", Mapper.of("id", id));
		if (count > 0) {
			Log.info("Plan " + id + " deleted");
		}
	}

	@Override
	public void updateDonePercent(UUID id) {
		String query = """
				SELECT id t_id, plan_id t_planId, title t_title, done t_done, sort_order t_sortOrder,
					created_at t_createdAt, updated_at t_updatedAt FROM tasks t WHERE t.plan_id = :id""";
		List<Task> tasks = db.selectList(query, Task.class, Mapper.of("id", id));

		long done = tasks.stream().filter(task -> task.done).count();
		int notDone = tasks.size();
		String donePercent = done + "/" + notDone;
		String updateQuery = "UPDATE plans SET done_percent = :donePercent WHERE id = :id";
		db.update(updateQuery, Mapper.of("donePercent", donePercent, "id", id));
	}

	@Override
	public void removeFromOrder(UUID userId, UUID id) {
		String query = """
				UPDATE plans
				SET sort_order = sort_order - 1
				WHERE
				    user_id = :userId AND
				    type = (SELECT type FROM plans WHERE id = :id) AND
				    sort_order > (SELECT sort_order FROM plans WHERE id = :id)
				""";

		db.update(query, Mapper.of("userId", userId, "id", id));
	}

	@Override
	public void updateOrder(UUID userId, String type, int oldOrder, int newOrder) {
		String query = """
				UPDATE plans SET sort_order =
				    CASE
				        WHEN sort_order = :oldOrder THEN :newOrder
				        WHEN sort_order > :oldOrder AND sort_order <= :newOrder THEN sort_order - 1
				        WHEN sort_order >= :newOrder AND sort_order < :oldOrder THEN sort_order + 1
				        ELSE sort_order
				    END
				WHERE
				    user_id = :userId AND
				    type = :type
				""";

		db.update(
				query,
				Mapper.of(
						"userId", userId,
						"type", type,
						"oldOrder", oldOrder,
						"newOrder", newOrder));
	}

	@Override
	public void updateType(UUID userId, UUID id, String type) {
		String query = """
				UPDATE plans
				SET type = :type,
				    sort_order = (SELECT COUNT(1) FROM plans WHERE user_id = :userId AND type = :type),
				    updated_at = current_timestamp
				WHERE id = :id
				""";

		db.update(query, Mapper.of("userId", userId, "id", id, "type", type));
	}

	@Override
	public long getCount(UUID userId, String type) {
		String queryCount = "SELECT COUNT(*) FROM plans WHERE user_id = :userId AND type = :type";
		return db.selectCount(queryCount, Mapper.of("userId", userId, "type", type));
	}

	@Override
	public int updateUserId(UUID oldUserId, UUID newUserId) {
		String query = """
				UPDATE plans
				SET user_id = :newUserId,
				    sort_order = (sort_order + (SELECT count(1) FROM plans WHERE user_id = :newUserId)),
				    updated_at = current_timestamp
				WHERE user_id = :oldUserId
				""";

		return db.update(query, Mapper.of("oldUserId", oldUserId, "newUserId", newUserId));
	}

	private Timestamp toTimestamp(LocalDate dateTime) {
		if (dateTime == null) {
			return null;
		}
		return java.sql.Timestamp.valueOf(dateTime.atStartOfDay());
	}
}
