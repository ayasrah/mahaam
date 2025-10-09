import { Injectable } from '@nestjs/common';
import { Traffic } from './monitor.model';
import { DB } from '../db';

export interface TrafficRepo {
  create(traffic: Traffic): Promise<void>;
}

@Injectable()
export class DefaultTrafficRepo implements TrafficRepo {
  async create(traffic: Traffic): Promise<void> {
    setImmediate(async () => {
      try {
        await DB.sql`INSERT INTO x_traffic (id, health_id, method, path, code, elapsed, headers, request, response, created_at) 
	VALUES (${traffic.id}, ${traffic.healthId}, ${traffic.method}, ${traffic.path}, ${traffic.code || null}, ${traffic.elapsed || null}, 
	${traffic.headers || null}, ${traffic.request || null}, ${traffic.response || null}, current_timestamp)`;
      } catch (error) {
        console.error(error);
      }
    });
  }
}
