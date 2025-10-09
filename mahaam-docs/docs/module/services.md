# Services

### Overview

Services are the business logic (domain) layer.

### Job

- Service layer contains the business logic classes.
- Transaction management is done in service layer.
- Service layer orchestrates business operations and enforce business rules.
- Service layer are called from controllers.
- Service layer calls repositories for data access.
- Services layer calls external sevices like twilio.

### Implementation

Mahaam defines **interface** and **implementation** for each service in one file, for readability and decoupling.

### Mahaam Services

- PlanService
- TaskService
- UserService
- HealthService

### Example

**PlanService** interface

::: code-group

```C#
public interface IPlanService
{
    Plan GetOne(Guid planId);
    List<Plan> GetMany(string? type);
    Guid Create(PlanIn plan);
    void Update(PlanIn plan);
    void Delete(Guid id);
    void Share(Guid id, string email);
    void Unshare(Guid id, string email);
    void Leave(Guid id);
    void UpdateType(Guid id, string type);
    void ReOrder(string type, int oldIndex, int newIndex);
    void ValidateUserOwnsThePlan(Guid planId);
}
```

```Java
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
    void reOrder(String type, int oldIndex, int newIndex);
    void validateUserOwnsThePlan(UUID planId);
}
```

```Go
type PlanService interface {
    GetOne(planID UUID) *Plan
    GetMany(userID UUID, planType string) []Plan
    Create(userID UUID, plan PlanIn) UUID
    Update(userID UUID, plan *PlanIn)
    Delete(userID UUID, id UUID)
    Share(userID UUID, id UUID, email string)
    Unshare(userID UUID, id UUID, email string)
    Leave(userID UUID, id UUID)
    UpdateType(userID UUID, id UUID, planType string)
    ReOrder(userID UUID, planType string, oldIndex, newIndex int)
    ValidateUserOwnsThePlan(userID UUID, planID UUID)
}
```

```TypeScript
export interface PlansService {
  getOne(id: string): Promise<Plan | null>;
  getMany(type?: string): Promise<Plan[]>;
  create(plan: PlanIn): Promise<string>;
  update(plan: PlanIn): Promise<void>;
  delete(id: string): Promise<void>;
  share(id: string, email: string): Promise<void>;
  unshare(id: string, email: string): Promise<void>;
  leave(id: string): Promise<void>;
  updateType(id: string, type: string): Promise<void>;
  reOrder(type: string, oldIndex: number, newIndex: number): Promise<void>;
  validateUserOwnsThePlan(planId: string): Promise<void>;
}
```

```Python
class PlanService(Protocol):
    def get_one(self, plan_id: UUID) -> Plan: ...
    def get_many(self, type: str | None) -> list[Plan]: ...
    def create(self, plan: PlanIn) -> UUID: ...
    def update(self, plan: PlanIn) -> None: ...
    def delete(self, id: UUID) -> None: ...
    def share(self, id: UUID, email: str) -> None: ...
    def unshare(self, id: UUID, email: str) -> None: ...
    def leave(self, id: UUID) -> None: ...
    def update_type(self, id: UUID, type: str) -> None: ...
    def reorder(self, type: str, old_index: int, new_index: int) -> None: ...
    def validate_user_owns_the_plan(self, plan_id: UUID) -> None: ...
```

:::

**PlanService** Implementation

::: code-group

```C#
public class PlanService : IPlanService
{
    public Guid Create(PlanIn plan)
    {
        var userId = Request<Guid>.Get("userId");
        var count = App.PlanRepo.GetCount(userId, "Main");
        if (count >= 100) throw new LogicException("max_is_100", "Max is 100");

        using var scope = new TransactionScope();
        var id = App.PlanRepo.Create(plan);
        scope.Complete();
        return id;
    }
}
```

```Java
@ApplicationScoped
class DefaultPlanService implements PlanService {
    @Inject
    PlanRepo planRepo;

    @Override
    @Transactional
    public UUID create(PlanIn plan) {
        UUID userId = ReqData.get("userId");
        var count = planRepo.getCount(userId, "Main");
        if (count >= PlanConstants.MAX_GROUPS_PER_USER) {
            throw new LogicException("max_is_100", "Max is 100");
        }

        return planRepo.create(plan);
    }
}
```

```Go
type planService struct {
    planRepo           repo.PlanRepo
}

func NewPlanService(db *sqlx.DB, planRepo repo.PlanRepo) PlanService {
    return &planService{
        planRepo:  planRepo
    }
}

func (s *planService) Create(userID UUID, plan PlanIn) UUID {
    count := s.planRepo.GetCount(userID, string(models.PlanTypeMain))
    if count >= 100 {
        panic(models.HttpErr{Code: 409, Message: "maximum of 100 plans reached"})
    }

    var planID UUID
    err := dbs.WithTx(func(tx *sqlx.Tx) error {
        planID = s.planRepo.Create(tx, userID, plan)
        return nil
    })

    if err != nil {
        return uuid.Nil
    }
    return planID
}
```

```TypeScript
@Injectable()
export class DefaultPlansService implements PlansService {
  constructor(@Inject("PlansRepo") private readonly plansRepo: PlansRepo) {}

  async create(plan: PlanIn): Promise<string> {
    const userId = ReqCxt.get<string>("userId")!;
    const id = await DB.withTrx(async (trx) => {
      const count = await this.plansRepo.getCount(trx, userId, "Main");
      if (count >= 100) throw new LogicError("max_is_100", "Max is 100");

      return await this.plansRepo.create(trx, plan);
    });
    return id;
  }
}
```

```Python
class DefaultPlanService(metaclass=ProtocolEnforcer, protocol=PlanService):
    def __init__(self, plan_repo: PlanRepo) -> None:
        self.plan_repo = plan_repo

    def create(self, plan: PlanIn) -> UUID:
        user_id = Req.get("userId", UUID)
        count = self.plan_repo.get_count(user_id, "Main")
        if count >= 100:
            raise LogicException("max_is_100", "Max is 100")

        with db.DB.transaction_scope() as conn:
            id = self.plan_repo.create(plan, conn)
        return id
```

:::

### Key Patterns Across Languages

1. **Interface Definition**: All languages define a service interface with the same methods
2. **Dependency Injection**: Services receive repositories as dependencies
3. **Transaction Management**: Each language handles transactions appropriately:
   - C#: `TransactionScope`
   - Java: `@Transactional` annotation
   - Go: `dbs.WithTx` function
   - Python: `db.DB.transaction_scope()` context manager
   - TypeScript: `DB.withTrx` async function
4. **Business Logic**: Services contain validation and business rules
5. **Error Handling**: Consistent error handling patterns across languages
