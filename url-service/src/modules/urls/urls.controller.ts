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
import {
  ErrorResponseDto,
  ValidationErrorResponseDto,
} from 'src/common/dto/error-response.dto';
import { CreateUrlDto } from './dto/create-url.dto';
import { CreateUrlResponseDto } from './dto/urls-response.dto';
import { CurrentUser } from '../auth/decorators/current-user.decorator';
import { UrlService } from './urls.service';

@ApiTags('Urls')
@Controller('urls')
@ApiBearerAuth('JWT-auth')
export class UrlsController {
    constructor(private readonly urlService: UrlService) {}
  
  @Roles('user')
  @UseGuards(JwtAuthGuard)
  @Post('create')
  @HttpCode(HttpStatus.CREATED)
  @ApiBody({ type: CreateUrlDto })
  @ApiOperation({
    summary: 'Create Short URL',
    description: 'Create a short URL for user.',
  })
  @ApiResponse({
    status: HttpStatus.UNAUTHORIZED,
    description: 'Invalid or expired token',
    type: ErrorResponseDto,
  })
  @ApiResponse({
    status: HttpStatus.CREATED,
    description: 'Short URL created successfully',
    type: CreateUrlResponseDto,
  })
  @ApiResponse({
    status: HttpStatus.BAD_REQUEST,
    description: 'Validation error',
    type: ValidationErrorResponseDto,
  })
  async createShortUrl(
    @CurrentUser() user: { id: string },
    @Body() createUrlDto: CreateUrlDto,
  ) {
    const payload = { idUser: user.id, longUrl: createUrlDto.longUrl };
    const shortUrl = await this.urlService.createShortUrl(payload);
    return {
      message: 'Short URL created successfully',
      data: 'ASD',
    };
  }
}
