import {
  Injectable,
  Logger,
  OnModuleInit,
  OnModuleDestroy,
} from '@nestjs/common';
import { InjectConnection } from '@nestjs/mongoose';
import { Connection } from 'mongoose';

@Injectable()
export class DatabaseService implements OnModuleInit, OnModuleDestroy {
  private readonly logger = new Logger(DatabaseService.name);

  constructor(@InjectConnection() private readonly connection: Connection) {}

  async onModuleInit() {
    const { host, port, name } = this.connection;
    this.logger.log(`Database connected to ${host}:${port}/${name}`);
  }

  async onModuleDestroy() {
    await this.connection.close();
    this.logger.log('Database connection closed');
  }

  getDbHandle(): Connection {
    return this.connection;
  }

  async checkConnection(): Promise<boolean> {
    return this.connection.readyState === 1;
  }

  async getStats() {
    const db = this.connection.db;
    const stats = await db?.stats();
    return {
      database: db?.databaseName,
      collections: stats?.collections,
      dataSize: `${(stats?.dataSize / 1024 / 1024).toFixed(2)} MB`,
      storageSize: `${(stats?.storageSize / 1024 / 1024).toFixed(2)} MB`,
      indexes: stats?.indexes,
      indexSize: `${(stats?.indexSize / 1024 / 1024).toFixed(2)} MB`,
    };
  }

  async dropDatabase() {
    if (process.env.NODE_ENV === 'production') {
      throw new Error('Cannot drop database in production');
    }
    await this.connection.dropDatabase();
    this.logger.warn('Database dropped');
  }
}
