import {
  Injectable,
  BadRequestException,
  NotFoundException,
} from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { CreateUrlDto } from './dto/create-url.dto';
import * as crypto from 'crypto';
import { Url } from './schemas/urls.schema';
import { UrlRepository } from './urls.repository';

@Injectable()
export class UrlService {
  private readonly shortCodeLength: number;
  private readonly maxRetries: number;

  constructor(
    private readonly urlRepository: UrlRepository,
    private readonly configService: ConfigService, 
  ) {
    this.shortCodeLength = 6;
    this.maxRetries = 5;
  }

  async createShortUrl(
    createUrlDto: CreateUrlDto,
    userId: string,
  ): Promise<Url> {
    try {
      const normalizedUrl = this.normalizeUrl(createUrlDto.originalUrl);

      if (createUrlDto.customAlias) {
        const aliasExists = await this.urlRepository.checkCustomAliasExists(
          createUrlDto.customAlias,
        );
        if (aliasExists) {
          throw new BadRequestException('Custom alias already exists');
        }
      }

      const shortCode =
        createUrlDto.customAlias || (await this.generateShortCode());
      const urlObj = new URL(normalizedUrl);
      const metadata = {
        domain: urlObj.hostname,
        protocol: urlObj.protocol,
        path: urlObj.pathname,
      };

      const urlData: Partial<Url> = {
        originalUrl: normalizedUrl,
        shortCode,
        customAlias: createUrlDto.customAlias,
        userId: new (require('mongoose').Types.ObjectId)(userId),
        title: createUrlDto.title,
        description: createUrlDto.description,
        expiresAt: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000),
        metadata,
        clickCount: 0,
        isActive: true,
      };

      const createdUrl = await this.urlRepository.create(urlData);
      return createdUrl;
    } catch (error) {
      throw error;
    }
  }

  private async generateShortCode(): Promise<string> {
    let attempts = 0;

    while (attempts < this.maxRetries) {
      let shortCode: string;

      switch (attempts) {
        case 0:
          shortCode = this.generateBase62ShortCode(this.shortCodeLength);
          break;
        case 1:
          shortCode = this.generateCryptoShortCode(this.shortCodeLength);
          break;
        case 2:
          shortCode = this.generateTimestampBasedCode(this.shortCodeLength);
          break;
        case 3:
          shortCode = this.generateUuidBasedCode(this.shortCodeLength);
          break;
        default:
          shortCode = this.generateBase62ShortCode(
            this.shortCodeLength + attempts - 3,
          );
          break;
      }

      const exists = await this.urlRepository.checkShortCodeExists(shortCode);

      if (!exists) {
        return shortCode;
      }

      attempts++;
    }

    throw new BadRequestException(
      'Failed to generate unique short code. Please try again.',
    );
  }

  // Strategy 1: Base62 encoding (0-9, a-z, A-Z)
  private generateBase62ShortCode(length: number): string {
    const chars =
      '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ';
    let result = '';

    for (let i = 0; i < length; i++) {
      const randomIndex = Math.floor(Math.random() * chars.length);
      result += chars[randomIndex];
    }

    return result;
  }

  // Strategy 2: Crypto-based random generation
  private generateCryptoShortCode(length: number): string {
    const chars =
      '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ';
    const randomBytes = crypto.randomBytes(length);
    let result = '';

    for (let i = 0; i < length; i++) {
      result += chars[randomBytes[i] % chars.length];
    }

    return result;
  }

  // Strategy 3: Timestamp-based with random suffix
  private generateTimestampBasedCode(length: number): string {
    const timestamp = Date.now().toString(36); // Convert to base36
    const randomSuffix = this.generateBase62ShortCode(
      Math.max(1, length - timestamp.length),
    );
    const combined = timestamp + randomSuffix;

    return combined.slice(-length); // Take last 'length' characters
  }

  // Strategy 4: UUID-based short code
  private generateUuidBasedCode(length: number): string {
    const uuid = crypto.randomUUID().replace(/-/g, '');
    const hash = crypto.createHash('sha256').update(uuid).digest('base64url');
    return hash.slice(0, length);
  }

  private normalizeUrl(url: string): string {
    try {
      const urlObj = new URL(url);

      let pathname = urlObj.pathname;
      if (pathname.endsWith('/') && pathname.length > 1) {
        pathname = pathname.slice(0, -1);
      }

      return `${urlObj.protocol}//${urlObj.host}${pathname}${urlObj.search}${urlObj.hash}`;
    } catch (error) {
      throw new BadRequestException('Invalid URL format');
    }
  }

  /**
   * Get original URL by short code
   */
  async getOriginalUrl(shortCode: string): Promise<Url> {
    const url = await this.urlRepository.findByShortCode(shortCode);

    if (!url) {
      throw new NotFoundException('Short URL not found or expired');
    }

    return url;
  }

  /**
   * Track click and redirect
   */
  async trackAndRedirect(shortCode: string): Promise<string> {
    const url = await this.getOriginalUrl(shortCode);

    return url.originalUrl;
  }

  /**
   * Get URL statistics
   */
  async getUrlStats(shortCode: string, userId?: string): Promise<any> {
    const url = await this.urlRepository.findByShortCode(shortCode);

    if (!url) {
      throw new NotFoundException('Short URL not found');
    }

    // Check ownership if userId provided
    if (userId && url.userId?.toString() !== userId) {
      throw new BadRequestException(
        'You do not have access to this URL statistics',
      );
    }

    // Calculate analytics
    // const analytics = url.analytics || [];
    // const clicksByDay = this.groupClicksByDay(analytics);
    // const topReferers = this.getTopReferers(analytics);
    // const topCountries = this.getTopCountries(analytics);

    return {
      url: {
        id: url.id,
        originalUrl: url.originalUrl,
        shortCode: url.shortCode,
        shortUrl: 'url.shortUrl',
        title: url.title,
        createdAt: url.createdAt,
      },
      statistics: {
        totalClicks: url.clickCount,
        // clicksByDay,
        // topReferers,
        // topCountries,
        // lastClickedAt: analytics[analytics.length - 1]?.clickedAt || null,
      },
      // recentClicks: analytics.slice(-10).reverse(),
    };
  }
}
