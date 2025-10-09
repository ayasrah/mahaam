package mahaam.infra.monitor;

import java.util.UUID;

public class MonitorModel {
	public static class Health {
		public UUID id;
		public String apiName;
		public String apiVersion;
		public String nodeIP;
		public String nodeName;
		public String envName;
	}

	public static class Traffic {
		public UUID id;
		public UUID healthId;
		public String method;
		public String path;
		public Integer code;
		public Long elapsed;
		public String headers;
		public String request;
		public String response;

		public Traffic() {
		}

		public Traffic(
				UUID id,
				UUID healthId,
				String method,
				String path,
				Integer code,
				Long elapsed,
				String headers,
				String request,
				String response) {
			this.id = id;
			this.healthId = healthId;
			this.method = method;
			this.path = path;
			this.code = code;
			this.elapsed = elapsed;
			this.headers = headers;
			this.request = request;
			this.response = response;
		}
	}

	public static class TrafficHeaders {
		public UUID userId;
		public UUID deviceId;
		public String appVersion;
		public String appStore;
	}
}