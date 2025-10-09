namespace Mahaam.Infra;

using System;
using System.Net;
using System.Net.Sockets;
using System.Text.RegularExpressions;
using Mahaam.Infra.Monitoring;
using Npgsql;

public class Starter
{
	public static void Init(WebApplication? app)
	{
		InitDB();
		Email.Init();

		var health = new Health()
		{
			Id = Guid.NewGuid(),
			ApiName = Config.ApiName,
			ApiVersion = Config.ApiVersion,
			NodeIP = GetNodeIP(),
			NodeName = Environment.MachineName,
			EnvName = Config.EnvName
		};
		App.HealthService.ServerStarted(health);
		Cache.Init(health);
		var startMsg = $"✓ {Config.ApiName}-v{Config.ApiVersion}/{Cache.NodeIP}-{Cache.NodeName} started with healthID={Cache.HealthId}";
		Log.Info(startMsg);
		Thread.Sleep(2000);
		App.HealthService.StartSendingPulses();
	}

	private static void InitDB()
	{
		using (var connection = new NpgsqlConnection(Config.DbUrl))
		{
			connection.Open();
			connection.Close();
		}
		string pattern = "Host=([^;]+)";

		Match match = Regex.Match(Config.DbUrl, pattern);
		var host = match.Success ? match.Groups[1].Value : "";
		Log.Info($"✓ Connected to DB on server {host}");
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