package mahaam.infra;

import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.util.List;
import java.util.UUID;

import jakarta.annotation.Priority;
import jakarta.inject.Inject;
import jakarta.ws.rs.container.ContainerRequestContext;
import jakarta.ws.rs.container.ContainerRequestFilter;
import jakarta.ws.rs.container.ContainerResponseContext;
import jakarta.ws.rs.container.ContainerResponseFilter;
import jakarta.ws.rs.core.Response;
import jakarta.ws.rs.ext.ExceptionMapper;
import jakarta.ws.rs.ext.Provider;
import mahaam.infra.Exceptions.AppException;
import mahaam.infra.Exceptions.NotFoundException;
import mahaam.infra.monitor.MonitorModel.Traffic;
import mahaam.infra.monitor.MonitorModel.TrafficHeaders;
import mahaam.infra.monitor.TrafficRepo;

@Provider
@Priority(1000)
public class Filters implements ContainerRequestFilter, ContainerResponseFilter {

	@Inject
	TrafficRepo trafficRepo;

	@Inject
	Security security;

	private static final List<String> BYPASS_AUTH_PATHS = List.of("/swagger", "/health", "/users/create", "/audit/info",
			"/audit/error");

	private static final List<String> NO_TRAFFIC_PATHS = List.of("/swagger", "/health", "/audit");

	@Override
	public void filter(ContainerRequestContext reqCtx) throws IOException {
		Req.setTrafficId(UUID.randomUUID());
		Req.setStartTime(System.currentTimeMillis());
		String path = reqCtx.getUriInfo().getPath();
		String pathBase = reqCtx.getUriInfo().getBaseUri().getPath();

		// Validate path base
		if (!pathBase.contains("mahaam-api")) {
			throw new NotFoundException("Invalid path base");
		}

		// Check required headers
		String appStore = reqCtx.getHeaderString("x-app-store");
		String appVersion = reqCtx.getHeaderString("x-app-version");

		// if (appStore == null || appVersion == null) {
		// throw new UnauthorizedException("Required headers not exists");
		// }

		// Store app info in request data
		Req.setAppStore(appStore);
		Req.setAppVersion(appVersion);

		// Handle authentication
		if (!shouldBypassAuth(path) || false) {
			var authResult = security.validateAndExtractJwt(reqCtx);
			Req.setUserId(authResult.getUserId());
			Req.setDeviceId(authResult.getDeviceId());
			Req.setLoggedIn(authResult.isLoggedIn());
		}

		// Read and store payload
		if (Config.logReqEnabled) {
			String reqBody = getReqBody(reqCtx);
			reqCtx.setProperty("reqBody", reqBody);
		}
	}

	@Override
	public void filter(
			ContainerRequestContext requestCtx, ContainerResponseContext responseCtx)
			throws IOException {
		try {

			String path = requestCtx.getUriInfo().getPath();

			// Log traffic if not in excluded paths
			if (!shouldSkipTraffic(path)) {
				logTraffic(requestCtx, responseCtx);
			}

		} catch (Exception e) {
			Log.error("Error in response filter: " + e.toString());
		}
	}

	private boolean shouldBypassAuth(String path) {
		return BYPASS_AUTH_PATHS.stream().anyMatch(path::startsWith);
	}

	private boolean shouldSkipTraffic(String path) {
		return NO_TRAFFIC_PATHS.stream().anyMatch(path::startsWith);
	}

	private String getReqBody(ContainerRequestContext requestContext) throws IOException {
		if (requestContext.hasEntity()) {
			InputStream entityStream = requestContext.getEntityStream();
			ByteArrayOutputStream baos = new ByteArrayOutputStream();
			byte[] buffer = new byte[1024];
			int bytesRead;

			while ((bytesRead = entityStream.read(buffer)) != -1) {
				baos.write(buffer, 0, bytesRead);
			}

			byte[] entity = baos.toByteArray();
			requestContext.setEntityStream(new ByteArrayInputStream(entity));

			String payload = new String(entity);
			return payload.isEmpty() ? null : payload;
		}
		return null;
	}

	private void logTraffic(
			ContainerRequestContext requestContext,
			ContainerResponseContext responseContext) {
		try {
			String path = requestContext.getUriInfo().getPath();
			String queryString = requestContext.getUriInfo().getRequestUri().getQuery();
			String fullPath = path + (queryString != null ? "?" + queryString : "");

			String request = null;
			String response = null;

			// Don't log request/response for successful requests
			boolean isFailedResponse = responseContext.getStatus() >= 400;
			if (isFailedResponse) {
				request = (String) requestContext.getProperty("reqBody");
				response = getResponseEntity(responseContext);
			}

			// Don't log response for user endpoints or if empty
			if (path.startsWith("/user") || (response != null && response.isEmpty())) {
				response = null;
			}

			// Create traffic headers
			TrafficHeaders trafficHeaders = new TrafficHeaders();
			trafficHeaders.userId = (UUID) Req.getUserId();
			trafficHeaders.deviceId = (UUID) Req.getDeviceId();
			trafficHeaders.appVersion = (String) Req.getAppVersion();
			trafficHeaders.appStore = (String) Req.getAppStore();

			long elapsed = System.currentTimeMillis() - Req.getStartTime();

			Traffic traffic = new Traffic(
					Req.getTrafficId(),
					Cache.getHealthId(),
					requestContext.getMethod(),
					fullPath,
					responseContext.getStatus(),
					elapsed,
					Json.toString(trafficHeaders),
					request,
					response);

			trafficRepo.create(traffic);

		} catch (Exception e) {
			Log.error("Error logging traffic: " + e.toString());
		}
	}

	private String getResponseEntity(ContainerResponseContext responseContext) {
		if (responseContext.hasEntity()) {
			try {
				return responseContext.getEntity().toString();
			} catch (Exception e) {
				Log.error("Error getting response entity: " + e.toString());
				return null;
			}
		}
		return null;
	}

}

@Provider
class ExceptionFilter implements ExceptionMapper<Throwable> {

	@Override
	public Response toResponse(Throwable exception) {

		String response = Json.toString(exception.getMessage());
		int statusCode = Http.ServerError;
		exception.printStackTrace();
		Log.error(exception.toString());
		if (exception instanceof AppException) {
			AppException appException = (AppException) exception;
			statusCode = appException.getHttpCode();
			String key = appException.getKey();
			if (key != null && !key.isEmpty()) {
				ErrorResponse errorResponse = new ErrorResponse(key, exception.getMessage());
				response = Json.toString(errorResponse);
			}
		}

		return Response.status(statusCode).entity(response).type(Http.JsonMedia).build();
	}

	private static class ErrorResponse {
		private String key;
		private String error;

		public ErrorResponse(String key, String error) {
			this.key = key;
			this.error = error;
		}

		// Getters (used by JSON serialization)
		@SuppressWarnings("unused")
		public String getKey() {
			return key;
		}

		@SuppressWarnings("unused")
		public String getError() {
			return error;
		}
	}
}
