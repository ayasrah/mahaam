import { Injectable } from '@nestjs/common';
import { Device } from './users.model';
import { DB, Trx } from '../../infra/db';
import { randomUUID } from 'crypto';

export interface DeviceRepo {
  create(trx: Trx, device: Device): Promise<string>;
  delete(trx: Trx, id: string): Promise<number>;
  deleteByUser(trx: Trx, userId: string, exceptDeviceId: string): Promise<number>;
  deleteByFingerprint(trx: Trx, fingerprint: string): Promise<number>;
  getOne(trx: Trx, deviceId: string): Promise<Device | null>;
  getMany(trx: Trx, userId: string): Promise<Device[]>;
  updateUserId(trx: Trx, deviceId: string, userId: string): Promise<number>;
}

@Injectable()
export class DefaultDeviceRepo implements DeviceRepo {
  async create(trx: Trx, device: Device): Promise<string> {
    const id = randomUUID();
    await trx`
      INSERT INTO devices (id, user_id, platform, fingerprint, info, created_at) 
      VALUES (${id}, ${device.userId}, ${device.platform || null}, ${device.fingerprint}, ${device.info || null}, current_timestamp)
    `;
    return id;
  }

  async deleteByFingerprint(trx: Trx, fingerprint: string): Promise<number> {
    const result = await trx`DELETE FROM devices WHERE fingerprint = ${fingerprint}`;
    return result.count || 0;
  }

  async deleteByUser(trx: Trx, userId: string, exceptDeviceId: string): Promise<number> {
    const result = await trx`DELETE FROM devices WHERE user_id = ${userId} AND id != ${exceptDeviceId}`;
    return result.count || 0;
  }

  async delete(trx: Trx, deviceId: string): Promise<number> {
    const result = await trx`DELETE FROM devices WHERE id = ${deviceId}`;
    return result.count || 0;
  }

  async updateUserId(trx: Trx, deviceId: string, userId: string): Promise<number> {
    const result = await trx`UPDATE devices SET user_id = ${userId}, updated_at = current_timestamp WHERE id = ${deviceId}`;
    return result.count || 0;
  }

  async getOne(trx: Trx, deviceId: string): Promise<Device | null> {
    const result =
      await trx`SELECT id, user_id, platform, fingerprint, info, created_at FROM devices WHERE id = ${deviceId} ORDER BY created_at DESC`;
    if (result.length === 0) return null;
    return DB.as<Device>(result[0]);
  }

  async getMany(trx: Trx, userId: string): Promise<Device[]> {
    const result =
      await trx`SELECT id, user_id, platform, fingerprint, info, created_at FROM devices WHERE user_id = ${userId} ORDER BY created_at DESC`;
    return result.map((row) => DB.as<Device>(row));
  }
}
