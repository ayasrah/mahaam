# Validation

### Overview

Validation ensures data integrity and business rules.

### Example

- **Required**: Validates that a field is not null, empty, or whitespace.
- **OneAtLeastRequired**: Ensures at least one field from a list is provided.
- **FailIf**: Throws exception when a condition is true.
- **In/Contains**: Validates a value exists in a list.

::: code-group

```C#
Rule.Required(userId, "userId");
Rule.ValidateEmail(email);
Rule.In(type, PlanType.All);
```

```Java
Rule.required(userId, "userId");
Rule.validateEmail(email);
Rule.in(type, PlanType.All);
```

```TypeScript
rule.required(userId, "userId");
rule.validateEmail(email);
rule.isIn(type, PlanType.All);
```

```Python
Rule.required(user_id, "user_id")
Rule.validate_email(email)
Rule.contains(PlanType.All, type)
```

:::

## Validation Levels

- Controllers: Validate required inputs and types.
- Service: Validate business rules.
- Middlewares: Validate security and authorization rules.
