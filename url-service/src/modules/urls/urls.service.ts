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
  private readonly baseRetryDelayMs: number;
  private collisionCount: number = 0;

  constructor(
    private readonly urlRepository: UrlRepository,
    private readonly configService: ConfigService,
  ) {
    this.shortCodeLength = 6;
    this.maxRetries = 10;
    this.baseRetryDelayMs = 10;
  }

  async createShortUrl(
    createUrlDto: CreateUrlDto,
    userId: number,
  ): Promise<Url> {
    const normalizedUrl = this.normalizeUrl(createUrlDto.originalUrl);

    if (createUrlDto.customAlias) {
      return this.createWithCustomAlias(createUrlDto, normalizedUrl, userId);
    }

    return this.createWithGeneratedShortCode(
      createUrlDto,
      normalizedUrl,
      userId,
    );
  }

  private async createWithCustomAlias(
    createUrlDto: CreateUrlDto,
    normalizedUrl: string,
    userId: number,
  ): Promise<Url> {
    const urlObj = new URL(normalizedUrl);
    const urlData: Partial<Url> = {
      originalUrl: normalizedUrl,
      shortCode: createUrlDto.customAlias!,
      customAlias: createUrlDto.customAlias,
      userId: userId,
      title: createUrlDto.title,
      description: createUrlDto.description,
      expiresAt: new Date(
        Date.now() +
          this.configService.get('url.defaultTtlDays', 7) * 24 * 60 * 60 * 1000,
      ),
      metadata: {
        domain: urlObj.hostname,
        protocol: urlObj.protocol,
        path: urlObj.pathname,
      },
      clickCount: 0,
      isActive: true,
    };

    try {
      const createdUrl = await this.urlRepository.create(urlData);
      return createdUrl;
    } catch (error: any) {
      if (error.code === 11000 || error.message?.includes('duplicate')) {
        throw new BadRequestException('Custom alias already exists');
      }
      throw error;
    }
  }

  private async createWithGeneratedShortCode(
    createUrlDto: CreateUrlDto,
    normalizedUrl: string,
    userId: number,
  ): Promise<Url> {
    const urlObj = new URL(normalizedUrl);
    const baseUrlData: Partial<Url> = {
      originalUrl: normalizedUrl,
      customAlias: createUrlDto.customAlias,
      userId: userId,
      title: createUrlDto.title,
      description: createUrlDto.description,
      expiresAt: new Date(
        Date.now() +
          this.configService.get('url.defaultTtlDays', 7) * 24 * 60 * 60 * 1000,
      ),
      metadata: {
        domain: urlObj.hostname,
        protocol: urlObj.protocol,
        path: urlObj.pathname,
      },
      clickCount: 0,
      isActive: true,
    };

    let attempts = 0;
    let currentLength = this.shortCodeLength;

    while (attempts < this.maxRetries) {
      const shortCode = this.generateShortCodeByStrategy(
        attempts,
        currentLength,
      );

      try {
        const createdUrl = await this.urlRepository.create({
          ...baseUrlData,
          shortCode,
        });

        if (attempts > 0) {
          console.warn(
            `✅ Short code generated after ${attempts + 1} attempts (collision detected)`,
          );
          this.collisionCount += attempts;
        }

        return createdUrl;
      } catch (error: any) {
        const isDuplicateError =
          error.code === 11000 || error.message?.includes('duplicate');

        if (!isDuplicateError) {
          throw error;
        }

        attempts++;

        if (attempts % 3 === 0) {
          currentLength++;
        }

        if (attempts <= 3) {
          console.warn(
            `⚠️  Short code collision detected (attempt ${attempts}/${this.maxRetries})`,
          );
        }

        if (attempts < this.maxRetries) {
          await this.exponentialBackoff(attempts);
        }
      }
    }

    console.error(
      `❌ Failed to generate unique short code after ${this.maxRetries} attempts`,
    );
    throw new BadRequestException(
      'Unable to generate unique short code. Please try again later.',
    );
  }

  private generateShortCodeByStrategy(attempt: number, length: number): string {
    switch (attempt % 4) {
      case 0:
        return this.generateCryptoShortCode(length);
      case 1:
        return this.generateUuidBasedCode(length);
      case 2:
        return this.generateTimestampBasedCode(length);
      case 3:
        return this.generateCryptoShortCode(length); // Use crypto again (most secure)
      default:
        return this.generateCryptoShortCode(length);
    }
  }

  private async exponentialBackoff(attempt: number): Promise<void> {
    const delay = this.baseRetryDelayMs * Math.pow(2, attempt);

    const jitter = Math.random() * delay * 0.5;
    const totalDelay = delay + jitter;

    const cappedDelay = Math.min(totalDelay, 500);

    await new Promise((resolve) => setTimeout(resolve, cappedDelay));
  }

  private generateCryptoShortCode(length: number): string {
    const chars =
      '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ';
    const randomBytes = crypto.randomBytes(length * 2);
    let result = '';

    for (let i = 0; i < length; i++) {
      const randomValue = randomBytes.readUInt16BE(i * 2);
      result += chars[randomValue % chars.length];
    }

    return result;
  }

  private generateTimestampBasedCode(length: number): string {
    const timestamp = Date.now().toString(36);
    const randomBytes = crypto.randomBytes(
      Math.max(1, length - timestamp.length),
    );
    const chars =
      '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ';

    let randomSuffix = '';
    for (let i = 0; i < randomBytes.length; i++) {
      randomSuffix += chars[randomBytes[i] % chars.length];
    }

    const combined = timestamp + randomSuffix;
    return combined.slice(-length);
  }

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

  getCollisionMetrics(): { totalCollisions: number } {
    return {
      totalCollisions: this.collisionCount,
    };
  }
}
