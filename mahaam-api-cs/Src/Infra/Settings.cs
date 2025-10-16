

namespace Mahaam.Infra;

public class Settings
{
	public ApiSettings Api { get; set; }
	public EmailSettings Email { get; set; }
	public LoggingSettings Logging { get; set; }
	public string DbUrl { get; set; }
	public bool LogReqEnabled { get; set; }
}

public class LoggingSettings
{
	public string File { get; set; }
	public int FileSizeLimit { get; set; }
	public int FileCountLimit { get; set; }
	public string OutputTemplate { get; set; }
}

public class EmailSettings
{
	public string AccountSid { get; set; }
	public string VerificationServiceSid { get; set; }
	public string AuthToken { get; set; }
	public List<string> TestEmails { get; set; }
	public string TestSID { get; set; }
	public string TestOTP { get; set; }
}


public class ApiSettings
{
	public string Name { get; set; }
	public string Version { get; set; }
	public string EnvName { get; set; }
	public int HttpPort { get; set; }
	public string TokenSecretKey { get; set; }
}