import { Injectable } from '@nestjs/common';
import { InjectModel } from '@nestjs/mongoose';
import { Model } from 'mongoose';
import { Url } from './schemas/urls.schema';
import { BaseRepository } from 'src/database/repositories/base.repository';

@Injectable()
export class UrlRepository extends BaseRepository<Url> {
  constructor(@InjectModel(Url.name) private urlModel: Model<Url>) {
    super(urlModel);
  }

  async findByShortCode(shortCode: string): Promise<Url | null> {
    return this.urlModel
      .findOne({
        shortCode,
        isDeleted: false,
        isActive: true,
        $or: [
          { expiresAt: { $exists: false } },
          { expiresAt: { $gt: new Date() } },
        ],
      })
      .exec();
  }
}
