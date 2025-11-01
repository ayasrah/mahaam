using System.Net;
using System.Reflection;
using System.Threading.RateLimiting;
using Mahaam.Feat.Plans;
using Mahaam.Feat.Tasks;
using Mahaam.Feat.Users;
using Mahaam.Infra;
using Mahaam.Infra.Monitoring;
using Microsoft.Extensions.Options;
using Serilog;

var builder = WebApplication.CreateBuilder(args);
var services = builder.Services;
services.AddControllers();
services.AddLogging(a => a.SetMinimumLevel(LogLevel.Warning));
services.AddEndpointsApiExplorer();
services.AddSwaggerGen();
services.AddMvc().AddJsonOptions(options =>
{
	options.JsonSerializerOptions.DefaultIgnoreCondition = System.Text.Json.Serialization.JsonIgnoreCondition.WhenWritingNull;
});

var configBuilder = new ConfigurationBuilder()
	.SetBasePath(Path.GetDirectoryName(Assembly.GetEntryAssembly().Location))
	.AddJsonFile("appsettings.json", optional: false, reloadOnChange: true);
var config = configBuilder.Build();

services.Configure<Settings>(config.GetSection("app"));
services.AddSingleton(sp => sp.GetRequiredService<IOptions<Settings>>().Value);

var settings = services.BuildServiceProvider().GetRequiredService<Settings>();

builder.WebHost.UseKestrel(opts =>
{
	opts.Listen(IPAddress.Parse("0.0.0.0"), settings.Api.HttpPort);
});


Serilog.Log.Logger = new LoggerConfiguration()
	.WriteTo.Console(outputTemplate: "{Timestamp:yyyy-MM-dd HH:mm:ss.fff} [{Level:u3}] {Message:lj}{NewLine}{Exception}")
	.WriteTo.File(settings.Logging.File!,
		rollingInterval: RollingInterval.Infinite,
		rollOnFileSizeLimit: true,
		fileSizeLimitBytes: settings.Logging.FileSizeLimit,
		retainedFileCountLimit: settings.Logging.FileCountLimit,
		outputTemplate: settings.Logging.OutputTemplate
	)
	.CreateLogger();

services.AddCors(options =>
{
	options.AddPolicy("AllowAll", builder =>
	{
		builder.AllowAnyOrigin()
			.AllowAnyMethod()
			.AllowAnyHeader();
	});
});

services.AddRateLimiter(options =>
{
	options.AddPolicy("PerUserRateLimit", context =>
	{
		// Use User.Identity.Name or a custom identifier to apply per-user rate limits
		string userId = context.Connection.RemoteIpAddress.ToString();

		// Define a fixed window rate limiter per user
		return RateLimitPartition.GetTokenBucketLimiter(userId, _ => new TokenBucketRateLimiterOptions
		{
			TokenLimit = 5, // Max 5 requests
			TokensPerPeriod = 5,
			ReplenishmentPeriod = TimeSpan.FromMinutes(1),
			AutoReplenishment = true
		});
	});
});

_addServices(services);

// Configure the HTTP request pipeline.
var app = builder.Build();
app.UsePathBase(new PathString("/mahaam-api"));
app.UseMiddleware<AppMiddleware>();
app.UseCors("AllowAll");
if ("local".Equals(settings.Api.EnvName))
{
	app.UseSwagger();
	app.UseSwaggerUI();
}
app.UseRouting();
app.UseRateLimiter();
app.MapControllers();

Starter.Init(app);

// for dapper
Dapper.DefaultTypeMap.MatchNamesWithUnderscores = true;

app.Run();
app.Services.GetService<IHealthService>()?.ServerStopped();


static void _addServices(IServiceCollection services)
{
	services.AddSingleton<ICache, Cache>();
	services.AddSingleton<IAuth, Auth>();
	services.AddSingleton<ILog, Mahaam.Infra.Log>();
	services.AddSingleton<IEmail, Email>();
	services.AddSingleton<IDB, DB>();
	services.AddSingleton<IPlanRepo, PlanRepo>();
	services.AddSingleton<IPlanMembersRepo, PlanMembersRepo>();
	services.AddSingleton<ITaskRepo, TaskRepo>();
	services.AddSingleton<IUserRepo, UserRepo>();
	services.AddSingleton<IDeviceRepo, DeviceRepo>();
	services.AddSingleton<ISuggestedEmailsRepo, SuggestedEmailsRepo>();
	services.AddSingleton<IHealthRepo, HealthRepo>();
	services.AddSingleton<ITrafficRepo, TrafficRepo>();
	services.AddSingleton<ILogRepo, LogRepo>();
	services.AddSingleton<IPlanService, PlanService>();
	services.AddSingleton<ITaskService, TaskService>();
	services.AddSingleton<IUserService, UserService>();
	services.AddSingleton<IHealthService, HealthService>();
}