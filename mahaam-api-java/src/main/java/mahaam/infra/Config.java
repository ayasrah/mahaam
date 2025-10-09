package mahaam.infra;

import java.util.List;

import org.eclipse.microprofile.config.ConfigProvider;

/**
 * Configuration class that reads from .env file using MicroProfile Config.
 *
 * <p>
 * Configuration priority (highest to lowest): 1. System properties
 * (-Dkey=value) 2. Environment
 * variables 3. .env file (in project root) 4. application.properties
 *
 * <p>
 * To set up: 1. Copy env.example to .env 2. Update .env with your actual values
 * 3. The .env file
 * is gitignored for security
 */
public class Config {

	public static final String apiName = ConfigProvider.getConfig().getValue("apiName", String.class);
	public static final String apiVersion = ConfigProvider.getConfig().getValue("apiVersion", String.class);
	public static final String envName = ConfigProvider.getConfig().getValue("envName", String.class);
	public static final String dbUrl = ConfigProvider.getConfig().getValue("dbUrl", String.class);
	public static final String logFile = ConfigProvider.getConfig().getValue("logFile", String.class);
	public static final int logFileSizeLimit = ConfigProvider.getConfig().getValue("logFileSizeLimit", Integer.class);
	public static final int logFileCountLimit = ConfigProvider.getConfig().getValue("logFileCountLimit", Integer.class);
	public static final String logFileOutputTemplate = ConfigProvider.getConfig().getValue("logFileOutputTemplate",
			String.class);
	public static final int httpPort = ConfigProvider.getConfig().getValue("httpPort", Integer.class);
	public static final String tokenSecretKey = ConfigProvider.getConfig().getValue("tokenSecretKey", String.class);
	public static final String emailAccountSid = ConfigProvider.getConfig().getValue("emailAccountSid", String.class);
	public static final String emailVerificationServiceSid = ConfigProvider.getConfig()
			.getValue("emailVerificationServiceSid", String.class);
	public static final String emailAuthToken = ConfigProvider.getConfig().getValue("emailAuthToken", String.class);
	private static final String testEmailsStr = ConfigProvider.getConfig().getValue("testEmails", String.class);
	public static final List<String> testEmails = List.of(testEmailsStr.split(","));
	public static final String testSID = ConfigProvider.getConfig().getValue("testSID", String.class);
	public static final String testOTP = ConfigProvider.getConfig().getValue("testOTP", String.class);
	public static final Boolean logReqEnabled = ConfigProvider.getConfig().getValue("logReqEnabled", Boolean.class);
}
