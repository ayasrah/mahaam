using System.Diagnostics;
using System.Text;
using Mahaam.Infra.Monitoring;

namespace Mahaam.Infra;

public class AppMiddleware(RequestDelegate next)
{
	private readonly RequestDelegate _next = next;

	public async Task Invoke(HttpContext context)
	{
		var stopwatch = new Stopwatch();
		stopwatch.Start();
		Req.TrafficId = Guid.NewGuid();

		string? request = null;
		string? response = null;

		try
		{

			if (Config.LogReqEnabled)
			{
				// expensive operation, only used if needed
				context.Request.EnableBuffering();
			}
			AuthenticateReq(context);
			await _next(context);
		}
		catch (Exception e)
		{
			if (Config.LogReqEnabled)
			{
				context.Request.Body.Position = 0;
				request = await GetReqBody(context.Request);
			}
			response = await HandleException(context, e);
		}
		finally
		{
			var path = context.Request.Path.Value ?? "";
			var notTrafficPaths = path.StartsWith("/swagger") || path.Equals("/health") || path.StartsWith("/audit");
			if (!notTrafficPaths)
			{
				CreateTraffic(context, request, response, stopwatch);
			}
		}
	}

	private static void AuthenticateReq(HttpContext context)
	{
		var path = context.Request.Path.Value ?? "";
		var pathBase = context.Request.PathBase.Value;
		if (!"/mahaam-api".Equals(pathBase))
		{
			throw new NotFoundException("Invalid path base");
		}

		string? appStore = context.Request.Headers["x-app-store"];
		string? appVersion = context.Request.Headers["x-app-version"];
		if ((appStore == null || appVersion == null) && !path.StartsWith("/swagger"))
		{
			Log.Error($"Required headers not exists, appStore: {appStore}, appVersion: {appVersion}, path: {path}");
			throw new UnauthorizedException("Required headers not exists");
		}
		Req.AppStore = appStore!;
		Req.AppVersion = appVersion!;

		var bypassAuthPaths = new List<string> {
			"/swagger",
			"/health",
			"/audit",
			"/users/create"
		};

		if (!bypassAuthPaths.Exists(path.StartsWith))
		{
			(Guid userId, Guid deviceId, bool isLoggedIn) = Auth.ValidateAndExtractJwt(context);
			Req.UserId = userId;
			Req.DeviceId = deviceId;
			Req.IsLoggedIn = isLoggedIn;
		}
	}

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

	private static void CreateTraffic(HttpContext context, string? request, string? response, Stopwatch stopwatch)
	{
		var method = context.Request.Method;
		var path = $"{context.Request.Path.Value}{context.Request.QueryString.ToString()}";
		var code = context.Response.StatusCode;
		var elapsed = stopwatch.ElapsedMilliseconds;
		var headers = new TrafficHeaders
		{
			UserId = Req.UserId,
			DeviceId = Req.DeviceId,
			AppStore = Req.AppStore,
			AppVersion = Req.AppVersion,
		};
		var isSuccessResponse = code < 400;
		var isUserPath = path.StartsWith("/user");
		if (string.IsNullOrEmpty(response) || isUserPath)
		{
			response = null;
		}

		if (isSuccessResponse)
		{
			request = null;
			response = null;
		}

		stopwatch.Stop();

		var traffic = new Traffic
		{
			Id = Req.TrafficId,
			Method = method,
			Path = path,
			Code = code,
			Elapsed = elapsed,
			Headers = Json.Serialize(headers),
			Request = request,
			Response = string.IsNullOrEmpty(response) ? null : response,
			HealthId = Cache.HealthId
		};
		App.TrafficRepo.Create(traffic);
	}

	private static async Task<string?> GetReqBody(HttpRequest request)
	{
		using var reader = new StreamReader(request.Body, Encoding.UTF8, true, 1024, true);
		var body = await reader.ReadToEndAsync();
		return string.IsNullOrEmpty(body) ? null : body;
	}
}
