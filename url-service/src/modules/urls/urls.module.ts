import { Module } from '@nestjs/common';
import { MongooseModule } from '@nestjs/mongoose';
import { UrlController } from './urls.controller';
import { UrlRepository } from './urls.repository';
import { Url, UrlSchema } from './schemas/urls.schema';
import { UrlService } from './urls.service';
import { UsersModule } from '../users/users.module';

@Module({
  imports: [
    MongooseModule.forFeature([{ name: Url.name, schema: UrlSchema }]),
    UsersModule,
  ],
  controllers: [UrlController],
  providers: [UrlRepository, UrlService],
  exports: [],
})
export class UrlsModule {}
