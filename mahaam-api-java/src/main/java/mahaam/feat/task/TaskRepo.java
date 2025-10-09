package mahaam.feat.task;

import java.util.List;
import java.util.UUID;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import mahaam.infra.DB;
import mahaam.infra.Exceptions.InputException;
import mahaam.infra.Mapper;

public interface TaskRepo {
	List<Task> getAll(UUID planId);

	Task getOne(UUID id);

	UUID create(UUID planId, String title);

	void deleteOne(UUID id);

	void deleteAll(UUID planId);

	void updateDone(UUID id, boolean done);

	void updateTitle(UUID id, String title);

	void updateOrder(UUID planId, int oldOrder, int newOrder);

	void updateOrderBeforeDelete(UUID planId, UUID id);

	long getCount(UUID planId);
}

@ApplicationScoped
class DefaultTaskRepo implements TaskRepo {

	@Inject
	DB db;

	@Override
	public List<Task> getAll(UUID planId) {
		String query = """
				SELECT
					 t.id t_id,
					 t.plan_id t_planId,
					 t.title t_title,
					 t.done t_done,
					 t.sort_order t_sortOrder,
					 t.created_at t_createdAt,
					 t.updated_at t_updatedAt
				FROM tasks t
				WHERE t.plan_id = :planId
				ORDER BY t.sort_order DESC
				""";
		return db.selectList(query, Task.class, Mapper.of("planId", planId));
	}

	@Override
	public Task getOne(UUID id) {
		String query = """
				SELECT
					t.id t_id,
					t.plan_id t_planId,
					t.title t_title,
					t.done t_done,
					t.sort_order t_sortOrder,
					t.created_at t_createdAt,
					t.updated_at t_updatedAt
				FROM tasks t WHERE t.id = :id""";

		return db.selectOne(query, Task.class, Mapper.of("id", id));
	}

	@Override
	public UUID create(UUID planId, String title) {
		String query = """
				INSERT INTO tasks (id, plan_id, title, done, sort_order, created_at)
				VALUES (:id, :planId, :title, :done, (SELECT COUNT(1) FROM tasks WHERE plan_id = :planId), current_timestamp)
				""";

		UUID id = UUID.randomUUID();
		db.insert(
				query,
				Mapper.of(
						"id", id,
						"planId", planId,
						"title", title,
						"done", false));
		return id;
	}

	@Override
	public void deleteOne(UUID id) {
		String query = "DELETE FROM tasks WHERE id = :id";
		int deletedRows = db.delete(query, Mapper.of("id", id));
		if (deletedRows == 0) {
			throw new InputException("task id=" + id + " not found");
		}
	}

	@Override
	public void deleteAll(UUID planId) {
		String query = "DELETE FROM tasks WHERE plan_id = :planId";
		db.delete(query, Mapper.of("planId", planId));
	}

	@Override
	public void updateDone(UUID id, boolean done) {
		String query = "UPDATE tasks SET done = :done, updated_at = current_timestamp WHERE id = :id";
		int updatedRows = db.update(query, Mapper.of("id", id, "done", done));
		if (updatedRows == 0) {
			throw new InputException("task id=" + id + " not found");
		}
	}

	@Override
	public void updateTitle(UUID id, String title) {
		String query = "UPDATE tasks SET title = :title, updated_at = current_timestamp WHERE id = :id";
		int updatedRows = db.update(query, Mapper.of("id", id, "title", title));
		if (updatedRows == 0) {
			throw new InputException("task id=" + id + " not found");
		}
	}

	@Override
	public void updateOrder(UUID planId, int oldOrder, int newOrder) {
		String query = """
				UPDATE tasks SET sort_order =
				    CASE
				        WHEN sort_order = :oldOrder THEN :newOrder
				        WHEN sort_order > :oldOrder AND sort_order <= :newOrder THEN sort_order - 1
				        WHEN sort_order >= :newOrder AND sort_order < :oldOrder THEN sort_order + 1
				        ELSE sort_order
				    END
				WHERE plan_id = :planId
				""";

		db.update(
				query,
				Mapper.of(
						"planId", planId,
						"oldOrder", oldOrder,
						"newOrder", newOrder));
	}

	@Override
	public void updateOrderBeforeDelete(UUID planId, UUID id) {
		String query = """
				UPDATE tasks
				SET sort_order = sort_order - 1
				WHERE
				    plan_id = :planId AND
				    sort_order > (SELECT sort_order FROM tasks WHERE id = :id)
				""";
		db.update(query, Mapper.of("planId", planId, "id", id));
	}

	@Override
	public long getCount(UUID planId) {
		String queryCount = "SELECT COUNT(1) FROM tasks WHERE plan_id = :planId";
		return db.selectCount(queryCount, Mapper.of("planId", planId));
	}
}
