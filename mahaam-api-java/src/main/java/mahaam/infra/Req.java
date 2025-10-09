package mahaam.infra;

import java.util.UUID;

import io.vertx.core.Vertx;

class ReqContext {
	public static <T> void set(String key, T value) {
		var ctx = Vertx.currentContext();
		if (ctx != null)
			ctx.putLocal(key, value);
	}

	@SuppressWarnings("unchecked")
	public static <T> T get(String key) {
		var ctx = Vertx.currentContext();
		return ctx != null ? (T) ctx.getLocal(key) : null;
	}

	public static void clear(String key) {
		var ctx = Vertx.currentContext();
		if (ctx != null)
			ctx.removeLocal(key);
	}
}

public class Req {
	public static UUID getTrafficId() {
		return ReqContext.get("trafficId");
	}

	public static void setTrafficId(UUID value) {
		ReqContext.set("trafficId", value);
	}

	public static UUID getUserId() {
		return ReqContext.get("userId");
	}

	public static void setUserId(UUID value) {
		ReqContext.set("userId", value);
	}

	public static UUID getDeviceId() {
		return ReqContext.get("deviceId");
	}

	public static void setDeviceId(UUID value) {
		ReqContext.set("deviceId", value);
	}

	public static String getAppStore() {
		return ReqContext.get("appStore");
	}

	public static void setAppStore(String value) {
		ReqContext.set("appStore", value);
	}

	public static String getAppVersion() {
		return ReqContext.get("appVersion");
	}

	public static void setAppVersion(String value) {
		ReqContext.set("appVersion", value);
	}

	public static boolean isLoggedIn() {
		return ReqContext.get("isLoggedIn");
	}

	public static void setLoggedIn(boolean value) {
		ReqContext.set("isLoggedIn", value);
	}

	public static void setStartTime(long value) {
		ReqContext.set("startTime", value);
	}

	public static long getStartTime() {
		return ReqContext.get("startTime");
	}
}
