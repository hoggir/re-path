import { Injectable } from '@nestjs/common';
import { UrlRepository } from './urls.repository';
import { ConflictException } from 'src/common/exceptions/custom-exceptions';

@Injectable()
export class UrlService {
  constructor(private readonly urlRepository: UrlRepository) {}

  async createShortUrl(data: {
    longUrl: string;
    idUser: string;
  }): Promise<string> {
    try {
      const { longUrl, idUser } = data;
      const countUrl = await this.urlRepository.countUrl(idUser);
      console.log("🚀 ~ UrlService ~ createShortUrl ~ countUrl:", countUrl)
      const existingUrl = await this.urlRepository.findByShortCode(longUrl);
      if (existingUrl) {
        throw new ConflictException('URL already exists');
      }
      return longUrl;
    } catch (error) {
      throw error;
    }
  }
}
