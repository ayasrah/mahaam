using System.Net.Mail;

namespace Mahaam.Infra;
public class Rule
{
	public static void Required(string? value, string name)
	{
		if (string.IsNullOrWhiteSpace(value)) throw new InputException(name + " is required");
	}

	public static void OneAtLeastRequired(List<object?> value, string message)
	{
		if (value.All((item) => item is null || string.IsNullOrWhiteSpace(item.ToString()))) throw new InputException(message);
	}

	public static void Required(Guid? value, string name)
	{
		if ((value is null) || Guid.Empty.Equals(value)) throw new InputException(name + " is required");
	}

	public static void Required(bool? value, string name)
	{
		if (value == null) throw new InputException(name + " is required");
	}

	public static void Required(int? value, string name)
	{
		if (value == null) throw new InputException(name + " is required");
	}

	public static void Required(object? value, string name)
	{
		if (value == null) throw new InputException(name + " is required");
	}

	public static void In(string item, List<string> list)
	{
		if (!list.Contains(item))
			throw new InputException($"{item} is not in [{string.Join(",", list)}]");
	}

	public static void FailIf(bool condition, string message)
	{
		if (condition) throw new InputException(message);
	}

	public static void ValidateEmail(string email)
	{
		Required(email, "email");
		try
		{
			var addr = new MailAddress(email);
			FailIf(!email.Equals(addr.Address), "Invalid email");
		}
		catch
		{
			throw new InputException("Invalid email");
		}
	}
}
