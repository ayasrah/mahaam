package mahaam.infra;

import jakarta.ws.rs.core.MediaType;

public class Http {
	// HTTP Status Codes
	public static final int OK = 200;
	public static final int Created = 201;
	public static final int NoContent = 204;
	public static final int BadRequest = 400;
	public static final int UnAuthorized = 401;
	public static final int Forbidden = 403;
	public static final int NotFound = 404;
	public static final int Conflict = 409;
	public static final int ServerError = 500;

	// Media Types
	public static final String JsonMedia = MediaType.APPLICATION_JSON;
	public static final String TextMedia = MediaType.TEXT_PLAIN;
	public static final String FormMedia = MediaType.APPLICATION_FORM_URLENCODED;
}
