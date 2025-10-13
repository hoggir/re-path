import { Module, Global } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { MongooseModule } from '@nestjs/mongoose';
import { ConfigService } from '@nestjs/config';
import { DatabaseService } from './database.service';
import { PostgresService } from './postgres.service';

@Global()
@Module({
  imports: [
    TypeOrmModule.forRootAsync({
      inject: [ConfigService],
      useFactory: async (configService: ConfigService) => {
        const dbConfig = configService.get('database');

        console.log('✅ Initializing PostgreSQL connection for Users/Auth...');

        return {
          ...dbConfig,
          autoLoadEntities: true,
        };
      },
    }),
    MongooseModule.forRootAsync({
      inject: [ConfigService],
      useFactory: async (configService: ConfigService) => {
        const uri = configService.get<string>('mongodb.uri');
        const options = configService.get('mongodb.options');

        console.log('✅ Initializing MongoDB connection for URLs...');

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

            if (process.env.NODE_ENV !== 'production') {
              // connection.set('debug', true);
            }

            return connection;
          },
        };
      },
    }),
  ],
  providers: [DatabaseService, PostgresService],
  exports: [DatabaseService, PostgresService],
})
export class DatabaseModule {}
