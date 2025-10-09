package mahaam.feat.user;

import java.util.List;
import java.util.UUID;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import jakarta.transaction.Transactional;
import mahaam.feat.plan.PlanRepo;
import mahaam.feat.user.UserModel.CreatedUser;
import mahaam.feat.user.UserModel.Device;
import mahaam.feat.user.UserModel.SuggestedEmail;
import mahaam.feat.user.UserModel.User;
import mahaam.feat.user.UserModel.VerifiedUser;
import mahaam.infra.Config;
import mahaam.infra.Email;
import mahaam.infra.Exceptions.InputException;
import mahaam.infra.Exceptions.UnauthorizedException;
import mahaam.infra.Log;
import mahaam.infra.Req;
import mahaam.infra.Security;

public interface UserService {
	CreatedUser create(Device device);

	String sendMeOtp(String email);

	VerifiedUser verifyOtp(String email, String sid, String otp);

	VerifiedUser refreshToken();

	void updateName(String name);

	void logout(UUID deviceId);

	void delete(String sid, String otp);

	List<Device> getDevices();

	List<SuggestedEmail> getSuggestedEmails();

	void deleteSuggestedEmail(UUID suggestedEmailId);
}

@ApplicationScoped
class DefaultUserService implements UserService {

	@Inject
	UserRepo userRepo;

	@Inject
	DeviceRepo deviceRepo;

	@Inject
	SuggestedEmailsRepo suggestedEmailsRepo;

	@Inject
	PlanRepo planRepo;

	@Override
	@Transactional
	public CreatedUser create(Device device) {
		UUID userId = userRepo.create();

		// add device
		Device deviceWithUserId = new Device();
		deviceWithUserId.id = device.id;
		deviceWithUserId.userId = userId;
		deviceWithUserId.platform = device.platform;
		deviceWithUserId.fingerprint = device.fingerprint;
		deviceWithUserId.info = device.info;
		deviceWithUserId.createdAt = device.createdAt;
		deviceRepo.deleteByFingerprint(device.fingerprint);
		UUID deviceId = deviceRepo.create(deviceWithUserId);

		String jwt = Security.createToken(userId.toString(), deviceId.toString());

		Log.info("User Created with id:" + userId + ", deviceId:" + deviceId + ".");
		CreatedUser createdUser = new CreatedUser();
		createdUser.id = userId;
		createdUser.deviceId = deviceId;
		createdUser.jwt = jwt;
		return createdUser;
	}

	@Override
	public String sendMeOtp(String email) {
		String verifySid;
		if (Config.testEmails.contains(email)) {
			verifySid = Config.testSID;
		} else {
			verifySid = Email.sendOtp(email);
		}

		if (verifySid != null) {
			Log.info("OTP sent to " + email);
		}

		return verifySid;
	}

	@Override
	@Transactional
	public VerifiedUser verifyOtp(String email, String sid, String otp) {
		String otpStatus;
		if (Config.testEmails.contains(email)
				&& Config.testSID.equals(sid)
				&& Config.testOTP.equals(otp)) {
			otpStatus = "approved";
		} else {
			otpStatus = Email.verifyOtp(otp, sid, email);
		}

		if (!"approved".equals(otpStatus)) {
			throw new InputException(
					"OTP not verified for " + email + ", status: " + otpStatus);
		}

		User user = userRepo.getOne(email);
		UUID userId = Req.getUserId();
		UUID deviceId = Req.getDeviceId();

		if (user == null) {
			userRepo.updateEmail(userId, email);
			Log.info("User loggedIn for " + email);
		} else {
			// move plans of current user to the one with email
			planRepo.updateUserId(userId, user.id);
			List<Device> devices = deviceRepo.getMany(user.id);
			if (devices != null && devices.size() >= 5) {
				Device lastDevice = devices.get(devices.size() - 1);
				deviceRepo.delete(lastDevice.id);
			}

			deviceRepo.updateUserId(deviceId, user.id);
			userRepo.delete(userId);
			Log.info("Merging userId:" + userId + " to " + user.id);
		}

		UUID newUserId = user == null ? userId : user.id;
		String jwt = Security.createToken(newUserId.toString(), deviceId.toString());

		Log.info("OTP verified for " + email);
		VerifiedUser verifiedUser = new VerifiedUser();
		verifiedUser.userId = newUserId;
		verifiedUser.deviceId = deviceId;
		verifiedUser.jwt = jwt;
		verifiedUser.userFullName = user != null ? user.name : null;
		verifiedUser.email = email;
		return verifiedUser;
	}

	@Override
	public VerifiedUser refreshToken() {
		UUID userId = Req.getUserId();
		UUID deviceId = Req.getDeviceId();
		User user = userRepo.getOne(userId);
		String jwt = Security.createToken(userId.toString(), deviceId.toString());

		VerifiedUser verifiedUser = new VerifiedUser();
		verifiedUser.userId = userId;
		verifiedUser.deviceId = deviceId;
		verifiedUser.jwt = jwt;
		verifiedUser.userFullName = user != null ? user.name : null;
		verifiedUser.email = user != null ? user.email : null;
		return verifiedUser;
	}

	@Override
	public void updateName(String name) {
		userRepo.updateName(Req.getUserId(), name);
	}

	@Override
	public void logout(UUID deviceId) {
		Device device = deviceRepo.getOne(deviceId);
		if (device == null || !device.userId.equals(Req.getUserId())) {
			throw new UnauthorizedException("Invalid deviceId");
		}
		deviceRepo.delete(deviceId);
	}

	@Override
	public void deleteSuggestedEmail(UUID suggestedEmailId) {
		SuggestedEmail suggestedEmail = suggestedEmailsRepo.getOne(suggestedEmailId);
		if (suggestedEmail == null || !suggestedEmail.userId.equals(Req.getUserId())) {
			throw new UnauthorizedException("Invalid suggestedEmailId");
		}
		suggestedEmailsRepo.delete(suggestedEmailId);
	}

	@Override
	@Transactional
	public void delete(String sid, String otp) {
		User user = userRepo.getOne(Req.getUserId());

		String otpStatus;
		if (Config.testEmails.contains(user.email)
				&& Config.testSID.equals(sid)
				&& Config.testOTP.equals(otp)) {
			otpStatus = "approved";
		} else {
			otpStatus = Email.verifyOtp(otp, sid, user.email);
		}

		if (!"approved".equals(otpStatus)) {
			throw new InputException(
					"OTP not approved for " + user.email + ", status: " + otpStatus);
		}

		suggestedEmailsRepo.deleteManyByEmail(user.email);
		userRepo.delete(Req.getUserId());
	}

	@Override
	public List<Device> getDevices() {
		return deviceRepo.getMany(Req.getUserId());
	}

	@Override
	public List<SuggestedEmail> getSuggestedEmails() {
		return suggestedEmailsRepo.getMany(Req.getUserId());
	}
}
