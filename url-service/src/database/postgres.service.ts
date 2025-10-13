import {
  Injectable,
  Logger,
  OnModuleInit,
  OnModuleDestroy,
} from '@nestjs/common';
import { InjectDataSource } from '@nestjs/typeorm';
import { DataSource } from 'typeorm';

@Injectable()
export class PostgresService implements OnModuleInit, OnModuleDestroy {
  private readonly logger = new Logger(PostgresService.name);

  constructor(@InjectDataSource() private readonly dataSource: DataSource) {}

  async onModuleInit() {
    const { host, port, database } = this.dataSource.options as any;
    this.logger.log(`PostgreSQL connected to ${host}:${port}/${database}`);
  }

  async onModuleDestroy() {
    this.logger.log('🔌 Closing PostgreSQL connection...');
    if (this.dataSource.isInitialized) {
      await this.dataSource.destroy();
      this.logger.log('✅ PostgreSQL connection closed successfully');
    }
  }

  getDataSource(): DataSource {
    return this.dataSource;
  }

  async checkConnection(): Promise<boolean> {
    return this.dataSource.isInitialized;
  }

  async query(sql: string, parameters?: any[]): Promise<any> {
    return this.dataSource.query(sql, parameters);
  }
}
