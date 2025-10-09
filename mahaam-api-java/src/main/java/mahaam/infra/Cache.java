package mahaam.infra;

import java.util.UUID;

import jakarta.enterprise.context.ApplicationScoped;
import mahaam.infra.monitor.MonitorModel.Health;

@ApplicationScoped
public class Cache {

	public static void init(Health health) {
		_health = health;
	}

	private static Health _health;

	public static String getNodeIP() {
		return _health != null ? _health.nodeIP : "";
	}

	public static String getNodeName() {
		return _health != null ? _health.nodeName : "";
	}

	public static String getApiName() {
		return _health != null ? _health.apiName : "";
	}

	public static String getApiVersion() {
		return _health != null ? _health.apiVersion : "";
	}

	public static String getEnvName() {
		return _health != null ? _health.envName : "";
	}

	public static UUID getHealthId() {
		return _health != null ? _health.id : null;
	}
}
