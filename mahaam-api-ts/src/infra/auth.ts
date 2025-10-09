import * as jwt from 'jsonwebtoken';
import { Request } from 'express';
import { UnauthorizedError, ForbiddenError } from './errors';
import * as Rule from './rule';
import Config from './config';
import { UsersRepo } from 'src/feat/users/users.repo';
import { UUID } from 'crypto';
import { DB } from './db';
import { DeviceRepo } from 'src/feat/users/device.repo';

export interface AuthResult {
  userId: string;
  deviceId: string;
  isLoggedIn: boolean;
}

export class Auth {
  public static async validateAndExtractJwt(
    request: Request,
    userRepo: UsersRepo,
    deviceRepo: DeviceRepo,
  ): Promise<AuthResult> {
    const path = request.path;
    const authorization = request.headers.authorization;

    if (!authorization) {
      throw new UnauthorizedError('Authorization header not exists');
    }

    if (!authorization.startsWith('Bearer ')) {
      throw new UnauthorizedError('Invalid Authorization header format');
    }

    // Remove 'Bearer ' to get the jwt token
    const tokenString = authorization.substring(7);

    JWT.validate(tokenString);

    try {
      const decoded = jwt.decode(tokenString) as any;

      if (!decoded || typeof decoded === 'string') {
        throw new UnauthorizedError('Invalid token format');
      }

      const userId: UUID = decoded.userId as UUID;
      Auth.nonEmptyUuid(userId, 'userId');

      const deviceId = decoded.deviceId;
      Auth.nonEmptyUuid(deviceId, 'deviceId');

      await DB.withTrx(async (trx) => {
        const device = await deviceRepo.getOne(trx, deviceId);
        if ((device === null || device.userId !== userId) && path !== '/user/logout') {
          throw new UnauthorizedError('Invalid user info');
        }
      });

      let isLoggedIn = false;
      await DB.withTrx(async (trx) => {
        const user = await userRepo.getOneById(trx, userId);
        if (!user) throw new UnauthorizedError('Invalid user info');
        isLoggedIn = !!user.email;
      });

      return {
        userId,
        deviceId,
        isLoggedIn,
      };
    } catch (error) {
      if (error instanceof UnauthorizedError || error instanceof ForbiddenError) {
        throw error;
      }
      throw new UnauthorizedError(`Invalid JWT token: ${error}`);
    }
  }

  public static createToken(userId: string, deviceId: string): string {
    return JWT.create(userId, deviceId);
  }

  private static isValidUUID(uuid: string): boolean {
    const uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;
    return uuidRegex.test(uuid);
  }

  private static EMPTY_GUID = '00000000-0000-0000-0000-000000000000';

  public static nonEmptyUuid(uuidString: string | null | undefined, name: string): void {
    Rule.required(uuidString, name);
    if (!uuidString || uuidString.trim() === '' || uuidString === Auth.EMPTY_GUID) {
      throw new ForbiddenError(`${name} is Empty`);
    }
    if (!this.isValidUUID(uuidString)) {
      throw new UnauthorizedError(`${name} is not a valid UUID`);
    }
  }
}

class JWT {
  /**
   * Creates a JWT token with specified claims
   * @param userId User identifier
   * @param deviceId Device identifier
   * @returns JWT token string
   */
  public static create(userId: string, deviceId: string): string {
    try {
      const payload = {
        userId,
        deviceId,
      };

      const options: jwt.SignOptions = {
        expiresIn: '7d',
        issuer: 'mahaam-api',
        algorithm: 'HS256',
      };

      return jwt.sign(payload, this.getSecurityKey(), options);
    } catch (exception) {
      throw new Error(`JWT creation failed: ${exception}`);
    }
  }

  /**
   * Validates a JWT token
   * @param token JWT token to validate
   */
  public static validate(token: string): void {
    try {
      const validationParams = this.getValidationParams();
      jwt.verify(token, this.getSecurityKey(), validationParams);
    } catch (error) {
      throw new UnauthorizedError(`Invalid token: ${error}`);
    }
  }

  /**
   * Gets validation parameters for JWT verification
   * @returns JWT verification options
   */
  private static getValidationParams(): jwt.VerifyOptions {
    return {
      algorithms: ['HS256'],
      issuer: 'mahaam-api',
      // audience: 'Sample', // Commented out as not used in original
      clockTolerance: 30, // 30 seconds tolerance for clock skew
    };
  }

  /**
   * Gets the security key for signing/verifying tokens
   * @returns Secret key as string
   */
  private static getSecurityKey(): string {
    const secretKey = Config.tokenSecretKey;
    if (!secretKey) {
      throw new Error('Token secret key is not configured');
    }
    return secretKey;
  }
}
