using System;
using System.Collections.Generic;
using Mahaam.Infra.Monitoring;

namespace Mahaam.Infra;

public interface ILog
{
	void Info(string info);
	void Error(string error);
}

public class Log(ILogRepo logRepo) : ILog
{
	public void Info(string info)
	{
		var message = Req.TrafficId != Guid.Empty ? $"TrafficId: {Req.TrafficId}, {info}" : info;

		Serilog.Log.Information(message);
		logRepo.Create("Info", info, Req.TrafficId);
	}

	public void Error(string error)
	{
		var message = Req.TrafficId != Guid.Empty ? $"TrafficId: {Req.TrafficId}, {error}" : error;

		Serilog.Log.Error(message);
		logRepo.Create("Error", error, Req.TrafficId);
	}
}
