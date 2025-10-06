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

  async findByOriginalUrl(
    originalUrl: string,
    userId?: string,
  ): Promise<Url | null> {
    const query: any = { originalUrl, isDeleted: false };
    if (userId) {
      query.userId = userId;
    }
    return this.urlModel.findOne(query).exec();
  }

  async findByCustomAlias(customAlias: string): Promise<Url | null> {
    return this.urlModel.findOne({ customAlias, isDeleted: false }).exec();
  }

  // async incrementClickCount(shortCode: string): Promise<void> {
  //   await this.urlModel
  //     .findOneAndUpdate({ shortCode }, { $inc: { clickCount: 1 } })
  //     .exec();
  // }

  // async addAnalytics(
  //   shortCode: string,
  //   analyticsData: {
  //     ipAddress: string;
  //     userAgent: string;
  //     referer: string;
  //     country?: string;
  //     city?: string;
  //   },
  // ): Promise<void> {
  //   await this.urlModel
  //     .findOneAndUpdate(
  //       { shortCode },
  //       {
  //         $push: {
  //           analytics: {
  //             ...analyticsData,
  //             clickedAt: new Date(),
  //           },
  //         },
  //         $inc: { clickCount: 1 },
  //       },
  //     )
  //     .exec();
  // }

  // async findByUserId(userId: string, page: number = 1, limit: number = 10) {
  //   return this.paginate({ userId, isActive: true }, page, limit, {
  //     createdAt: -1,
  //   });
  // }

  // async getTotalClicksByUserId(userId: string): Promise<number> {
  //   const result = await this.urlModel.aggregate([
  //     { $match: { userId, isDeleted: false } },
  //     { $group: { _id: null, total: { $sum: '$clickCount' } } },
  //   ]);
  //   return result[0]?.total || 0;
  // }

  async checkShortCodeExists(shortCode: string): Promise<boolean> {
    return this.exists({ shortCode });
  }

  async checkCustomAliasExists(customAlias: string): Promise<boolean> {
    const count = await this.urlModel
      .countDocuments({ shortCode: customAlias, isDeleted: false })
      .exec();
    return count > 0;
  }
}
