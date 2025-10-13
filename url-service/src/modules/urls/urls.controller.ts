import {
  Body,
  Controller,
  HttpCode,
  HttpStatus,
  Post,
  UseGuards,
} from '@nestjs/common';
import {
  ApiBearerAuth,
  ApiBody,
  ApiOperation,
  ApiResponse,
  ApiTags,
} from '@nestjs/swagger';
import { JwtAuthGuard } from '../auth/guards/jwt-auth.guard';
import { Roles } from '../auth/decorators/roles.decorator';
import { ErrorResponseDto } from 'src/common/dto/error-response.dto';
import { CreateUrlDto } from './dto/create-url.dto';
import { UrlResponseDto } from './dto/urls-response.dto';
import { CurrentUser } from '../auth/decorators/current-user.decorator';
import { UrlService } from './urls.service';

@ApiTags('URL Shortener')
@Controller('url')
@ApiBearerAuth('JWT-auth')
export class UrlController {
  constructor(private readonly urlService: UrlService) {}

  @Roles('user')
  @Post('create')
  @UseGuards(JwtAuthGuard)
  @ApiBody({ type: CreateUrlDto })
  @ApiBearerAuth('JWT-auth')
  @ApiOperation({
    summary: 'Create shortened URL',
    description: 'Create a new shortened URL from a long URL',
  })
  @ApiResponse({
    status: HttpStatus.CREATED,
    description: 'URL shortened successfully',
    type: UrlResponseDto,
  })
  @ApiResponse({
    status: HttpStatus.BAD_REQUEST,
    description: 'Invalid URL or custom alias already exists',
    type: ErrorResponseDto,
  })
  async createShortUrl(
    @Body() createUrlDto: CreateUrlDto,
    @CurrentUser('userId') userId: number,
  ) {
    const url = await this.urlService.createShortUrl(createUrlDto, userId);

    return {
      message: 'URL shortened successfully',
      data: url,
    };
  }
}
