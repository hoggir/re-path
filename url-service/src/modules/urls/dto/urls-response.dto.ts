// src/modules/url-shortener/dto/url-response.dto.ts
import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';

export class UrlResponseDto {
  @ApiProperty({ example: '507f1f77bcf86cd799439011' })
  id: string;

  @ApiProperty({ example: 'https://www.example.com/very/long/url/path' })
  originalUrl: string;

  @ApiProperty({ example: 'abc123' })
  shortCode: string;

  @ApiProperty({ example: 'http://localhost:3010/abc123' })
  shortUrl: string;

  @ApiPropertyOptional({ example: 'my-custom-link' })
  customAlias?: string;

  @ApiProperty({ example: 0 })
  clickCount: number;

  @ApiProperty({ example: true })
  isActive: boolean;

  @ApiPropertyOptional({ example: 'My Awesome Link' })
  title?: string;

  @ApiPropertyOptional({ example: 'This is a description' })
  description?: string;

  @ApiPropertyOptional({ example: '2025-12-31T23:59:59.000Z' })
  expiresAt?: Date;

  @ApiProperty({ example: '2025-10-03T10:00:00.000Z' })
  createdAt: Date;

  @ApiProperty({ example: '2025-10-03T10:00:00.000Z' })
  updatedAt: Date;
}
