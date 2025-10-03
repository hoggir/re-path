import { ApiProperty } from '@nestjs/swagger';

export class CreateUrlResponseDto {
  @ApiProperty({ example: '507f1f77bcf86cd799439011' })
  id: string;

  @ApiProperty({ example: 'abcd1234' })
  shortCode: string;

  @ApiProperty({ example: 'https://www.example.com/some/long/url' })
  longUrl: string;

  @ApiProperty({ example: '2025-10-03T10:00:00.000Z' })
  activeUntil: Date;

  @ApiProperty({ example: true })
  isActive: boolean;

  @ApiProperty({ example: 666 })
  accessCount: number;
}
