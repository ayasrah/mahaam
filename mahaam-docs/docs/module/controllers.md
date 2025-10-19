# Controllers

#### Overview

Controllers are the API definition and routing layer.

#### Job

- Define APIs.
- Routing Requests.
- Methods should not have logic.
- Input validation, like required and data types.
- Delegate business logic to services.

#### Implementation

Mahaam defines **interface** and **implementation** for each controller in one file, for readability.

#### Mahaam Controllers

- PlanController
- TaskController
- UserController
- AuditController

#### Sample

**IPlanController** interface

::: code-group

```C#
public interface IPlanController
{
    IActionResult Create(PlanIn plan);
    IActionResult Update(PlanIn plan);
    IActionResult Delete(Guid id);
    IActionResult Share(Guid id, string email);
    IActionResult Unshare(Guid id, string email);
    IActionResult Leave(Guid id);
    IActionResult UpdateType(Guid id, string type);
    IActionResult ReOrder(string type, int oldIndex, int newIndex);
    IActionResult GetOne(Guid planId);
    IActionResult GetMany(string? type);
}
```

```Java
public interface PlanController {
    Response create(PlanIn plan);
    Response update(PlanIn plan);
    Response delete(UUID id);
    Response share(UUID id, String email);
    Response unshare(UUID id, String email);
    Response leave(UUID id);
    Response updateType(UUID id, String type);
    Response reOrder(String type, int oldIndex, int newIndex);
    Response getOne(UUID planId);
    Response getMany(String type);
}
```

```Go
type PlanHandler interface {
    Create(c *ginC)
    Update(c *ginC)
    Delete(c *ginC)
    Share(c *ginC)
    Unshare(c *ginC)
    Leave(c *ginC)
    UpdateType(c *ginC)
    ReOrder(c *ginC)
    GetOne(c *ginC)
    GetMany(c *ginC)
}
```

```TypeScript
export interface PlansController {
  create(plan: PlanIn, res: Response): Promise<void>;
  update(plan: PlanIn, res: Response): Promise<void>;
  delete(id: string, res: Response): Promise<void>;
  share(id: string, email: string, res: Response): Promise<void>;
  unshare(id: string, email: string, res: Response): Promise<void>;
  leave(id: string, res: Response): Promise<void>;
  updateType(id: string, type: string, res: Response): Promise<void>;
  reOrder(type: string, oldIndex: number, newIndex: number, res: Response): Promise<void>;
  getOne(planId: string, res: Response): Promise<void>;
  getMany(res: Response, type?: string): Promise<void>;
}
```

```Python
class PlanRouter(Protocol):
    def create(self, plan: PlanIn = Body(...)) -> Response: ...
    def update(self, plan: PlanIn = Body(...)) -> Response: ...
    def delete(self, id: UUID = Path(...)) -> Response: ...
    def share(self, id: UUID = Path(...), email: str = Form(...)) -> Response: ...
    def unshare(self, id: UUID = Path(...), email: str = Form(...)) -> Response: ...
    def leave(self, id: UUID = Path(...)) -> Response: ...
    def update_type(self, id: UUID = Path(...), type: str = Form(...)) -> Response: ...
    def reorder(self, type: str = Form(...), old_index: int = Form(...), new_index: int = Form(...)) -> Response: ...
    def get_one(self, plan_id: UUID = Path(...)) -> Response: ...
    def get_many(self, type: str | None = Query(None)) -> Response: ...
```

:::

**PlanController** Implementation

::: code-group

```C#
[ApiController]
[Route("plans")]
public class PlanController : ControllerBase, IPlanController
{
    [HttpGet]
    [Route("{planId}")]
    public IActionResult GetOne(Guid planId)
    {
        Rule.Required(planId, "planId");
        var plan = _planService.GetOne(planId);
        return StatusCode(Http.Ok, plan);
    }
}
```

```Java
@ApplicationScoped
@Path("/plans")
@Consumes(Http.JsonMedia)
@Produces(Http.JsonMedia)
class DefaultPlanController implements PlanController {

    @GET
    @Path("/{planId}")
    public Response getOne(@PathParam("planId") UUID planId) {
        Rule.required(planId, "planId");
        Plan plan = planService.getOne(planId);
        return Response.status(Http.OK).entity(Json.toString(plan)).build();
    }
}
```

```Go
type planHandler struct {
    planService service.PlanService
}

func NewPlanHandler(service service.PlanService) PlanHandler {
    return &planHandler{planService: service}
}

func (h *planHandler) GetOne(c *ginC) {
    id := req.PathUuid(c, "planId")
    plan := h.planService.GetOne(id)
    c.JSON(OK, plan)
}
```

```TypeScript
@Controller("plans")
export class PlansController {
  constructor(@Inject("PlansService") private readonly plansService: PlansService) {}

  @Get(":planId")
  async getOne(@Param("planId") planId: string, @Res() res: Response) {
    rule.required(planId, "planId");
    const plan = await this.plansService.getOne(planId);
    res.status(200).json(plan);
  }
}
```

```Python
@cbv(router)
class DefaultPlanRouter(metaclass=ProtocolEnforcer, protocol=PlanRouter):
    def __init__(self, plan_service: PlanService = Depends(get_plan_service)):
        self.plan_service = plan_service

    @router.get("/plans/{plan_id}", response_model=Plan)
    def get_one(self, plan_id: UUID = Path(...)) -> Response:
        Rule.required(plan_id, "planId")
        plan = self.plan_service.get_one(plan_id)
        return plan
}
```

:::

#### Http Methods

- GET: get resource
- POST: Create resource
- PUT: Replace resource
- PATCH: Update resource partially
- DELETE: Delete resource

#### Http status codes

Mahaam uses status codes as follows:

- 200: OK, **success** GET, PATCH, PUT
- 201: Created, **success** POST
- 204: NoContent, **success** DELETE
- 400: BadRequest, **failed** in input, GET, POST, PATCH, PUT, DELETE
- 401: Unauthorized, **failed** in identity, GET, POST, PATCH, PUT, DELETE
- 403: Forbidden, **failed** in role, GET, POST, PATCH, PUT, DELETE
- 404: NotFound, **failed** in resource, GET, POST, PATCH, PUT, DELETE
- 409: Conflict, **failed** in logic, GET, POST, PATCH, PUT, DELETE
- 500: ServerError, **failed** in server, GET, POST, PATCH, PUT, DELETE
