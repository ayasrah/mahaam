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
	private readonly ILogRepo _logRepo = logRepo;
	public void Info(string info)
	{
		var message = Req.TrafficId != Guid.Empty ? $"TrafficId: {Req.TrafficId}, {info}" : info;

		Serilog.Log.Information(message);
		_logRepo.Create("Info", info, Req.TrafficId);
	}

	public void Error(string error)
	{
		var message = Req.TrafficId != Guid.Empty ? $"TrafficId: {Req.TrafficId}, {error}" : error;

		Serilog.Log.Error(message);
		_logRepo.Create("Error", error, Req.TrafficId);
	}
}
