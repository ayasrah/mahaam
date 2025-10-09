package mahaam.infra;

import java.util.List;
import java.util.UUID;
import java.util.regex.Pattern;

import mahaam.infra.Exceptions.InputException;

public class Rule {

	private static final Pattern EMAIL_PATTERN = Pattern.compile(
			"^[a-zA-Z0-9_+&*-]+(?:\\.[a-zA-Z0-9_+&*-]+)*@(?:[a-zA-Z0-9-]+\\.)+[a-zA-Z]{2,7}$");

	public static void required(String value, String name) {
		if (value == null || value.trim().isEmpty())
			throw new InputException(name + " is required");
	}

	public static void oneAtLeastRequired(List<?> values, String message) {
		if (values.stream().allMatch(item -> item == null || item.toString().trim().isEmpty()))
			throw new InputException(message);
	}

	public static void required(UUID value, String name) {
		if (value == null)
			throw new InputException(name + " is required");
	}

	public static void required(Boolean value, String name) {
		if (value == null)
			throw new InputException(name + " is required");
	}

	public static void required(Integer value, String name) {
		if (value == null)
			throw new InputException(name + " is required");
	}

	public static void required(Object value, String name) {
		if (value == null)
			throw new InputException(name + " is required");
	}

	public static void in(String item, List<String> list) {
		if (!list.contains(item))
			throw new InputException(item + " is not in [" + String.join(",", list) + "]");
	}

	public static void failIf(boolean condition, String message) {
		if (condition)
			throw new InputException(message);
	}

	public static void validateEmail(String email) {
		required(email, "email");
		if (!EMAIL_PATTERN.matcher(email).matches())
			throw new InputException("Invalid email");
	}
}
