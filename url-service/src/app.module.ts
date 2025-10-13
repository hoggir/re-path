import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import databaseConfig from './config/database.config';
import mongodbConfig from './config/mongodb.config';
import authConfig from './config/auth.config';
import redisConfig from './config/redis.config';
import { DatabaseModule } from './database/database.module';
import { CacheModule } from './cache/cache.module';
import { AuthModule } from './modules/auth/auth.module';
import { UsersModule } from './modules/users/users.module';
import { UrlsModule } from './modules/urls/urls.module';
import { HealthModule } from './modules/health/health.module';
import { IdEncryptionService } from './common/utils/id-encryption.service';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
      load: [databaseConfig, mongodbConfig, authConfig, redisConfig],
      envFilePath: ['.env.local', '.env'],
    }),
    CacheModule,
    HealthModule,
    DatabaseModule,
    AuthModule,
    UsersModule,
    UrlsModule,
  ],
  controllers: [],
  providers: [IdEncryptionService],
  exports: [IdEncryptionService],
})
export class AppModule {}
