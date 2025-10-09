using System.Collections.Concurrent;

namespace Mahaam.Infra;

public static class Req
{
	public static Guid TrafficId
	{
		get => ReqContext<Guid>.Get("trafficId")!;
		set => ReqContext<Guid>.Set("trafficId", value);
	}

	public static Guid UserId
	{
		get => ReqContext<Guid>.Get("userId");
		set => ReqContext<Guid>.Set("userId", value);
	}

	public static Guid DeviceId
	{
		get => ReqContext<Guid>.Get("deviceId");
		set => ReqContext<Guid>.Set("deviceId", value);
	}

	public static string AppStore
	{
		get => ReqContext<string>.Get("appStore")!;
		set => ReqContext<string>.Set("appStore", value);
	}

	public static string AppVersion
	{
		get => ReqContext<string>.Get("appVersion")!;
		set => ReqContext<string>.Set("appVersion", value);
	}

	public static bool IsLoggedIn
	{
		get => ReqContext<bool>.Get("isLoggedIn")!;
		set => ReqContext<bool>.Set("isLoggedIn", value);
	}
}

/// <summary>
/// ReqContext is a thread-local storage for request Context.
/// It is used to store data that is needed for the current request.
/// </summary>
/// <typeparam name="T">The type of the data to store</typeparam>
static class ReqContext<T>
{
	private static readonly ConcurrentDictionary<string, AsyncLocal<T>> state = new();

	/// <summary>
	/// Middlewares below can get this value. The above middlewares cannot
	/// </summary>
	/// <param name="name"></param>
	/// <param name="data"></param>
	public static void Set(string name, T data)
		=> state.GetOrAdd(name, _ => new AsyncLocal<T>()).Value = data;

	public static T? Get(string name) =>
		state.TryGetValue(name, out AsyncLocal<T>? data) ? data.Value : default;

}


