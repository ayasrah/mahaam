package mahaam.infra;

import java.nio.charset.StandardCharsets;
import java.security.Key;
import java.util.Date;
import java.util.UUID;

import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.SignatureAlgorithm;
import io.jsonwebtoken.security.Keys;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import jakarta.ws.rs.container.ContainerRequestContext;
import mahaam.feat.user.DeviceRepo;
import mahaam.feat.user.UserModel.User;
import mahaam.feat.user.UserRepo;
import mahaam.infra.Exceptions.ForbiddenException;
import mahaam.infra.Exceptions.UnauthorizedException;

@ApplicationScoped
public class Security {

	@Inject
	DeviceRepo deviceRepo;

	@Inject
	UserRepo userRepo;

	@Inject
	Config config;

	public static class AuthResult {
		private UUID userId;
		private UUID deviceId;
		private boolean isLoggedIn;

		public AuthResult(UUID userId, UUID deviceId, boolean isLoggedIn) {
			this.userId = userId;
			this.deviceId = deviceId;
			this.isLoggedIn = isLoggedIn;
		}

		public UUID getUserId() {
			return userId;
		}

		public UUID getDeviceId() {
			return deviceId;
		}

		public boolean isLoggedIn() {
			return isLoggedIn;
		}
	}

	private static final String EMPTY_UUID = "00000000-0000-0000-0000-000000000000";

	public static void nonEmptyUuid(String uuidString, String name) {
		Rule.required(uuidString, name);
		if (uuidString == null || uuidString.trim().isEmpty() || EMPTY_UUID.equals(uuidString)) {
			throw new ForbiddenException(name + " is Empty");
		}
	}

	public AuthResult validateAndExtractJwt(ContainerRequestContext context)
			throws UnauthorizedException, ForbiddenException {
		String path = context.getUriInfo().getPath();
		String authorization = context.getHeaderString("Authorization");

		if (authorization == null || authorization.isEmpty()) {
			throw new UnauthorizedException("Authorization header not exists");
		}

		if (!authorization.startsWith("Bearer ")) {
			throw new UnauthorizedException("Invalid Authorization header format");
		}

		String tokenString = authorization.substring(7); // Remove 'Bearer ' to get the jwt

		try {
			validate(tokenString);
			Claims claims = Jwts.parserBuilder()
					.setSigningKey(getSecurityKey())
					.build()
					.parseClaimsJws(tokenString)
					.getBody();

			String userIdClaim = claims.get("userId", String.class);
			nonEmptyUuid(userIdClaim, "userId");

			String deviceIdClaim = claims.get("deviceId", String.class);
			nonEmptyUuid(deviceIdClaim, "deviceId");

			UUID userId = UUID.fromString(userIdClaim);
			UUID deviceId = UUID.fromString(deviceIdClaim);

			// Use the device repository to get device info
			var device = deviceRepo.getOne(deviceId);
			if ((device == null || !userId.equals(device.userId)) && !"/user/logout".equals(path)) {
				throw new UnauthorizedException("Invalid user info");
			}

			User user = userRepo.getOne(userId);
			boolean isLoggedIn = user.email != null;

			return new AuthResult(userId, deviceId, isLoggedIn);

		} catch (Exception e) {
			if (e instanceof UnauthorizedException || e instanceof ForbiddenException) {
				throw e;
			}
			throw new UnauthorizedException("Invalid JWT token: " + e.getMessage());
		}
	}

	public String createToken(String userId, String deviceId) {
		try {
			Key key = getSecurityKey();
			Date now = new Date();
			Date expiration = new Date(now.getTime() + (7 * 24 * 60 * 60 * 1000)); // 7 days

			return Jwts.builder()
					.claim("userId", userId)
					.claim("deviceId", deviceId)
					.setIssuedAt(now)
					.setExpiration(expiration)
					.setIssuer("mahaam-api")
					.signWith(key, SignatureAlgorithm.HS256)
					.compact();

		} catch (Exception e) {
			throw new RuntimeException("Error creating JWT token: " + e.getMessage(), e);
		}
	}

	private void validate(String token) {
		try {
			Jwts.parserBuilder()
					.setSigningKey(getSecurityKey())
					.requireIssuer("mahaam-api")
					.build()
					.parseClaimsJws(token);
		} catch (Exception e) {
			throw new RuntimeException("JWT validation failed: " + e.getMessage(), e);
		}
	}

	private Key getSecurityKey() {
		byte[] keyBytes = config.tokenSecretKey().getBytes(StandardCharsets.UTF_8);
		return Keys.hmacShaKeyFor(keyBytes);
	}

}
