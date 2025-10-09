/* eslint-disable */

import * as twilio from 'twilio';
import config from './config';

const client = twilio(config.emailAccountSid, config.emailAuthToken);

//    function init()  {
//     try {

//       this.verificationServiceSid = Config.emailVerificationServiceSid;
//     } catch (error: unknown) {
//       console.error('Failed to initialize Twilio client:', error);
//       throw error;
//     }
//   }

async function sendOtp(email: string): Promise<string | null> {
  try {
    if (!client || !config.emailVerificationServiceSid) {
      throw new Error('Email service not initialized. Call Email.init() first.');
    }

    const verification = await client.verify.v2.services(config.emailVerificationServiceSid).verifications.create({
      to: email,
      channel: 'email',
    });

    return verification.sid;
  } catch (error: unknown) {
    console.error('Failed to send OTP:', error);
    return null;
  }
}

async function verifyOtp(otp: string, sid: string, email: string): Promise<string | null> {
  try {
    if (!client || !config.emailVerificationServiceSid) {
      throw new Error('Email service not initialized. Call Email.init() first.');
    }

    const verificationCheck = await client.verify.v2.services(config.emailVerificationServiceSid).verificationChecks.create({
      to: email,
      code: otp,
    });

    return verificationCheck.status;
  } catch (error: unknown) {
    console.error('Failed to verify OTP:', error);
    return null;
  }
}

export default {
  sendOtp,
  verifyOtp,
} as const;
