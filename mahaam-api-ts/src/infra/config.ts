/* eslint-disable prettier/prettier */

/**
 * Configuration module that reads from environment variables and .env file.
 *
 * Configuration priority (highest to lowest):
 * 1. Environment variables
 * 2. .env file (in project root)
 * 3. Default values
 *
 * To set up:
 * 1. Copy env.example to .env
 * 2. Update .env with your actual values
 * 3. The .env file is gitignored for security
 */

const configObject = {
  get apiName(): string {
    return process.env.apiName || '';
  },

  get apiVersion(): string {
    return process.env.apiVersion || '';
  },

  get envName(): string {
    return process.env.envName || 'development';
  },

  get dbUrl(): string {
    return process.env.dbUrl || '';
  },

  get logFile(): string {
    return process.env.logFile || '';
  },

  get httpPort(): number {
    return parseInt(process.env.httpPort!, 10);
  },

  get baseUrl(): string {
    return process.env.baseUrl || '';
  },

  get tokenSecretKey(): string {
    return process.env.tokenSecretKey || '';
  },

  get emailAccountSid(): string {
    return process.env.emailAccountSid || '';
  },

  get emailVerificationServiceSid(): string {
    return process.env.EMAIL_VERIFICATION_SERVICE_SID || process.env.emailVerificationServiceSid || '';
  },

  get emailAuthToken(): string {
    return process.env.emailAuthToken || '';
  },

  get testEmails(): string[] {
    const testEmailsStr = process.env.testEmails || '';
    return testEmailsStr ? testEmailsStr.split(',').map((email) => email.trim()) : [];
  },

  get testSID(): string {
    return process.env.testSID || '';
  },

  get testOTP(): string {
    return process.env.testOTP || '';
  },

  get logReqEnabled(): boolean {
    return process.env.logReqEnabled === 'true';
  },
} as const;

export default configObject;
