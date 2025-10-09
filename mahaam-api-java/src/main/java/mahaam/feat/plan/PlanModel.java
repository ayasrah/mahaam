package mahaam.feat.plan;

import java.time.Instant;
import java.time.LocalDate;
import java.util.List;
import java.util.UUID;

import mahaam.feat.user.UserModel.User;

public class PlanModel {

	public static class Plan {
		public UUID id;
		public String title;
		public String type;
		public int sortOrder;
		public Instant starts;
		public Instant ends;
		public String donePercent;
		public Instant createdAt;
		public Instant updatedAt;
		public List<User> members;
		public boolean isShared;
		public User user;
	}

	public static class PlanIn {
		public UUID id;
		public String title;
		public LocalDate starts;
		public LocalDate ends;
	}

	public static class PlanType {
		public static final String MAIN = "Main";
		public static final String ARCHIVED = "Archived";
		public static final List<String> ALL = List.of(MAIN, ARCHIVED);
	}
}
