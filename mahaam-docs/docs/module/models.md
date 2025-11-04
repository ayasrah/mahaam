# Models

### Overview

Data models are related properties collected in a data type.

### Purpose

These are some of the usage of data models:

- Reflect database tables.
- Fill queries result into.
- Fill named params into.
- Grouping function inputs/outputs into.
- Maintain API requests and responses, DTOs (Data Transfer Objects).

### Definition

Following is the `Plan` data model definition:

::: code-group

```C#
public class Plan
{
	public Guid Id { get; set; }
	public string? Title { get; set; }
	public string? Type { get; set; }
	public int SortOrder { get; set; }
	public DateTime? Starts { get; set; }
	public DateTime? Ends { get; set; }
	public string? DonePercent { get; set; }
	public DateTime? CreatedAt { get; set; }
	public DateTime? UpdatedAt { get; set; }
	public List<User>? Members { get; set; }
	public bool IsShared { get; set; }
	public User User { get; set; }
}
```

```Java
public static class Plan {
	public UUID id;
	public String title;
	public String type;
	public int sortOrder;
	public Instant starts;
	public Instant ends;
	public String donePercent;
	public Instant createdAt;
	public Instant updatedAt;
	public List<User> members;
	public boolean isShared;
	public User user;
}
```

```Go
type Plan struct {
	ID          uuid.UUID  `json:"id,omitempty"`
	Title       *string    `json:"title,omitempty"`
	Type        *string    `json:"type,omitempty"`
	SortOrder   int        `json:"sortOrder,omitempty" db:"sort_order"`
	Starts      *time.Time `json:"starts,omitempty"`
	Ends        *time.Time `json:"ends,omitempty"`
	DonePercent *string    `json:"donePercent,omitempty" db:"done_percent"`
	CreatedAt   *time.Time `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
	Members  	[]User     `json:"members,omitempty" db:"shared_with"`
	IsShared    bool       `json:"isShared,omitempty" db:"is_shared"`
	User        User       `json:"user,omitempty" db:"user"`
}
```

```TypeScript
export interface Plan {
  id: string;
  title?: string | null;
  type?: string | null;
  sortOrder: number;
  starts?: Date | null;
  ends?: Date | null;
  donePercent?: string | null;
  createdAt?: Date | null;
  updatedAt?: Date | null;
  members?: User[] | null;
  isShared: boolean;
  user: User;
}
```

```Python
@dataclass
class Plan:
    id: UUID
    title: str | None = None
    type: str | None = None
    sort_order: int = 0
    starts: datetime | None = None
    ends: datetime | None = None
    done_percent: str | None = None
    created_at: datetime | None = None
    updated_at: datetime | None = None
    members: list[User] | None = None
    is_shared: bool = False
    user: User | None = None
```

:::

### See

- Plan model in: [C#](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-cs/Src/Feat/Plan/Plan.cs), [Java](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-java/src/main/java/mahaam/feat/plan/PlanModel.java), [Go](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-go/app/models/plan.go), [TypeScript](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-ts/src/feat/plans/plans.model.ts), [Python](https://github.com/ayasrah/mahaam/blob/main/mahaam-api-py/feat/plan/plan_model.py)
