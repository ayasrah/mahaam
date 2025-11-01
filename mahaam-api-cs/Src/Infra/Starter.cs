namespace Mahaam.Infra;

using System;
using System.Net;
using System.Net.Sockets;
using System.Text.RegularExpressions;
using Mahaam.Feat.Users;
using Mahaam.Infra.Monitoring;
using Microsoft.Extensions.Options;
using Npgsql;

public class Starter
{
	public static void Init(WebApplication? app)
	{
		var settings = app.Services.GetRequiredService<Settings>();
		InitDB(app, settings);
		app.Services.GetService<IEmail>()?.Init();

		var health = new Health()
		{
			Id = Guid.NewGuid(),
			ApiName = settings.Api.Name,
			ApiVersion = settings.Api.Version,
			NodeIP = GetNodeIP(),
			NodeName = Environment.MachineName,
			EnvName = settings.Api.EnvName
		};
		app.Services.GetService<IHealthService>()?.ServerStarted(health);
		var cache = app.Services.GetService<ICache>();
		cache?.Init(health);
		var startMsg = $"✓ {settings.Api.Name}-v{settings.Api.Version}/{cache.NodeIP()}-{cache.NodeName()} started with healthID={cache.HealthId()}";
		app.Services.GetService<ILog>()?.Info(startMsg);
		Thread.Sleep(2000);
		app.Services.GetService<IHealthService>()?.StartSendingPulses();
	}

	private static void InitDB(WebApplication app, Settings settings)
	{
		using (var connection = new NpgsqlConnection(settings.DbUrl))
		{
			connection.Open();
			connection.Close();
		}
		string pattern = "Host=([^;]+)";

		Match match = Regex.Match(settings.DbUrl, pattern);
		var host = match.Success ? match.Groups[1].Value : "";
		app.Services.GetService<ILog>()?.Info($"✓ Connected to DB on server {host}");
	}

	private static string GetNodeIP()
	{
		string? ipAddress = null;
		try
		{
			using var socket = new Socket(AddressFamily.InterNetwork, SocketType.Dgram, ProtocolType.Udp);
			socket.Connect(IPAddress.Parse("8.8.8.8"), 10002);
			ipAddress = ((IPEndPoint)socket.LocalEndPoint).Address.ToString();
		}
		catch (Exception e)
		{
			Console.WriteLine($"An error occurred while getting the local IP address: {e.Message}");
		}
		return ipAddress;
	}
}