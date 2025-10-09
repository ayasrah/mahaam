import config from './config';
import * as postgres from 'postgres';
import camelcaseKeys from 'camelcase-keys';
import Log from './log';

export type Trx = postgres.TransactionSql;

export class DB {
  public static sql: postgres.Sql;

  public static async init(): Promise<void> {
    this.sql = postgres(config.dbUrl);
  }

  public static as<T>(result: any): T {
    return camelcaseKeys(result, { deep: true }) as T;
  }

  public static withTrx<T>(callback: (trx: Trx) => Promise<T>) {
    return this.sql.begin(callback);
  }
}
