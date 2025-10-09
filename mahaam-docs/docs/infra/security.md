# Security

### Overview

Mahaam security ensures that app is used safely, correctly, and only by trusted users.

### Mahaam Security

Mahaam provide anynomous user, by using device as identity. And provide email login with OTP.

Mahaam uses JWT for authentication.

::: code-group

```C#
public static string CreateToken(string userId, string deviceId);
public static (Guid, Guid, bool) ValidateAndExtractJwt(HttpContext context);
```

```Java
public static String createToken(String userId, String deviceId);
public AuthResult validateAndExtractJwt(ContainerRequestContext context);
```

```Go
func (s *AuthService) CreateToken(userId, deviceId uuid.UUID) (string, error);
func (s *AuthService) ValidateAndExtractJwt(r *gin.Context) (uuid.UUID, uuid.UUID, bool);
```

```TypeScript
public static createToken(userId: string, deviceId: string): string;
public static validateAndExtractJwt(request: Request, userRepo: UsersRepo, deviceRepo: DeviceRepo): AuthResult;
```

```Python
def create_token(self, user_id: str, device_id: str) -> str:
def validate_and_extract_jwt(self, request: Request) -> Tuple[uuid.UUID, uuid.UUID, bool]:
```

:::

##### Creating Tokens

- Generates JWT after user creation based on device as identity.
- Generates JWT after user login with email and OTP.

##### Validate And Extract Jwt

- Read JWT from `Authorization` header.
- Validate JWT is issued by Mahaam.
- Extracts JWT claims like userId, deviceId.

##### Auth Middleware

Auth Middleware calls `validateAndExtractJwt` and throws `UnauthorizedException` for invalid tokens or `ForbiddenException` for security violations.

##### Generating JWT secret key (signing key)

Used to generate and verify JWTs.

- Never commit to version control or expose in logs
- Should be rotated periodically for enhanced security
- Save/Read it safe place (vault, env variable) see config section.
- Use different keys for different environments (dev, staging, prod)
- Store keys in secure key management systems in production

**Generation Methods:**

::: code-group

```bash [OpenSSL]
openssl rand -base64 64
```

```bash [Node]
node -e "console.log(require('crypto').randomBytes(64).toString('base64'))"
```

```bash [Python]
python -c "import secrets; print(secrets.token_urlsafe(64))"
```

```bash [Linux]
head -c 64 /dev/urandom | base64
```

:::
