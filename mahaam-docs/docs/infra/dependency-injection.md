# Dependency Injection

### Overview

Dependency Injection is giving classes the instances that depends on rather than letting these classes creating their dependencies inside.

Dependnecies are usually interfaces (contracts) not classes.

### Importance

- Decouple components.
- More testable, you can inject mocks.

### Example

- `PlanController` depends on `IPlanService` interface.
- An implementation of the `IPlanService` interface is injected, given, provided to PlanController.

::: code-group

```C#
public class PlanController(IPlanService planService) : ControllerBase, IPlanController
{
	// planService instance is injected to PlanController
	// private readonly IPlanService _planService = new PlanService(); // This is dependency concrete creation (tightly coupled)
}

// Which instance should be passed? configured in Program.cs:
services.AddSingleton<IPlanService, PlanService>();
```

```Java
@ApplicationScoped
class DefaultPlanController implements PlanController {

	@Inject
	PlanService planService; // PlanService is the refernce type (the interface), planService is the instance (the object) which it's value is of type DefaultPlanService which is an implementation to that interface.
	// PlanService planService = new DefaultPlanService(); // This is dependency concrete creation (tightly coupled)
}
```

```Go
type planHandler struct {
	planService service.PlanService
	logger      logs.Logger
}

func NewPlanHandler(service service.PlanService, logger logs.Logger) PlanHandler {
	return &planHandler{planService: service, logger: logger}
}
// an instance of type PlanService is injected to PlanHandler in main using NewPlanHandler
```

```TypeScript
@Controller('plans')
export class DefaultPlansController implements PlansController {
  constructor(@Inject('PlansService') private readonly plansService: PlansService) {}
}

// Definition in factory.ts

@Module({
  imports: [],
  controllers: [
    DefaultPlansController,
  ],
  providers: [
    { provide: 'PlansService', useClass: DefaultPlansService },
  ],
})
export class AppModule implements NestModule {
  configure(consumer: MiddlewareConsumer) {
    consumer.apply(AppMiddleware).exclude('/health', '/swagger', '/api-docs').forRoutes('*path');
  }
}
```

```Python
def get_plan_service() -> PlanService:
    from infra.factory import App
    return App.plan_service

@cbv(router)
class DefaultPlanRouter(metaclass=ProtocolEnforcer, protocol=PlanRouter):
    def __init__(self, plan_service: PlanService = Depends(get_plan_service)):
        self.plan_service = plan_service

# Definition in factory.py
from feat.plan.plan_service import DefaultPlanService, PlanService

class App():
    plan_repo: PlanRepo = DefaultPlanRepo()
    plan_members_repo: PlanMembersRepo = DefaultPlanMembersRepo()
    user_repo: UserRepo = DefaultUserRepo()
    plan_service: PlanService = DefaultPlanService(plan_repo=plan_repo, plan_members_repo=plan_members_repo, user_repo=user_repo,

```

:::

### See

- Dependency injection in PlanController in: [C#](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-cs/Src/Feat/Plan/PlanController.cs), [Java](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-java/src/main/java/mahaam/feat/plan/PlanController.java), [Go](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-go/app/handler/plan.go), [TypeScript](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-ts/src/feat/plans/plans.controller.ts), [Python](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-py/feat/plan/plan_router.py)
