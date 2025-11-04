# Transaction Management

#### Overview

Transaction management is treating a set of db operations as a single unit, **all succeed or none**.

#### Purpose

To ensure data consistency.

#### Management Steps

- If all operation succeed, commit is done.
- If any operation fails, rollback is done.
- Mahaam manages transactions in the service layer, since all db operations are initiated there.
- In `C# and Java`, the framework is doning commit and rollback (ambient transaction) via `TransactionScope` and `@Transactional`.
- In Go, Javascript, and Python, commit and rollback are done manually: `transaction.begin`,`transaction.commit` or `transaction.rollback`.
- Log and traffic repos are execluded from transactions (`suppressed`), so audits are created to db even the transaction is rolled back for that request. Its recommeded to place the log at the end of the method, eg: after `scope.Complete();` to make sure transaction is completed.

#### Mahaam Code

Mahaam implements transaction management in the service layer across all language versions:

::: code-group

```C#
// C# uses TransactionScope for declarative transaction management
public Guid CreateTask(Guid planId, string title)
{
	using var scope = new TransactionScope();
	var id = _taskRepo.Create(planId, title);
	_planRepo.UpdateDonePercent(planId);
	scope.Complete();
	return id;
}
```

```Java
// Java uses @Transactional annotation for method-level transaction management
@Override
@Transactional
public UUID createTask(UUID planId, String title) {
	UUID id = taskRepo.create(planId, title);
	planRepo.updateDonePercent(planId);
	return id;
}
```

```Go
// Go uses a higher-order function that wraps transaction logic
// infra/dbs/db.go - WithTx function
func WithTx(fn func(tx *sqlx.Tx) error) error {
    tx, err := appDB.Beginx()
    if err != nil {
        return err
    }
    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
            panic(p)
        } else if err != nil {
            tx.Rollback()
        } else {
            err = tx.Commit()
        }
    }()
    return fn(tx)
}

// Service layer usage
func (s *taskService) CreateTask(planID UUID, title string) UUID {
	var id UUID
	err := dbs.WithTx(func(tx *sqlx.Tx) error {
		id = s.taskRepo.Create(tx, planID, title)
		s.planRepo.UpdateDonePercent(tx, planID)
		return nil
	})
	if err != nil {
		panic(models.LogicErr(err.Error(), "error_creating_task"))
	}
	return id
}
```

```TypeScript
// TypeScript/Node.js uses transaction callbacks with connection passing
async createTask(planId: string, title: string): Promise<string> {
	return await DB.withTrx(async (trx) => {
		const count = await this.tasksRepo.getCount(trx, planId);
		if (count >= 100) throw new LogicError('max_is_100', 'Max is 100');

		const id = await this.tasksRepo.create(trx, planId, title);
		await this.plansRepo.updateDonePercent(trx, planId);

		return id;
	});
}
```

```Python
# Python uses a context manager that automatically handles transaction lifecycle
# infra/db.py - Transaction scope implementation
@staticmethod
@contextmanager
def transaction_scope():
    conn = DB.get_engine().connect()
    try:
        conn.begin()
        yield conn
        conn.commit()
    except Exception:
        conn.rollback()
        raise

# Service layer usage
def create(self, plan_id: UUID, title: str) -> UUID:
	with db.DB.transaction_scope() as conn:
		count = self.task_repo.get_count(plan_id, conn)
		if count >= 100:
			raise LogicException("max_is_100", "Max is 100")
		id = self.task_repo.create(plan_id, title, conn)
		self.plan_repo.update_done_percent(plan_id, conn)
	return id
```

:::

### See

- Transaction management in TaskService in: [C#](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-cs/Src/Feat/Task/TaskService.cs), [Java](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-java/src/main/java/mahaam/feat/task/TaskService.java), [Go](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-go/app/service/task.go), [TypeScript](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-ts/src/feat/tasks/tasks.service.ts), [Python](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-py/feat/task/task_service.py)
