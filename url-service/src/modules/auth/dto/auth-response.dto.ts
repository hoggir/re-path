import { ApiProperty } from '@nestjs/swagger';

export class UserResponseDto {
  @ApiProperty({ example: '507f1f77bcf86cd799439011' })
  id: string;

  @ApiProperty({ example: 'John Doe' })
  name: string;

  @ApiProperty({ example: 'john.doe@example.com' })
  email: string;

  @ApiProperty({ example: 'user', enum: ['user', 'admin', 'moderator'] })
  role: string;

  @ApiProperty({ example: true })
  isActive: boolean;

  @ApiProperty({ example: '2025-10-03T10:00:00.000Z', required: false })
  lastLoginAt?: Date;
}

export class TokensDto {
  @ApiProperty({
    example: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...',
    description: 'JWT access token valid for 15 minutes',
  })
  accessToken: string;

  @ApiProperty({
    example: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...',
    description: 'JWT refresh token valid for 7 days',
  })
  refreshToken: string;

  @ApiProperty({ example: 'Bearer' })
  tokenType: string;
}

export class AuthResponseDataDto {
  @ApiProperty({ type: UserResponseDto })
  user: UserResponseDto;

  @ApiProperty()
  accessToken: string;

  @ApiProperty()
  refreshToken: string;

  @ApiProperty()
  tokenType: string;
}

export class AuthResponseDto {
  @ApiProperty({ example: true })
  success: boolean;

  @ApiProperty({ example: 200 })
  statusCode: number;

  @ApiProperty({ example: 'Login successful' })
  message: string;

  @ApiProperty({ type: AuthResponseDataDto })
  data: AuthResponseDataDto;

  @ApiProperty({ example: '2025-10-03T10:00:00.000Z' })
  timestamp: string;
}

export class RefreshTokenResponseDto {
  @ApiProperty({ example: true })
  success: boolean;

  @ApiProperty({ example: 200 })
  statusCode: number;

  @ApiProperty({ example: 'Token refreshed successfully' })
  message: string;

  @ApiProperty({ type: TokensDto })
  data: TokensDto;

  @ApiProperty({ example: '2025-10-03T10:00:00.000Z' })
  timestamp: string;
}

export class LogoutResponseDto {
  @ApiProperty({ example: true })
  success: boolean;

  @ApiProperty({ example: 200 })
  statusCode: number;

  @ApiProperty({ example: 'Logout successful' })
  message: string;

  @ApiProperty({ example: '2025-10-03T10:00:00.000Z' })
  timestamp: string;
}
