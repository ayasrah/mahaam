import { Injectable } from '@nestjs/common';
import { DB } from '../db';
import { Cache } from '../cache';

export interface LogRepo {
  create(trafficId: string | null, type: string, message: string): Promise<void>;
}

@Injectable()
export class DefaultLogRepo implements LogRepo {
  async create(trafficId: string | null, type: string, message: string): Promise<void> {
    setImmediate(async () => {
      try {
        await DB.sql`
		INSERT INTO x_log (traffic_id, type, message, node_ip, created_at) 
		VALUES (${trafficId || null}, ${type}, ${message}, ${Cache.getNodeIP()}, current_timestamp)
	  `;
      } catch (error) {
        console.error(error);
      }
    });
  }
}
