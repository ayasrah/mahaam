import { Injectable } from '@nestjs/common';
import { User } from './users.model';
import { DB, Trx } from '../../infra/db';
import { randomUUID } from 'crypto';

export interface UsersRepo {
  create(trx: Trx): Promise<string>;
  updateName(trx: Trx, id: string, name: string): Promise<void>;
  updateEmail(trx: Trx, id: string, email: string): Promise<void>;
  getOne(trx: Trx, email: string): Promise<User | null>;
  getOneById(trx: Trx, id: string): Promise<User | null>;
  delete(trx: Trx, id: string): Promise<number>;
}

@Injectable()
export class DefaultUsersRepo implements UsersRepo {
  async create(trx: Trx): Promise<string> {
    const id = randomUUID();
    await trx`INSERT INTO users (id, created_at) VALUES (${id}, current_timestamp)`;
    return id;
  }

  async updateName(trx: Trx, id: string, name: string): Promise<void> {
    await trx`UPDATE users SET name = ${name}, updated_at = current_timestamp WHERE id = ${id}`;
  }

  async updateEmail(trx: Trx, id: string, email: string): Promise<void> {
    await trx`UPDATE users SET email = ${email}, updated_at = current_timestamp WHERE id = ${id}`;
  }

  async getOne(trx: Trx, email: string): Promise<User | null> {
    const result = await trx`SELECT id, name, email FROM users WHERE email = ${email}`;
    if (result.length === 0) return null;
    return DB.as<User>(result[0]);
  }

  async getOneById(trx: Trx, id: string): Promise<User | null> {
    const result = await trx`SELECT id, name, email FROM users WHERE id = ${id}`;
    if (result.length === 0) return null;
    return DB.as<User>(result[0]);
  }

  async delete(trx: Trx, id: string): Promise<number> {
    const result = await trx`DELETE FROM users WHERE id = ${id}`;
    return result.count || 0;
  }
}
