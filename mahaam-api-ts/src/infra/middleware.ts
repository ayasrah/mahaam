import {
  Injectable,
  NestMiddleware,
  HttpStatus,
  Inject,
  NestInterceptor,
  CallHandler,
  ExecutionContext,
} from '@nestjs/common';
import { Request, Response, NextFunction } from 'express';
import { randomUUID } from 'crypto';
import { Auth } from './auth';
import { UnauthorizedError, NotFoundError, AppError } from './errors';
import { Req } from './req';
import Log from './log';
import { Cache } from './cache';
import { Traffic } from './monitor/monitor.model';
import { TrafficRepo } from './monitor/traffic.repo';

import config from './config';
import { UsersRepo } from 'src/feat/users/users.repo';
import { DeviceRepo } from 'src/feat/users/device.repo';
import { Observable, throwError } from 'rxjs';
import { catchError, finalize, tap } from 'rxjs/operators';

@Injectable()
export class AppMiddleware implements NestMiddleware {
  constructor(
    @Inject('UsersRepo') private readonly usersRepo: UsersRepo,
    @Inject('DeviceRepo') private readonly deviceRepo: DeviceRepo,
  ) {}

  async use(req: Request, res: Response, next: NextFunction): Promise<void> {
    return Req.run(async () => {
      Req.startTime = Date.now();
      Req.trafficId = randomUUID();
      await this.authenticateReq(req);
      next();
    });
  }

  private async authenticateReq(req: Request): Promise<void> {
    let path = '';
    if (config.baseUrl && req.baseUrl.slice(1).startsWith(config.baseUrl)) {
      path = req.baseUrl.slice(config.baseUrl.length + 1);
    }

    if (config.baseUrl && !req.baseUrl.slice(1).startsWith(config.baseUrl)) {
      throw new NotFoundError('Invalid path base');
    }

    const appStore = req.headers['x-app-store'] as string;
    const appVersion = req.headers['x-app-version'] as string;

    if ((!appStore || !appVersion) && !path.startsWith('/swagger')) {
      throw new UnauthorizedError('Required headers not exists');
    }

    Req.appStore = appStore;
    Req.appVersion = appVersion;

    const bypassAuthPaths = ['/swagger', '/health', '/audit', '/users/create'];

    if (!bypassAuthPaths.some((bypassPath) => path.startsWith(bypassPath))) {
      const { userId, deviceId, isLoggedIn } = await Auth.validateAndExtractJwt(req, this.usersRepo, this.deviceRepo);
      Req.userId = userId;
      Req.deviceId = deviceId;
      Req.isLoggedIn = isLoggedIn;
    }
  }
}

@Injectable()
export class TrafficInterceptor implements NestInterceptor {
  constructor(@Inject('TrafficRepo') private readonly trafficRepo: TrafficRepo) {}

  intercept(context: ExecutionContext, next: CallHandler): Observable<any> {
    let request: Request;
    let response: Response;
    let responseBody: string;

    return next.handle().pipe(
      catchError((err: unknown) => {
        Log.error(err instanceof Error ? err.toString() : String(err));

        const ctx = context.switchToHttp();
        response = ctx.getResponse<Response>();

        responseBody = JSON.stringify(err instanceof Error ? err.message : String(err));
        let code = HttpStatus.INTERNAL_SERVER_ERROR;

        if (err instanceof AppError) {
          const appException = err as AppError;
          const key = appException.key;
          code = appException.httpCode;
          if (key) {
            const responseObj = { key, error: err.message };
            responseBody = JSON.stringify(responseObj);
          }
        }
        response.status(code).contentType('application/json');
        response.send(responseBody);

        return throwError(() => err);
      }),
      finalize(() => {
        const ctx = context.switchToHttp();
        request = ctx.getRequest<Request>();
        response = ctx.getResponse<Response>();
        this.createTraffic(request, response, responseBody);
      }),
    );
  }

  private createTraffic(req: Request, res: Response, responseBody: string | null): void {
    const path = req.path;
    const notTrafficPaths = path.startsWith('/swagger') || path === '/health' || path.startsWith('/audit');
    if (notTrafficPaths) return;

    const elapsed = Date.now() - Req.startTime;

    const headers = {
      userId: Req.userId,
      deviceId: Req.deviceId,
      appStore: Req.appStore,
      appVersion: Req.appVersion,
    };

    const isSuccessResponse = res.statusCode < 400;
    const isUserPath = path.startsWith('/user');
    let requestBody: string | null = null;
    if (isUserPath) {
      responseBody = null;
    }

    if (!isSuccessResponse && config.logReqEnabled) {
      requestBody = this.getReqBody(req);
    }

    const traffic: Traffic = {
      id: Req.trafficId || randomUUID(),
      method: req.method,
      path,
      code: res.statusCode,
      elapsed,
      headers: JSON.stringify(headers),
      request: requestBody,
      response: responseBody,
      healthId: Cache.getHealthId(),
    };

    this.trafficRepo.create(traffic);
  }

  private getReqBody(request: Request): string | null {
    if (request.body) {
      const body = typeof request.body === 'string' ? request.body : JSON.stringify(request.body);
      return body || null;
    }
    return null;
  }
}
