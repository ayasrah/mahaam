package mahaam.infra;

import jakarta.ws.rs.core.Response;

public class Exceptions {

	public abstract static class AppException extends RuntimeException {
		private String key;
		private int httpCode;

		public String getKey() {
			return key;
		}

		public int getHttpCode() {
			return httpCode;
		}

		public AppException(String message, int httpCode) {
			super(message);
			this.httpCode = httpCode;
		}

		public AppException(String message, int httpCode, String key) {
			super(message);
			this.httpCode = httpCode;
			this.key = key;
		}
	}

	public static class UnauthorizedException extends AppException {
		public UnauthorizedException(String message) {
			super(message, Response.Status.UNAUTHORIZED.getStatusCode());
		}
	}

	public static class ForbiddenException extends AppException {
		public ForbiddenException(String message) {
			super(message, Response.Status.FORBIDDEN.getStatusCode());
		}
	}

	public static class LogicException extends AppException {
		public LogicException(String message) {
			super(message, Response.Status.CONFLICT.getStatusCode());
		}

		public LogicException(String message, String key) {
			super(message, Response.Status.CONFLICT.getStatusCode(), key);
		}
	}

	public static class NotFoundException extends AppException {
		public NotFoundException(String message) {
			super(message, Response.Status.NOT_FOUND.getStatusCode());
		}
	}

	public static class InputException extends AppException {
		public InputException(String message) {
			super(message, Response.Status.BAD_REQUEST.getStatusCode());
		}
	}
}
