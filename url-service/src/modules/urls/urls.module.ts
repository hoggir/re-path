import { Module } from '@nestjs/common';
import { MongooseModule } from '@nestjs/mongoose';
import { UrlController } from './urls.controller';
import { UrlRepository } from './urls.repository';
import { Url, UrlSchema } from './schemas/urls.schema';
import { User, UserSchema } from '../users/schemas/user.schema';
import { UrlService } from './urls.service';
import { UsersRepository } from '../users/users.repository';

@Module({
  imports: [
    MongooseModule.forFeature([
      { name: User.name, schema: UserSchema },
      { name: Url.name, schema: UrlSchema },
    ]),
  ],
  controllers: [UrlController],
  providers: [UrlRepository, UrlService, UsersRepository],
  exports: [],
})
export class UrlsModule {}
