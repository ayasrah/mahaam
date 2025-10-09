import { HttpException, HttpStatus } from '@nestjs/common';

export abstract class AppError extends HttpException {
  public readonly key?: string;
  public readonly httpCode: number;

  constructor(message: string, httpCode: number, key?: string) {
    super(message, httpCode);
    this.httpCode = httpCode;
    this.key = key;
  }
}

export class UnauthorizedError extends AppError {
  constructor(message: string) {
    super(message, HttpStatus.UNAUTHORIZED);
  }
}

export class ForbiddenError extends AppError {
  constructor(message: string) {
    super(message, HttpStatus.FORBIDDEN);
  }
}

export class LogicError extends AppError {
  constructor(message: string, key?: string) {
    super(message, HttpStatus.CONFLICT, key);
  }
}

export class NotFoundError extends AppError {
  constructor(message: string) {
    super(message, HttpStatus.NOT_FOUND);
  }
}

export class InputError extends AppError {
  constructor(message: string) {
    super(message, HttpStatus.BAD_REQUEST);
  }
}
