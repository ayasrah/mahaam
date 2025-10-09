

namespace Mahaam.Infra;

public class Config
{
	private static IConfiguration? _configuration;

	public static void Init(IConfiguration configuration)
	{
		_configuration = configuration;
	}

	public static string ApiName => GetValue("apiName");
	public static string ApiVersion => GetValue("apiVersion");
	public static string EnvName => GetValue("envName");
	public static string DbUrl => GetValue("dbUrl");
	public static string LogFile => GetValue("logFile");
	public static int LogFileSizeLimit => int.Parse(GetValue("logFileSizeLimit"));
	public static int LogFileCountLimit => int.Parse(GetValue("logFileCountLimit"));
	public static string LogFileOutputTemplate => GetValue("logFileOutputTemplate");
	public static int HttpPort => int.Parse(GetValue("httpPort"));
	public static string TokenSecretKey => GetValue("tokenSecretKey");
	public static string EmailAccountSid => GetValue("emailAccountSid");
	public static string EmailVerificationServiceSid => GetValue("emailVerificationServiceSid");
	public static string EmailAuthToken => GetValue("emailAuthToken");
	private static string TestEmailsStr => GetValue("testEmails");
	public static List<string> TestEmails => TestEmailsStr.Split(',').ToList();
	public static string TestSID => GetValue("testSID");
	public static string TestOTP => GetValue("testOTP");
	public static bool LogReqEnabled => bool.Parse(GetValue("logReqEnabled"));


	private static string GetValue(string key)
	{
		if (_configuration == null)
			throw new InvalidOperationException("Configuration not initialized. Call Config.Initialize() first.");

		return _configuration[key] ?? throw new ArgumentException($"Configuration key '{key}' not found.");
	}
}

