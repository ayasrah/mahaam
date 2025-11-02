# Repositories

### Overview

- Repos are the classes interacting with the DB.
- Repos forms the **foundation** layer that the app is built on.
- Each table has its own repo class.
- Join selects go to the child table repo.
- Methods in repos handle pure db operation, no app logic inside, and start with `Insert, Update, Delete, GetOne, GetMany`.
- Once a method is implemented in this file, feature foundation is ready, just build on top.

### DB libraries

Usually these are DB client libraries options:

- No ORM
- Micro/Lightweight ORM
- Full ORM

Mahaam chose the second option, which is Micro/Lightweight ORM libraries like Dapper. It's very important to allow developers to write raw SQL queries field by field. In my opinion, using a full ORM is like a company outsourcing its core business, no direct access, no full control. With micro ORMs, you still write raw SQL and offload data mapping to the library.

#### Mahaam cares about 3 things when choosing a DB library

- Writing raw SQL.
- Automatically mapping results to objects.
- Supporting named parameters.

#### Chosen Library by Language

- **C#**: Dapper (micro ORM)
- **Java**: jdbi
- **Go**: sqlx
- **Python**: sqlalchemy
- **TypeScript**: postgres

### Mahaam Repos

- Plan Module: `PlanRepo`, `PlanMembersRepo` for `plans`, `plan_members` tables.
- Task Module: `TaskRepo` for `tasks` table.
- User Module: `UserRepo`, `DeviceRepo`, `SuggestedEmailsRepo` for `users`, `devices`, `suggested_emails` tables.
- Monitor Module: `HealthRepo`, `TrafficRepo`, `LogRepo` for `health`, `traffic`, `logs` tables.

### Sample

**TaskRepo** interface

::: code-group

```C#
public interface ITaskRepo
{
    List<Task> GetAll(Guid planId);
    Guid Create(Guid planId, string title);
    void DeleteOne(Guid id);
    void DeleteAll(Guid planId);
    void UpdateDone(Guid id, Boolean done);
    void UpdateTitle(Guid id, string title);
    void UpdateOrder(Guid planId, int oldIndex, int newIndex);
    void UpdateOrderBeforeDelete(Guid planId, Guid id);
    int GetCount(Guid planId);
}
```

```Java
public interface TaskRepo {
    List<Task> getAll(UUID planId);
    UUID create(UUID planId, String title);
    void deleteOne(UUID id);
    void deleteAll(UUID planId);
    void updateDone(UUID id, boolean done);
    void updateTitle(UUID id, String title);
    void updateOrder(UUID planId, int oldIndex, int newIndex);
    void updateOrderBeforeDelete(UUID planId, UUID id);
    long getCount(UUID planId);
}
```

```Go
type TaskRepo interface {
    GetAll(planID UUID) []Task
    Create(tx *sqlx.Tx, planID UUID, title string) UUID
    DeleteOne(tx *sqlx.Tx, id UUID) int64
    UpdateDone(tx *sqlx.Tx, id UUID, done bool) int64
    UpdateTitle(id UUID, title string) int64
    UpdateOrder(tx *sqlx.Tx, planID UUID, oldIndex, newIndex int) int64
    UpdateOrderBeforeDelete(tx *sqlx.Tx, planID UUID, id UUID) int64
    GetCount(planID UUID) int64
}
```

```TypeScript
export interface TasksRepo {
  getAll(trx: Trx, planId: string): Promise<Task[]>;
  create(trx: Trx, planId: string, title: string): Promise<string>;
  deleteOne(trx: Trx, id: string): Promise<void>;
  deleteAll(trx: Trx, planId: string): Promise<void>;
  updateDone(trx: Trx, id: string, done: boolean): Promise<void>;
  updateTitle(trx: Trx, id: string, title: string): Promise<void>;
  updateOrder(trx: Trx, planId: string, oldIndex: number, newIndex: number): Promise<void>;
  updateOrderBeforeDelete(trx: Trx, planId: string, id: string): Promise<void>;
  getCount(trx: Trx, planId: string): Promise<number>;
}
```

```Python
class TaskRepo(Protocol):
    def select_many(self, plan_id: UUID) -> list[Task]: ...
    def create(self, plan_id: UUID, title: str, conn) -> UUID: ...
    def delete_one(self, id: UUID, conn) -> None: ...
    def delete_all(self, plan_id: UUID, conn) -> None: ...
    def update_done(self, id: UUID, done: bool, conn) -> None: ...
    def update_title(self, id: UUID, title: str) -> None: ...
    def update_order(self, plan_id: UUID, old_index: int, new_index: int, conn) -> None: ...
    def update_order_before_delete(self, plan_id: UUID, id: UUID, conn) -> None: ...
    def get_count(self, plan_id: UUID, conn) -> int: ...
```

:::

**TaskRepo** readAll implementation
Example

::: code-group

```C#
public List<Task> GetAll(Guid planId)
{
	var query = @"SELECT id, plan_id, title, done, sort_order, created_at, updated_at
		FROM tasks WHERE plan_id = @planId order by sort_order desc;";
	return DB.SelectMany<Task>(query, new { planId });
}
```

```Java
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
```

```Go
func (r *taskRepo) GetAll(planID UUID) []Task {
	query := `SELECT id, plan_id, title, done, sort_order, created_at, updated_at
		FROM tasks WHERE plan_id = :plan_id ORDER BY sort_order DESC`
	param := Param{"plan_id": planID}
	return dbs.SelectMany[Task](query, param)
}
```

```TypeScript
async getAll(trx: Trx, planId: string): Promise<Task[]> {
	const result = await trx`SELECT id, plan_id, title, done, sort_order,created_at, updated_at
		FROM tasks WHERE plan_id = ${planId} ORDER BY sort_order ASC`;
	return DB.as<Task[]>(result);
}
```

```Python
def select_many(self, plan_id: UUID, conn=None) -> list[Task]:
	sql = """
	SELECT id, plan_id, title, done, sort_order, created_at, updated_at
	FROM tasks WHERE plan_id = :plan_id ORDER BY sort_order DESC"""
	return db.DB.select_many(Task, sql, {"plan_id": str(plan_id)}, conn)
```

:::

#### TypeScript Note

Following code seems to have sql injction problem, as parameter seem to be directly substituted in the query, but that is not the case as the **postgres** library does handle them properly.

```TypeScript
await trx`SELECT * FROM tasks WHERE plan_id = ${planId}`;
```
