package mahaam.infra;

import org.jboss.logging.Logger;

import io.quarkus.arc.Arc;
import mahaam.infra.monitor.LogRepo;

public class Log {
	private static final Logger logger = Logger.getLogger(Log.class);

	public static void error(String error) {
		var trafficId = Req.getTrafficId();
		var message = trafficId != null ? "TrafficId: " + trafficId + ", " + error : error;

		logger.error(message);
		try {
			var logRepo = Arc.container().instance(LogRepo.class).get();
			if (logRepo != null) {
				logRepo.create(trafficId, "error", error);
			}
		} catch (Exception e) {
			logger.error("Error in Log.Error: " + e.toString());
		}
	}

	public static void info(String info) {
		var trafficId = Req.getTrafficId();
		var message = trafficId != null ? "TrafficId: " + trafficId + ", " + info : info;

		logger.info(message);
		try {
			var logRepo = Arc.container().instance(LogRepo.class).get();
			if (logRepo != null) {
				logRepo.create(trafficId, "info", info);
			}
		} catch (Exception e) {
			logger.error("Error in Log.Info: " + e.toString());
		}
	}
}