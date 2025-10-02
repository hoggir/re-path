// src/common/exceptions/custom-exceptions.ts
import { HttpException, HttpStatus } from '@nestjs/common';

export class BadRequestException extends HttpException {
  constructor(message: string | string[] = 'Bad Request') {
    super(
      {
        error: 'Bad Request',
        message,
      },
      HttpStatus.BAD_REQUEST,
    );
  }
}

export class NotFoundException extends HttpException {
  constructor(message: string = 'Resource not found') {
    super(
      {
        error: 'Not Found',
        message,
      },
      HttpStatus.NOT_FOUND,
    );
  }
}

export class UnauthorizedException extends HttpException {
  constructor(message: string = 'Unauthorized') {
    super(
      {
        error: 'Unauthorized',
        message,
      },
      HttpStatus.UNAUTHORIZED,
    );
  }
}

export class ForbiddenException extends HttpException {
  constructor(message: string = 'Forbidden') {
    super(
      {
        error: 'Forbidden',
        message,
      },
      HttpStatus.FORBIDDEN,
    );
  }
}

export class ConflictException extends HttpException {
  constructor(message: string = 'Conflict') {
    super(
      {
        error: 'Conflict',
        message,
      },
      HttpStatus.CONFLICT,
    );
  }
}