import {
  Body,
  Controller,
  Get,
  HttpStatus,
  Post,
} from '@nestjs/common';
import {
  ApiBearerAuth,
  ApiBody,
  ApiOperation,
  ApiResponse,
  ApiTags,
} from '@nestjs/swagger';
import { Roles } from '../auth/decorators/roles.decorator';
import { ErrorResponseDto } from 'src/common/dto/error-response.dto';
import { CreateUrlDto } from './dto/create-url.dto';
import { UrlResponseDto } from './dto/urls-response.dto';
import { CurrentUser } from '../auth/decorators/current-user.decorator';
import { UrlService } from './urls.service';
import { ResponseMessage } from '../../common/decorators/response-message.decorator';

@ApiTags('URL Shortener')
@Controller('url')
@ApiBearerAuth('JWT-auth')
export class UrlController {
  constructor(private readonly urlService: UrlService) {}

  @Post('create')
  @Roles('user', 'admin')
  @ResponseMessage('URL shortened successfully')
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
    return this.urlService.createShortUrl(createUrlDto, userId);
  }

  @Get('metrics/collisions')
  @Roles('admin')
  @ResponseMessage('Collision metrics retrieved successfully')
  @ApiOperation({
    summary: 'Get short code collision metrics',
    description: 'Monitor collision rate for short code generation (admin only)',
  })
  @ApiResponse({
    status: HttpStatus.OK,
    description: 'Collision metrics retrieved successfully',
  })
  async getCollisionMetrics() {
    return this.urlService.getCollisionMetrics();
  }
}
