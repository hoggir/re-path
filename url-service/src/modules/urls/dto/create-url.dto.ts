// src/modules/url-shortener/dto/create-url.dto.ts
import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';
import {
  IsUrl,
  IsString,
  IsOptional,
  MinLength,
  MaxLength,
  Matches,
} from 'class-validator';

export class CreateUrlDto {
  @ApiProperty({
    description: 'The original URL to be shortened',
    example: 'https://www.example.com/very/long/url/path',
  })
  @IsUrl({}, { message: 'Please provide a valid URL' })
  originalUrl: string;

  @ApiPropertyOptional({
    description:
      'Custom alias for the short URL (optional, 3-20 characters, alphanumeric and hyphens only)',
    example: 'my-custom-link',
    minLength: 3,
    maxLength: 20,
  })
  @IsOptional()
  @IsString()
  @MinLength(3, { message: 'Custom alias must be at least 3 characters' })
  @MaxLength(20, { message: 'Custom alias must not exceed 20 characters' })
  @Matches(/^[a-zA-Z0-9-_]+$/, {
    message:
      'Custom alias can only contain letters, numbers, hyphens, and underscores',
  })
  customAlias?: string;

  @ApiPropertyOptional({
    description: 'Title for the shortened URL',
    example: 'My Awesome Link',
  })
  @IsOptional()
  @IsString()
  @MaxLength(100)
  title?: string;

  @ApiPropertyOptional({
    description: 'Description for the shortened URL',
    example: 'This is a description of my link',
  })
  @IsOptional()
  @IsString()
  @MaxLength(500)
  description?: string;
}
