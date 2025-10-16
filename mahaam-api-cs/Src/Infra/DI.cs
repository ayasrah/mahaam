using Mahaam.Feat.Plans;
using Mahaam.Feat.Tasks;
using Mahaam.Feat.Users;
using Mahaam.Infra.Monitoring;

namespace Mahaam.Infra;

public static class DI
{
	public static void Init(IServiceCollection services)
	{
		services.AddSingleton<IAuth, Auth>();
		services.AddSingleton<ILog, Log>();
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
}
