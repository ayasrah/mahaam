package mahaam.infra;

import java.net.DatagramSocket;
import java.net.InetAddress;
import java.net.InetSocketAddress;
import java.util.UUID;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

import io.quarkus.runtime.ShutdownEvent;
import io.quarkus.runtime.StartupEvent;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.enterprise.event.Observes;
import jakarta.inject.Inject;
import mahaam.infra.monitor.HealthService;
import mahaam.infra.monitor.MonitorModel.Health;

@ApplicationScoped
public class Starter {

	@Inject
	HealthService healthService;

	@Inject
	DB db;

	void onStart(@Observes StartupEvent ev) {
		initDB();
		Email.init(Config.emailAccountSid, Config.emailVerificationServiceSid, Config.emailAuthToken);

		Health health = new Health();
		health.id = UUID.randomUUID();
		health.apiName = Config.apiName;
		health.apiVersion = Config.apiVersion;
		health.nodeIP = getNodeIP();
		health.nodeName = getNodeName();
		health.envName = Config.envName;

		healthService.serverStarted(health);
		Cache.init(health);

		String startMsg = String.format(
				"✓ %s-v%s/%s-%s started with healthID=%s",
				Config.apiName,
				Config.apiVersion,
				Cache.getNodeIP(),
				Cache.getNodeName(),
				Cache.getHealthId());
		Log.info(startMsg);

		try {
			Thread.sleep(2000);
		} catch (InterruptedException e) {
			Thread.currentThread().interrupt();
		}

		healthService.startSendingPulses();
	}

	void onStop(@Observes ShutdownEvent ev) {
		String stopMsg = String.format(
				"✓ %s-v%s/%s-%s stopped with healthID=%s",
				Config.apiName,
				Config.apiVersion,
				Cache.getNodeIP(),
				Cache.getNodeName(),
				Cache.getHealthId());
		Log.info(stopMsg);
		healthService.serverStopped();
	}

	private void initDB() {
		db.init();
		String pattern = "jdbc:postgresql://([^:/]+)";
		Pattern regex = Pattern.compile(pattern);
		Matcher match = regex.matcher(Config.dbUrl);
		String host = match.find() ? match.group(1) : "unknown";
		Log.info("✓ Connected to DB on server " + host);
	}

	private String getNodeIP() {
		try (DatagramSocket socket = new DatagramSocket()) {
			socket.connect(new InetSocketAddress("8.8.8.8", 10002));
			return socket.getLocalAddress().getHostAddress();
		} catch (Exception e) {
			Log.error("An error occurred while getting the local IP address: " + e.getMessage());
			return "127.0.0.1";
		}
	}

	private static String getNodeName() {
		try {
			return InetAddress.getLocalHost().getHostName();
		} catch (Exception e) {
			System.err.println("Failed to get node name: " + e.getMessage());
			throw new RuntimeException("Failed to get node name", e);
		}
	}
}
