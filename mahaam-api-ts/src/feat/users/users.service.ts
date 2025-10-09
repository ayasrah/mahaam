import { Injectable, Inject } from '@nestjs/common';
import { CreatedUser, VerifiedUser, Device, SuggestedEmail } from './users.model';
import { UsersRepo } from './users.repo';
import { DeviceRepo } from './device.repo';
import { SuggestedEmailRepo } from './suggested-email.repo';
import { PlansRepo } from '../plans/plans.repo';
import { Auth } from 'src/infra/auth';
import Email from 'src/infra/email';
import Log from 'src/infra/log';
import { Req } from 'src/infra/req';
import { InputError, NotFoundError, UnauthorizedError } from 'src/infra/errors';
import { DB } from 'src/infra/db';
import Config from 'src/infra/config';

export interface UsersService {
  create(device: Device): Promise<CreatedUser>;
  sendMeOtp(email: string): Promise<string | null>;
  verifyOtp(email: string, sid: string, otp: string): Promise<VerifiedUser>;
  refreshToken(): Promise<VerifiedUser>;
  updateName(name: string): Promise<void>;
  logout(deviceId: string): Promise<void>;
  delete(sid: string, otp: string): Promise<void>;
  getDevices(): Promise<Device[]>;
  getSuggestedEmails(): Promise<SuggestedEmail[]>;
  deleteSuggestedEmail(suggestedEmailId: string): Promise<void>;
}

@Injectable()
export class DefaultUsersService implements UsersService {
  constructor(
    @Inject('UsersRepo') private readonly usersRepo: UsersRepo,
    @Inject('DeviceRepo') private readonly deviceRepo: DeviceRepo,
    @Inject('SuggestedEmailRepo') private readonly suggestedEmailRepo: SuggestedEmailRepo,
    @Inject('PlansRepo') private readonly plansRepo: PlansRepo,
  ) {}

  async create(device: Device): Promise<CreatedUser> {
    return await DB.withTrx(async (trx) => {
      const userId = await this.usersRepo.create(trx);

      // Add device
      device.userId = userId;
      await this.deviceRepo.deleteByFingerprint(trx, device.fingerprint);
      const deviceId = await this.deviceRepo.create(trx, device);

      const jwt = Auth.createToken(userId, deviceId);

      Log.info(`User Created with id:${userId}, deviceId:${device.id}.`);
      return { id: userId, deviceId, jwt };
    });
  }

  async sendMeOtp(email: string): Promise<string | null> {
    const verifySid = Config.testEmails.includes(email) ? Config.testSID : await Email.sendOtp(email);
    if (verifySid) Log.info(`OTP sent to ${email}`);

    return verifySid;
  }

  async verifyOtp(email: string, sid: string, otp: string): Promise<VerifiedUser> {
    const isTest = Config.testEmails.includes(email) && sid === Config.testSID && otp === Config.testOTP;
    const otpStatus = isTest ? 'approved' : await Email.verifyOtp(otp, sid, email);
    if (otpStatus !== 'approved') throw new InputError(`OTP not verified for ${email}, status: ${otpStatus}`);

    return await DB.withTrx(async (trx) => {
      const user = await this.usersRepo.getOne(trx, email);

      if (!user) {
        await this.usersRepo.updateEmail(trx, Req.userId, email);
        Log.info(`User loggedIn for ${email}`);
      } else {
        // Move plans of current user to the one with email
        await this.plansRepo.updateUserId(trx, Req.userId, user.id);

        const devices = await this.deviceRepo.getMany(trx, user.id);
        if (devices && devices.length >= 5) {
          await this.deviceRepo.delete(trx, devices[devices.length - 1].id);
        }

        await this.deviceRepo.updateUserId(trx, Req.deviceId, user.id);
        await this.usersRepo.delete(trx, Req.userId);
        Log.info(`Merging userId:${Req.userId} to ${user.id}`);
      }

      const newUserId = user?.id || Req.userId;
      const jwt = Auth.createToken(newUserId, Req.deviceId);
      Log.info(`OTP verified for ${email}`);

      return {
        userId: newUserId,
        deviceId: Req.deviceId,
        jwt,
        userFullName: user?.name || null,
        email,
      };
    });
  }

  async refreshToken(): Promise<VerifiedUser> {
    const user = await DB.withTrx((trx) => this.usersRepo.getOneById(trx, Req.userId));

    const jwt = Auth.createToken(Req.userId, Req.deviceId);

    return {
      userId: Req.userId,
      deviceId: Req.deviceId,
      jwt,
      userFullName: user?.name || null,
      email: user?.email || null,
    };
  }

  async updateName(name: string): Promise<void> {
    await DB.withTrx((trx) => this.usersRepo.updateName(trx, Req.userId, name));
  }

  async logout(deviceId: string): Promise<void> {
    await DB.withTrx(async (trx) => {
      const device = await this.deviceRepo.getOne(trx, deviceId);
      if (!device || device.userId !== Req.userId) throw new UnauthorizedError('Invalid deviceId');
      await this.deviceRepo.delete(trx, deviceId);
    });
  }

  async deleteSuggestedEmail(suggestedEmailId: string): Promise<void> {
    await DB.withTrx(async (trx) => {
      const suggestedEmail = await this.suggestedEmailRepo.getOne(trx, suggestedEmailId);
      if (!suggestedEmail || suggestedEmail.userId !== Req.userId)
        throw new UnauthorizedError('Invalid suggestedEmailId');
      await this.suggestedEmailRepo.delete(trx, suggestedEmailId);
    });
  }

  async delete(sid: string, otp: string): Promise<void> {
    await DB.withTrx(async (trx) => {
      const user = await this.usersRepo.getOneById(trx, Req.userId);
      if (!user || !user.email) throw new NotFoundError('User or email not found');

      const isTest = Config.testEmails.includes(user.email) && sid === Config.testSID && otp === Config.testOTP;
      const otpStatus = isTest ? 'approved' : await Email.verifyOtp(otp, sid, user.email);
      if (otpStatus !== 'approved') throw new InputError(`OTP not approved for ${user.email}, status: ${otpStatus}`);

      await this.suggestedEmailRepo.deleteManyByEmail(trx, user.email!);
      await this.usersRepo.delete(trx, Req.userId);
    });
  }

  async getDevices(): Promise<Device[]> {
    return DB.withTrx((trx) => this.deviceRepo.getMany(trx, Req.userId));
  }

  async getSuggestedEmails(): Promise<SuggestedEmail[]> {
    return DB.withTrx((trx) => this.suggestedEmailRepo.getMany(trx, Req.userId));
  }
}
