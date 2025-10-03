import { ApiProperty } from '@nestjs/swagger';
import { IsString, MinLength, MaxLength, IsNotEmpty } from 'class-validator';

export class CreateUrlDto {
  @ApiProperty({
    description: 'Long URL to be shortened',
    example: 'https://www.example.com/some/long/url',
  })
  @IsNotEmpty({ message: 'Long URL is required' })
  @IsString()
  longUrl: string;
}
