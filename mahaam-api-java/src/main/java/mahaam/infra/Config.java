package mahaam.infra;

import java.util.List;

import io.quarkus.runtime.annotations.StaticInitSafe;
import io.smallrye.config.ConfigMapping;

@StaticInitSafe
@ConfigMapping(prefix = "mahaam", namingStrategy = ConfigMapping.NamingStrategy.VERBATIM)
public interface Config {
	String apiName();

	String apiVersion();

	String envName();

	String dbUrl();

	String tokenSecretKey();

	String emailAccountSid();

	String emailVerificationServiceSid();

	String emailAuthToken();

	List<String> testEmails();

	String testSID();

	String testOTP();

	Boolean logReqEnabled();
}