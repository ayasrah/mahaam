package mahaam.feat.task;

import java.util.List;
import java.util.UUID;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import jakarta.transaction.Transactional;
import mahaam.feat.plan.PlanRepo;
import mahaam.infra.Exceptions.InputException;
import mahaam.infra.Exceptions.LogicException;

public interface TaskService {
	UUID create(UUID planId, String title);

	List<Task> getList(UUID planId);

	void delete(UUID planId, UUID id);

	void updateDone(UUID planId, UUID id, boolean done);

	void updateTitle(UUID id, String title);

	void reOrder(UUID planId, int oldOrder, int newOrder);
}

@ApplicationScoped
class DefaultTaskService implements TaskService {

	@Inject
	TaskRepo taskRepo;

	@Inject
	PlanRepo planRepo;

	@Override
	@Transactional
	public UUID create(UUID planId, String title) {
		var count = taskRepo.getCount(planId);
		if (count >= 100) {
			throw new LogicException("max_is_100", "Max is 100");
		}

		UUID id = taskRepo.create(planId, title);
		planRepo.updateDonePercent(planId);
		return id;
	}

	@Override
	public List<Task> getList(UUID planId) {
		return taskRepo.getAll(planId);
	}

	@Override
	@Transactional
	public void delete(UUID planId, UUID id) {
		taskRepo.updateOrderBeforeDelete(planId, id);
		taskRepo.deleteOne(id);
		planRepo.updateDonePercent(planId);
	}

	@Override
	@Transactional
	public void updateDone(UUID planId, UUID id, boolean done) {
		taskRepo.updateDone(id, done);
		planRepo.updateDonePercent(planId);

		var count = taskRepo.getCount(planId);
		var task = taskRepo.getOne(id);
		reOrder(planId, task.sortOrder, done ? 0 : (int) count - 1);
	}

	@Override
	public void updateTitle(UUID id, String title) {
		taskRepo.updateTitle(id, title);
	}

	@Override
	public void reOrder(UUID planId, int oldOrder, int newOrder) {
		if (oldOrder == newOrder) {
			return;
		}
		var count = taskRepo.getCount(planId);
		if (oldOrder > count || newOrder > count) {
			throw new InputException("oldOrder and newOrder should be less than " + count);
		}
		taskRepo.updateOrder(planId, oldOrder, newOrder);
	}
}
