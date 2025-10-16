using System.IdentityModel.Tokens.Jwt;
using System.Security;
using System.Security.Claims;
using System.Text;
using Mahaam.Feat.Users;
using Microsoft.IdentityModel.Tokens;

namespace Mahaam.Infra;

public interface IAuth
{
	(Guid, Guid, bool) ValidateAndExtractJwt(HttpContext context);
	string CreateToken(string userId, string deviceId);
}

public class Auth(IDeviceRepo deviceRepo, IUserRepo userRepo) : IAuth
{
	private readonly IDeviceRepo _deviceRepo = deviceRepo;
	private readonly IUserRepo _userRepo = userRepo;
	public (Guid, Guid, bool) ValidateAndExtractJwt(HttpContext context)
	{
		string? authorization = context.Request.Headers.Authorization;
		if (string.IsNullOrEmpty(authorization))
		{
			throw new UnauthorizedException($"Authorization header not exists");
		}
		string tokenString = authorization[7..]; // Remove 'Bearer ' to get the jwt

		JWT.Validate(tokenString);
		var token = new JwtSecurityToken(tokenString);


		var userId = GetGuidClaim(token, "userId");
		var deviceId = GetGuidClaim(token, "deviceId");

		var device = _deviceRepo.GetOne(deviceId);
		if ((device is null || userId != device.UserId) && !context.Request.Path.Equals("/user/logout"))
			throw new UnauthorizedException($"Invalid user info");

		var user = _userRepo.GetOne(userId);
		var isLoggedIn = user is not null && user.Email is not null;
		return (userId, deviceId, isLoggedIn);
	}

	private static Guid GetGuidClaim(JwtSecurityToken token, string claimType)
	{
		var claim = token.Claims.First(claim => claim.Type == claimType);
		Rule.Required(claim.Value, claimType);

		if (!Guid.TryParse(claim.Value, out Guid id) || id == Guid.Empty)
		{
			throw new ForbiddenException($"{claimType} is empty");
		}
		return id;
	}

	public string CreateToken(string userId, string deviceId)
	{
		return JWT.Create(userId, deviceId);
	}
}


class JWT
{
	public static string Create(string userId, string deviceId)
	{
		try
		{
			var creds = new SigningCredentials(SecurityKey(), SecurityAlgorithms.HmacSha256);
			var descriptor = new SecurityTokenDescriptor()
			{
				Subject = new ClaimsIdentity([
					new Claim("userId", userId),
					new Claim("deviceId", deviceId),
				]),
				Expires = DateTime.Now.Add(TimeSpan.FromDays(7)),
				SigningCredentials = creds,
				IssuedAt = DateTime.Now,
				Issuer = "mahaam-api",
			};
			var handler = new JwtSecurityTokenHandler();
			var token = handler.CreateToken(descriptor);
			return handler.WriteToken(token);
		}
		catch (Exception exception)
		{
			throw new SecurityException(exception.ToString());
		}
	}

	public static void Validate(string token)
	{
		try
		{
			var tokenHandler = new JwtSecurityTokenHandler();
			var validationParams = GetValidationParams();
			tokenHandler.ValidateToken(token, validationParams, out _);
		}
		catch (SecurityTokenExpiredException)
		{
			throw new UnauthorizedException("Token has expired");
		}
		catch (SecurityTokenException)
		{
			throw new UnauthorizedException("Invalid token");
		}
		catch (Exception)
		{
			throw new UnauthorizedException("Token validation failed");
		}
	}

	private static TokenValidationParameters GetValidationParams()
	{
		return new TokenValidationParameters()
		{
			ValidateLifetime = true, // Because there is no expiration in the generated token
			ValidateAudience = false, // Because there is no audiance in the generated token
			ValidateIssuer = true,   // Because there is no issuer in the generated token
			ValidIssuer = "mahaam-api",
			IssuerSigningKey = SecurityKey()
		};
	}

	private static SymmetricSecurityKey SecurityKey()
	{
		var keyBytes = Encoding.ASCII.GetBytes(Config.TokenSecretKey);
		return new SymmetricSecurityKey(keyBytes);
	}
}