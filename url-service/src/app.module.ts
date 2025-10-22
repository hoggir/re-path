import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import databaseConfig from './config/database.config';
import mongodbConfig from './config/mongodb.config';
import authConfig from './config/auth.config';
import redisConfig from './config/redis.config';
import urlConfig from './config/url.config';
import { CommonModule } from './common/common.module';
import { DatabaseModule } from './database/database.module';
import { CacheModule } from './cache/cache.module';
import { AuthModule } from './modules/auth/auth.module';
import { UsersModule } from './modules/users/users.module';
import { UrlsModule } from './modules/urls/urls.module';
import { HealthModule } from './modules/health/health.module';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
      load: [databaseConfig, mongodbConfig, authConfig, redisConfig, urlConfig],
      envFilePath: ['.env.local', '.env'],
    }),
    CommonModule,
    CacheModule,
    DatabaseModule,
    HealthModule,
    AuthModule,
    UsersModule,
    UrlsModule,
  ],
  controllers: [],
  providers: [],
  exports: [],
})
export class AppModule {}
