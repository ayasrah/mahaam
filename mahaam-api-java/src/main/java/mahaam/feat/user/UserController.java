package mahaam.feat.user;

import java.util.List;
import java.util.UUID;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import jakarta.ws.rs.Consumes;
import jakarta.ws.rs.DELETE;
import jakarta.ws.rs.FormParam;
import jakarta.ws.rs.GET;
import jakarta.ws.rs.PATCH;
import jakarta.ws.rs.POST;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.core.MediaType;
import jakarta.ws.rs.core.Response;
import mahaam.feat.user.UserModel.CreatedUser;
import mahaam.feat.user.UserModel.Device;
import mahaam.feat.user.UserModel.SuggestedEmail;
import mahaam.feat.user.UserModel.VerifiedUser;
import mahaam.infra.Json;
import mahaam.infra.Rule;

public interface UserController {
	Response create(String platform, boolean isPhysicalDevice, String deviceFingerprint, String deviceInfo);

	Response sendMeOtp(String email);

	Response verifyOtp(String email, String sid, String otp);

	Response refreshToken();

	Response updateName(String name);

	Response logout(UUID deviceId);

	Response delete(String sid, String otp);

	Response getDevices();

	Response getSuggestedEmails();

	Response deleteSuggestedEmail(UUID suggestedEmailId);
}

@ApplicationScoped
@Path("/users")
@Consumes(MediaType.APPLICATION_FORM_URLENCODED)
@Produces(MediaType.APPLICATION_JSON)
class DefaultUserController implements UserController {

	@Inject
	UserService userService;

	@POST
	@Path("/create")
	public Response create(
			@FormParam("platform") String platform,
			@FormParam("isPhysicalDevice") boolean isPhysicalDevice,
			@FormParam("deviceFingerprint") String deviceFingerprint,
			@FormParam("deviceInfo") String deviceInfo) {

		Rule.required(platform, "platform");
		Rule.required(deviceFingerprint, "deviceFingerprint");
		Rule.required(deviceInfo, "deviceInfo");
		Rule.failIf(!isPhysicalDevice, "Device should be real not simulator");

		Device device = new Device();
		device.platform = platform;
		device.fingerprint = deviceFingerprint;
		device.info = deviceInfo;

		CreatedUser createdUser = userService.create(device);
		return Response.ok().entity(Json.toString(createdUser)).build();
	}

	@POST
	@Path("/send-me-otp")
	public Response sendMeOtp(@FormParam("email") String email) {
		Rule.validateEmail(email);
		String verificationSid = userService.sendMeOtp(email);
		return Response.ok().entity(Json.toString(verificationSid)).build();
	}

	@POST
	@Path("/verify-otp")
	public Response verifyOtp(
			@FormParam("email") String email,
			@FormParam("sid") String sid,
			@FormParam("otp") String otp) {

		Rule.required(email, "email");
		Rule.required(sid, "sid");
		Rule.required(otp, "otp");

		VerifiedUser verifiedUser = userService.verifyOtp(email, sid, otp);
		return Response.ok().entity(Json.toString(verifiedUser)).build();
	}

	@POST
	@Path("/refresh-token")
	public Response refreshToken() {
		VerifiedUser verifiedUser = userService.refreshToken();
		return Response.ok().entity(Json.toString(verifiedUser)).build();
	}

	@PATCH
	@Path("/name")
	public Response updateName(@FormParam("name") String name) {
		Rule.required(name, "name");
		userService.updateName(name);
		return Response.ok().build();
	}

	@POST
	@Path("/logout")
	public Response logout(@FormParam("deviceId") UUID deviceId) {
		Rule.required(deviceId, "deviceId");
		userService.logout(deviceId);
		return Response.ok().build();
	}

	@DELETE
	public Response delete(@FormParam("sid") String sid, @FormParam("otp") String otp) {
		Rule.required(sid, "sid");
		Rule.required(otp, "otp");
		userService.delete(sid, otp);
		return Response.noContent().build();
	}

	@GET
	@Path("/devices")
	public Response getDevices() {
		List<Device> devices = userService.getDevices();
		return Response.ok().entity(Json.toString(devices)).build();
	}

	@GET
	@Path("/suggested-emails")
	public Response getSuggestedEmails() {
		List<SuggestedEmail> suggestedEmails = userService.getSuggestedEmails();
		return Response.ok().entity(Json.toString(suggestedEmails)).build();
	}

	@DELETE
	@Path("/suggested-emails")
	public Response deleteSuggestedEmail(@FormParam("suggestedEmailId") UUID suggestedEmailId) {
		Rule.required(suggestedEmailId, "suggestedEmailId");
		userService.deleteSuggestedEmail(suggestedEmailId);
		return Response.noContent().build();
	}
}
