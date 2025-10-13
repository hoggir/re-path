import { Injectable, OnModuleInit, OnModuleDestroy } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import Redis from 'ioredis';

@Injectable()
export class RedisService implements OnModuleInit, OnModuleDestroy {
  private client: Redis;

  constructor(private configService: ConfigService) {}

  async onModuleInit() {
    const redisConfig = this.configService.get('redis');

    console.log('🔧 Redis Config:', {
      host: redisConfig.host,
      port: redisConfig.port,
      db: redisConfig.db,
    });

    this.client = new Redis({
      host: redisConfig.host,
      port: redisConfig.port,
      password: redisConfig.password || undefined,
      db: redisConfig.db,
      retryStrategy: (times) => {
        const delay = Math.min(times * 50, 2000);
        return delay;
      },
    });

    this.client.on('connect', () => {
      console.log('✅ Redis connected successfully');
    });

    this.client.on('error', (err) => {
      console.error('❌ Redis connection error:', err);
    });
  }

  async onModuleDestroy() {
    console.log('🔌 Disconnecting from Redis...');
    await this.client.quit();
    console.log('✅ Redis connection closed');
  }

  async get<T>(key: string): Promise<T | null> {
    const value = await this.client.get(key);
    if (!value) return null;

    try {
      return JSON.parse(value) as T;
    } catch {
      return value as T;
    }
  }

  async set(key: string, value: any, ttlSeconds?: number): Promise<void> {
    const stringValue =
      typeof value === 'string' ? value : JSON.stringify(value);

    if (ttlSeconds) {
      await this.client.setex(key, ttlSeconds, stringValue);
    } else {
      await this.client.set(key, stringValue);
    }

    console.log(`✅ Redis SET: ${key} (TTL: ${ttlSeconds || 'no expiry'}s)`);
  }

  async del(key: string): Promise<void> {
    await this.client.del(key);
    console.log(`🗑️  Redis DEL: ${key}`);
  }

  async keys(pattern: string): Promise<string[]> {
    return this.client.keys(pattern);
  }

  async ttl(key: string): Promise<number> {
    return this.client.ttl(key);
  }

  getClient(): Redis {
    return this.client;
  }
}
