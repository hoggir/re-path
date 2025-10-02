// src/database/database.module.ts
import { Module, Global } from '@nestjs/common';
import { MongooseModule } from '@nestjs/mongoose';
import { ConfigService } from '@nestjs/config';
import { DatabaseService } from './database.service';

@Global()
@Module({
  imports: [
    MongooseModule.forRootAsync({
      inject: [ConfigService],
      useFactory: async (configService: ConfigService) => {
        const uri = configService.get<string>('database.uri');
        const options = configService.get('database.options');

        return {
          uri,
          ...options,
          connectionFactory: (connection) => {
            connection.on('connected', () => {
              console.log('✅ MongoDB connected successfully');
            });

            connection.on('disconnected', () => {
              console.log('❌ MongoDB disconnected');
            });

            connection.on('error', (error) => {
              console.error('❌ MongoDB connection error:', error);
            });

            // Middleware untuk logging (development only)
            if (process.env.NODE_ENV !== 'production') {
              connection.set('debug', true);
            }

            return connection;
          },
        };
      },
    }),
  ],
  providers: [DatabaseService],
  exports: [DatabaseService],
})
export class DatabaseModule {}