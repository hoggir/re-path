import { Injectable } from '@nestjs/common';
import { BaseRepository } from 'src/database/repositories/base.repository';
import { Model } from 'mongoose';
import { Url } from './schemas/urls.schema';
import { InjectModel } from '@nestjs/mongoose';

@Injectable()
export class UrlRepository extends BaseRepository<Url> {
  constructor(@InjectModel(Url.name) private urlModel: Model<Url>) {
    super(urlModel);
  }

  async findByShortCode(shortCode: string): Promise<Url | null> {
    return this.urlModel
      .findOne({ shortCode, isDeleted: false })
      .select('+shortCode +longUrl')
      .exec();
  }

  async countUrl(idUser: string): Promise<number> {
    return this.urlModel.countDocuments({ _id: idUser }).exec();
  }
}
