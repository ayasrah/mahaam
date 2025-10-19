# Exceptions

### Overview

**Exceptions** are the errors occur during app execution.

**Exception Handling** is how the app manages these errors and the flow control.

### Mahaam Exceptions

Mahaam defined 5 type of exceptions:

- **`UnauthorizedException`**: Identity violation - `401 http status`
- **`ForbiddenException`**: Role violation - `403 http status`
- **`InputException`**: Invalid input - `400 http status`
- **`NotFoundException`**: Resource not found - `404 http status`
- **`LogicException`**: Business logic violation - `409 http status`

### Example

::: code-group

```C#
// In TaskService.cs
public Guid Create(Guid planId, string title)
{
	var count = _taskRepo.GetCount(planId);
	if (count >= 100) throw new LogicException("max_is_100", "Max is 100");

	//... rest of create logic
}
```

```Java
// In TaskService.java
@Override
@Transactional
public UUID create(UUID planId, String title) {
	var count = taskRepo.getCount(planId);
	if (count >= 100) {
		throw new LogicException("max_is_100", "Max is 100");
	}
	//... rest of create logic
}
```

```Go
// In task_service.go
func (s *taskService) Create(planID UUID, title string) UUID {
	count := s.taskRepo.GetCount(planID)
	if count >= 100 {
		panic(models.LogicErr("maximum number of tasks is 100", "max_is_100"))
	}
	//... rest of create logic
}
```

```TypeScript
// In tasks.service.ts
async create(planId: string, title: string): Promise<string> {
	return await DB.withTrx(async (trx) => {
		const count = await this.tasksRepo.getCount(trx, planId);
		if (count >= 100) throw new LogicError('max_is_100', 'Max is 100');

		//... rest of create logic
	});
}
```

```Python
# In task_service.py
def create(self, plan_id: UUID, title: str) -> UUID:
	with db.DB.transaction_scope() as conn:
		count = self.task_repo.get_count(plan_id, conn)
		if count >= 100:
			raise LogicException("max_is_100", "Max is 100")
		#... rest of create logic
```

:::

### Handling

Exception handling in software is generally scoped per library. This means:

- When a request calls a library (e.g., lib X), the library attempts to handle the logic internally.
- If an exception occurs within the library, the internal flow breaks and control is passed up to the top-level interface of the library.
- The top-level interface is responsible for returning a meaningful response to the caller.

Exceptions should be captured and logged, either at the source or, at the top level.

In **Mahaam**, any layer (repo, service, or controller) can throw an exception. Mahaam does not log or handle exceptions at their origin, instead when exception happens, it breaks the request flow, till reach the top level **Exception Handler Middleware** where it is logged and handled.

### Purpose

This approach enhances code readability and maintainability by:

- Centralizing logging and response formatting.
- Making the codebase cleaner and easier to debug.

### Example

::: code-group

```C#
	[HttpPatch]
	[Route("{id}/unshare")]
	public IActionResult Unshare(Guid id, [FromForm] string email)
	{
		Rule.Required(id, "id");
		Rule.Required(email, "email");
		_planService.Unshare(id, email);
		return StatusCode(Http.Ok);
	}
```

```Java
@PatchMapping("/{id}/unshare")
public ResponseEntity<Void> unshare(@PathVariable UUID id, @RequestParam String email) {
    Rule.required(id, "id");
    Rule.required(email, "email");
    planService.unshare(id, email);
    return ResponseEntity.ok().build();
}
```

```Go
func (h *PlanHandler) Unshare(c *gin.Context) {
	id := c.Param("id")
	email := c.PostForm("email")

	Rule.Required(id, "id")
	Rule.Required(email, "email")
	h.planService.Unshare(id, email)
	c.Status(http.StatusOK)
}
```

```TypeScript
@Patch(':id/unshare')
async unshare(@Param('id') id: string, @Body('email') email: string) {
    Rule.required(id, 'id');
    Rule.required(email, 'email');
    this.planService.unshare(id, email);
    return { status: HttpStatus.OK };
}
```

```Python
@router.patch("/{id}/unshare")
def unshare(id: str, email: str = Form(...)):
    Rule.required(id, "id")
    Rule.required(email, "email")
    plan_service.unshare(id, email)
    return {"status": HTTPStatus.OK}
```

:::

- `Rule.Required` can throw an exception if the input is invalid.
- `PlanService.Unshare` may also throw exceptions due to business logic or repository errors.

### Exception Handler Middleware

Any exception raised in the **Repo**, **Service**, or **Controller** layers will propagate upward and be caught by the **Exception Handler Middleware**, which in role:

- Detect exception type
- Log exception details
- Return a proper response

::: code-group

```C#
private static async Task<string> HandleException(HttpContext context, Exception e)
{
	var response = Json.Serialize(e.Message);
	var code = Http.ServerError;


	if (e is AppException)
	{
		var appException = e as AppException;
		var key = appException!.Key;
		code = appException.HttpCode;
		if (!string.IsNullOrEmpty(key))
		{
			var res = new { key, error = e.Message };
			response = Json.Serialize(res);
		}
	}

	Log.Error(e.ToString());
	context.Response.StatusCode = code;
	context.Response.ContentType = Http.json;
	await context.Response.WriteAsync(response);
	return response;
}
```

```Java
@Provider
class ExceptionFilter implements ExceptionMapper<Throwable> {

	@Override
	public Response toResponse(Throwable exception) {

		String response = Json.toString(exception.getMessage());
		int statusCode = Http.ServerError;
		exception.printStackTrace();
		Log.error(exception.toString());
		if (exception instanceof AppException) {
			AppException appException = (AppException) exception;
			statusCode = appException.getHttpCode();
			String key = appException.getKey();
			if (key != null && !key.isEmpty()) {
				ErrorResponse errorResponse = new ErrorResponse(key, exception.getMessage());
				response = Json.toString(errorResponse);
			}
		}

		return Response.status(statusCode).entity(response).type(Http.JsonMedia).build();
	}
}
```

```Go
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic
				trafficID, ok := c.Value("trafficID").(uuid.UUID)
				if !ok || trafficID == uuid.Nil {
					trafficID = uuid.Nil
				}

				if e, ok := err.(*models.HttpErr); ok {
					logs.Error(trafficID, e.Error())
					if e.Key == "" {
						c.JSON(e.Code, e.Message)
					} else {
						c.JSON(e.Code, gin.H{"error": e.Message, "key": e.Key})
					}
					c.Abort()
					return
				} else {
					logs.Error(trafficID, err)
				}

				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				c.Abort()
			}
		}()
		c.Next()
	}
}
```

```TypeScript
return next.handle().pipe(
	catchError((err: unknown) => {
	Log.error(err instanceof Error ? err.toString() : String(err));

	const ctx = context.switchToHttp();
	response = ctx.getResponse<Response>();

	responseBody = JSON.stringify(err instanceof Error ? err.message : String(err));
	let code = HttpStatus.INTERNAL_SERVER_ERROR;

	if (err instanceof AppError) {
		const appException = err as AppError;
		const key = appException.key;
		code = appException.httpCode;
		if (key) {
		const responseObj = { key, error: err.message };
		responseBody = JSON.stringify(responseObj);
		}
	}
	response.status(code).contentType('application/json');
	response.send(responseBody);

	return throwError(() => err);
	}),
);
```

```Python
def handle_exception( e: Exception, traffic_id: uuid.UUID) -> tuple[str, int]:
	response_status = http.SERVER_ERROR
	res_body = json.dumps(str(e))

	if isinstance(e, AppException):
		response_status = e.http_code
		if e.key:
			res_body = json.dumps(
				{"key": e.key, "error": str(e)})

	log.Log.error(str(e), traffic_id=traffic_id)
	return res_body, response_status
```

:::

### Go case

Go traditionally avoids exceptions in favor of returning errors. However, `mahaam-api-go` follows the same exception handling pattern used in other Mahaam backends like `mahaam-api-cs`, `mahaam-api-java`, etc.

From the official [Go blog on defer, panic, and recover](https://go.dev/blog/defer-panic-and-recover):

> For a real-world example of panic and recover, see the json package from the Go standard library. It encodes an interface with a set of recursive functions. If an error occurs when traversing the value, panic is called to unwind the stack to the top-level function call, which recovers from the panic and returns an appropriate error value...

This principle is applied similarly in Mahaam’s Go codebase, panics are used internally and recovered at the top level to transform them into structured error responses to the caller.
