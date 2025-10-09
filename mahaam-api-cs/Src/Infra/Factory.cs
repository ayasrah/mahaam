using Mahaam.Feat.Plans;
using Mahaam.Feat.Tasks;
using Mahaam.Feat.Users;
using Mahaam.Infra.Monitoring;

namespace Mahaam.Infra;

public class Factory
{

	private static readonly Dictionary<Type, object> services = new();

	public static void Put<T>(T instance) => services.Add(typeof(T), instance);

	public static T Get<T>() => (T)services.GetValueOrDefault(typeof(T));


	public static void Init()
	{
		AddRepos();
		AddServices();
	}

	public static void AddServices()
	{
		Put<IPlanService>(new PlanService());
		Put<ITaskService>(new TaskService());
		Put<IUserService>(new UserService());
		Put<IHealthService>(new HealthService());
	}

	public static void AddRepos()
	{
		Put<IPlanRepo>(new PlanRepo());
		Put<IPlanMembersRepo>(new PlanMembersRepo());
		Put<ITaskRepo>(new TaskRepo());
		Put<IUserRepo>(new UserRepo());
		Put<IDeviceRepo>(new DeviceRepo());
		Put<ISuggestedEmailsRepo>(new SuggestedEmailsRepo());
		Put<IHealthRepo>(new HealthRepo());
		Put<ITrafficRepo>(new TrafficRepo());
		Put<ILogRepo>(new LogRepo());
	}

}

public class App
{
	public static IPlanService PlanService => Factory.Get<IPlanService>();
	public static ITaskService TaskService => Factory.Get<ITaskService>();
	public static IUserService UserService => Factory.Get<IUserService>();
	public static IHealthService HealthService => Factory.Get<IHealthService>();

	public static IPlanRepo PlanRepo => Factory.Get<IPlanRepo>();
	public static IPlanMembersRepo PlanMembersRepo => Factory.Get<IPlanMembersRepo>();
	public static ITaskRepo TaskRepo => Factory.Get<ITaskRepo>();
	public static IUserRepo UserRepo => Factory.Get<IUserRepo>();
	public static IDeviceRepo DeviceRepo => Factory.Get<IDeviceRepo>();
	public static ISuggestedEmailsRepo SuggestedEmailsRepo => Factory.Get<ISuggestedEmailsRepo>();
	public static IHealthRepo HealthRepo => Factory.Get<IHealthRepo>();
	public static ITrafficRepo TrafficRepo => Factory.Get<ITrafficRepo>();
	public static ILogRepo LogRepo => Factory.Get<ILogRepo>();
}



