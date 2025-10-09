import { Injectable, Inject } from '@nestjs/common';
import { SuggestedEmail } from './users.model';
import { DB, Trx } from '../../infra/db';
import Log from '../../infra/log';
import { randomUUID } from 'crypto';

export interface SuggestedEmailRepo {
  create(trx: Trx, userId: string, email: string): Promise<string>;
  delete(trx: Trx, id: string): Promise<number>;
  getMany(trx: Trx, userId: string): Promise<SuggestedEmail[]>;
  getOne(trx: Trx, id: string): Promise<SuggestedEmail | null>;
  deleteManyByEmail(trx: Trx, email: string): Promise<number>;
}

@Injectable()
export class DefaultSuggestedEmailRepo implements SuggestedEmailRepo {
  async create(trx: Trx, userId: string, email: string): Promise<string> {
    const id = randomUUID();
    const result = await trx`INSERT INTO suggested_emails (id, user_id, email, created_at) 
        VALUES (${id}, ${userId}, ${email}, current_timestamp)
		ON CONFLICT (user_id, email) DO NOTHING`;
    return result.count > 0 ? id : '';
  }

  async delete(trx: Trx, id: string): Promise<number> {
    const result = await trx`DELETE FROM suggested_emails WHERE id = ${id}`;
    return result.count || 0;
  }

  async deleteManyByEmail(trx: Trx, email: string): Promise<number> {
    const result = await trx`DELETE FROM suggested_emails WHERE email = ${email}`;
    return result.count || 0;
  }

  async getMany(trx: Trx, userId: string): Promise<SuggestedEmail[]> {
    const result =
      await trx`SELECT id, user_id, email, created_at FROM suggested_emails WHERE user_id = ${userId} ORDER BY created_at DESC`;
    return result.map((row) => DB.as<SuggestedEmail>(row));
  }

  async getOne(trx: Trx, id: string): Promise<SuggestedEmail | null> {
    const result = await trx`SELECT id, user_id, email, created_at FROM suggested_emails WHERE id = ${id}`;
    if (result.length === 0) return null;
    return DB.as<SuggestedEmail>(result[0]);
  }
}
