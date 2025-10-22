import { Controller, Get, HttpStatus } from '@nestjs/common';
import {
  ApiTags,
  ApiOperation,
  ApiResponse,
  ApiBearerAuth,
} from '@nestjs/swagger';
import { Roles } from '../auth/decorators/roles.decorator';
import { CurrentUser } from '../auth/decorators/current-user.decorator';
import { Public } from '../auth/decorators/public.decorator';
import { ErrorResponseDto } from '../../common/dto/error-response.dto';
import { ResponseMessage } from '../../common/decorators/response-message.decorator';

@ApiTags('Users')
@Controller('users')
@ApiBearerAuth('JWT-auth')
export class UsersController {
  @Public()
  @Get('public')
  @ResponseMessage('Public endpoint accessed')
  @ApiOperation({
    summary: 'Public endpoint',
    description: 'This endpoint is accessible without authentication',
  })
  @ApiResponse({
    status: HttpStatus.OK,
    description: 'Success response',
    schema: {
      properties: {
        success: { type: 'boolean', example: true },
        statusCode: { type: 'number', example: 200 },
        message: { type: 'string', example: 'Success' },
        data: {
          type: 'object',
          properties: {
            message: { type: 'string', example: 'This is a public route' },
          },
        },
        timestamp: { type: 'string', example: '2025-10-03T10:00:00.000Z' },
      },
    },
  })
  publicRoute() {
    return { message: 'This is a public route' };
  }

  @Get('profile')
  @ResponseMessage('Profile retrieved')
  @ApiOperation({
    summary: 'Get user profile',
    description: 'Retrieve the profile of the authenticated user',
  })
  @ApiResponse({
    status: HttpStatus.OK,
    description: 'Profile retrieved successfully',
    schema: {
      properties: {
        success: { type: 'boolean', example: true },
        statusCode: { type: 'number', example: 200 },
        message: { type: 'string', example: 'Profile retrieved' },
        data: {
          type: 'object',
          properties: {
            userId: { type: 'string', example: '507f1f77bcf86cd799439011' },
            email: { type: 'string', example: 'john.doe@example.com' },
            role: { type: 'string', example: 'user' },
          },
        },
        timestamp: { type: 'string', example: '2025-10-03T10:00:00.000Z' },
      },
    },
  })
  @ApiResponse({
    status: HttpStatus.UNAUTHORIZED,
    description: 'Invalid or expired token',
    type: ErrorResponseDto,
  })
  getProfile(@CurrentUser() user: any) {
    return {
      userId: user.encryptedUserId,
      email: user.email,
      role: user.role,
    };
  }

  @Get('admin')
  @Roles('admin')
  @ResponseMessage('Admin access granted')
  @ApiOperation({
    summary: 'Admin only endpoint',
    description: 'This endpoint is only accessible by users with admin role',
  })
  @ApiResponse({
    status: HttpStatus.OK,
    description: 'Success response for admin',
    schema: {
      properties: {
        success: { type: 'boolean', example: true },
        statusCode: { type: 'number', example: 200 },
        message: { type: 'string', example: 'Admin access' },
        data: {
          type: 'object',
          properties: {
            userId: { type: 'string', example: '507f1f77bcf86cd799439011' },
            email: { type: 'string', example: 'admin@example.com' },
            role: { type: 'string', example: 'admin' },
          },
        },
        timestamp: { type: 'string', example: '2025-10-03T10:00:00.000Z' },
      },
    },
  })
  @ApiResponse({
    status: HttpStatus.UNAUTHORIZED,
    description: 'Invalid or expired token',
    type: ErrorResponseDto,
  })
  @ApiResponse({
    status: HttpStatus.FORBIDDEN,
    description: 'Insufficient permissions',
    type: ErrorResponseDto,
  })
  adminOnly(@CurrentUser() user: any) {
    return {
      userId: user.encryptedUserId,
      email: user.email,
      role: user.role,
    };
  }

  @Get('moderator')
  @Roles('admin', 'moderator')
  @ResponseMessage('Moderator access granted')
  @ApiOperation({
    summary: 'Moderator or Admin endpoint',
    description:
      'This endpoint is accessible by users with moderator or admin role',
  })
  @ApiResponse({
    status: HttpStatus.OK,
    description: 'Success response',
  })
  @ApiResponse({
    status: HttpStatus.UNAUTHORIZED,
    description: 'Invalid or expired token',
    type: ErrorResponseDto,
  })
  @ApiResponse({
    status: HttpStatus.FORBIDDEN,
    description: 'Insufficient permissions',
    type: ErrorResponseDto,
  })
  moderatorAccess() {
    return { access: 'Moderator or Admin' };
  }
}
