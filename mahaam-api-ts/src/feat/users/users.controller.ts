import { Controller, Get, Post, Delete, Patch, Inject, Param, Body, Res, ParseBoolPipe } from '@nestjs/common';
import { Response } from 'express';
import { UsersService } from './users.service';
import { Device } from './users.model';
import * as rule from '../../infra/rule';

export interface UsersController {
  sendMeOtp(email: string, res: Response): Promise<void>;
  create(
    platform: string,
    isPhysicalDevice: boolean,
    deviceFingerprint: string,
    deviceInfo: string,
    res: Response,
  ): Promise<void>;
  verifyOtp(email: string, sid: string, otp: string, res: Response): Promise<void>;
  refreshToken(res: Response): Promise<void>;
  updateName(name: string, res: Response): Promise<void>;
  logout(deviceId: string, res: Response): Promise<void>;
  delete(sid: string, otp: string, res: Response): Promise<void>;
  getDevices(res: Response): Promise<void>;
  getSuggestedEmails(res: Response): Promise<void>;
  deleteSuggestedEmail(suggestedEmailId: string, res: Response): Promise<void>;
}

@Controller('users')
export class DefaultUsersController implements UsersController {
  constructor(@Inject('UsersService') private readonly usersService: UsersService) {}

  @Post('send-me-otp')
  async sendMeOtp(@Body('email') email: string, @Res() res: Response) {
    // todo: add rate limit
    rule.validateEmail(email);
    const verificationSid = await this.usersService.sendMeOtp(email);
    res.status(200).json(verificationSid);
  }

  @Post('create')
  async create(
    @Body('platform') platform: string,
    @Body('isPhysicalDevice', ParseBoolPipe) isPhysicalDevice: boolean,
    @Body('deviceFingerprint') deviceFingerprint: string,
    @Body('deviceInfo') deviceInfo: string,
    @Res() res: Response,
  ) {
    rule.requiredBoolean(isPhysicalDevice, 'isPhysicalDevice');
    rule.required(platform, 'platform');
    rule.required(deviceFingerprint, 'deviceFingerprint');
    rule.required(deviceInfo, 'deviceInfo');
    rule.failIf(!isPhysicalDevice, 'Device should be real not simulator');

    const device = {
      platform,
      fingerprint: deviceFingerprint,
      info: deviceInfo,
    } as Device;
    const createdUser = await this.usersService.create(device);
    res.status(200).json(createdUser);
  }

  @Post('verify-otp')
  async verifyOtp(
    @Body('email') email: string,
    @Body('sid') sid: string,
    @Body('otp') otp: string,
    @Res() res: Response,
  ) {
    rule.required(email, 'email');
    rule.required(sid, 'sid');
    rule.required(otp, 'otp');

    const verifiedUser = await this.usersService.verifyOtp(email, sid, otp);
    res.status(200).json(verifiedUser);
  }

  @Post('refresh-token')
  async refreshToken(@Res() res: Response) {
    const verifiedUser = await this.usersService.refreshToken();
    res.status(200).json(verifiedUser);
  }

  @Patch('name')
  async updateName(@Body('name') name: string, @Res() res: Response) {
    rule.required(name, 'name');
    await this.usersService.updateName(name);
    res.sendStatus(200);
  }

  @Post('logout')
  async logout(@Body('deviceId') deviceId: string, @Res() res: Response) {
    rule.required(deviceId, 'deviceId');
    await this.usersService.logout(deviceId);
    res.sendStatus(200);
  }

  @Delete()
  async delete(@Body('sid') sid: string, @Body('otp') otp: string, @Res() res: Response) {
    rule.required(sid, 'sid');
    rule.required(otp, 'otp');
    await this.usersService.delete(sid, otp);
    res.sendStatus(204);
  }

  @Get('devices')
  async getDevices(@Res() res: Response) {
    const devices = await this.usersService.getDevices();
    res.status(200).json(devices);
  }

  @Get('suggested-emails')
  async getSuggestedEmails(@Res() res: Response) {
    const suggestedEmails = await this.usersService.getSuggestedEmails();
    res.status(200).json(suggestedEmails);
  }

  @Delete('suggested-emails')
  async deleteSuggestedEmail(@Body('suggestedEmailId') suggestedEmailId: string, @Res() res: Response) {
    rule.required(suggestedEmailId, 'suggestedEmailId');
    await this.usersService.deleteSuggestedEmail(suggestedEmailId);
    res.sendStatus(204);
  }
}
