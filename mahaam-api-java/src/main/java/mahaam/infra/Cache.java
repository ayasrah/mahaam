package mahaam.infra;

import java.util.UUID;

import jakarta.enterprise.context.ApplicationScoped;
import mahaam.infra.monitor.MonitorModel.Health;

@ApplicationScoped
public class Cache {

	public void init(Health health) {
		_health = health;
	}

	private Health _health;

	public String nodeIP() {
		return _health != null ? _health.nodeIP : "";
	}

	public String nodeName() {
		return _health != null ? _health.nodeName : "";
	}

	public UUID healthId() {
		return _health != null ? _health.id : null;
	}
}
