package mahaam.feat.user;

import java.time.Instant;
import java.util.UUID;

public class UserModel {

	public static class User {
		public UUID id;
		public String email;
		public String status;
		public String name;
	}

	public static class CreatedUser {
		public UUID id;
		public UUID deviceId;
		public String jwt;
	}

	public static class Device {
		public UUID id;
		public UUID userId;
		public String platform;
		public String fingerprint;
		public String info;
		public Instant createdAt;

		public Device() {
		}

		public Device(UUID id, UUID userId, String platform, String fingerprint, String info, Instant createdAt) {
			this.id = id;
			this.userId = userId;
			this.platform = platform;
			this.fingerprint = fingerprint;
			this.info = info;
			this.createdAt = createdAt;
		}
	}

	public static class SuggestedEmail {
		public UUID id;
		public UUID userId;
		public String email;
		public Instant createdAt;
	}

	public static class VerifiedUser {
		public UUID userId;
		public UUID deviceId;
		public String jwt;
		public String userFullName;
		public String email;
	}
}
