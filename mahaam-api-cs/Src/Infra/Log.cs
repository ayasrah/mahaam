using System;
using System.Collections.Generic;

namespace Mahaam.Infra;
public class Log
{
	public static void Info(string info)
	{
		var message = Req.TrafficId != Guid.Empty ? $"TrafficId: {Req.TrafficId}, {info}" : info;

		Serilog.Log.Information(message);
		App.LogRepo.Create("Info", info, Req.TrafficId);
	}

	public static void Error(string error)
	{
		var message = Req.TrafficId != Guid.Empty ? $"TrafficId: {Req.TrafficId}, {error}" : error;

		Serilog.Log.Error(message);
		App.LogRepo.Create("Error", error, Req.TrafficId);
	}
}
