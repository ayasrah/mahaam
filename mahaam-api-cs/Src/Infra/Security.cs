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

public class Auth(IDeviceRepo deviceRepo, IUserRepo userRepo, Settings settings) : IAuth
{
	public (Guid, Guid, bool) ValidateAndExtractJwt(HttpContext context)
	{
		string? authorization = context.Request.Headers.Authorization;
		if (string.IsNullOrEmpty(authorization))
		{
			throw new UnauthorizedException($"Authorization header not exists");
		}
		string tokenString = authorization[7..]; // Remove 'Bearer ' to get the jwt

		ValidateJwt(tokenString);
		var token = new JwtSecurityToken(tokenString);


		var userId = GetGuidClaim(token, "userId");
		var deviceId = GetGuidClaim(token, "deviceId");

		var device = deviceRepo.GetOne(deviceId).GetAwaiter().GetResult();
		if ((device is null || userId != device.UserId) && !context.Request.Path.Equals("/user/logout"))
			throw new UnauthorizedException($"Invalid user info");

		var user = userRepo.GetOne(userId).GetAwaiter().GetResult();
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

	public void ValidateJwt(string token)
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

	private TokenValidationParameters GetValidationParams()
	{
		return new TokenValidationParameters()
		{
			ValidateLifetime = true,
			ValidateAudience = false,
			ValidateIssuer = true,
			ValidIssuer = "mahaam-api",
			IssuerSigningKey = SecurityKey()
		};
	}


	private SymmetricSecurityKey SecurityKey()
	{
		var keyBytes = Encoding.ASCII.GetBytes(settings.Api.TokenSecretKey);
		return new SymmetricSecurityKey(keyBytes);
	}
}

