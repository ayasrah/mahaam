using System.Net;
using System.Reflection;
using System.Threading.RateLimiting;
using Mahaam.Infra;
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
	.AddJsonFile("config.json", optional: false, reloadOnChange: true);
var config = configBuilder.Build();
Config.Init(config);


builder.WebHost.UseKestrel(opts =>
{
	opts.Listen(IPAddress.Parse("0.0.0.0"), Config.HttpPort);
});

Serilog.Log.Logger = new LoggerConfiguration()
	.WriteTo.Console(outputTemplate: "{Timestamp:yyyy-MM-dd HH:mm:ss.fff} [{Level:u3}] {Message:lj}{NewLine}{Exception}")
	.WriteTo.File(Config.LogFile!,
		rollingInterval: RollingInterval.Infinite,
		rollOnFileSizeLimit: true,
		fileSizeLimitBytes: Config.LogFileSizeLimit,
		retainedFileCountLimit: Config.LogFileCountLimit,
		outputTemplate: Config.LogFileOutputTemplate
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

// Configure the HTTP request pipeline.
var app = builder.Build();
app.UsePathBase(new PathString("/mahaam-api"));

app.UseMiddleware<AppMiddleware>();

app.UseCors("AllowAll");
if ("local".Equals(Config.EnvName))
{
	app.UseSwagger();
	app.UseSwaggerUI();
}
app.UseRouting();
app.UseRateLimiter();
app.MapControllers();

Factory.Init();
Starter.Init(app);

// for dapper
Dapper.DefaultTypeMap.MatchNamesWithUnderscores = true;

app.Run();
// executed when app is stopped
App.HealthService.ServerStopped();
