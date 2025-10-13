import { registerAs } from '@nestjs/config';

export default registerAs('redis', () => ({
  host: process.env.REDIS_HOST || 'localhost',
  port: parseInt(process.env.REDIS_PORT as string, 10) || 6379,
  password: process.env.REDIS_PASSWORD || undefined,
  db: parseInt(process.env.REDIS_DB as string, 10) || 0,
  ttl: parseInt(process.env.REDIS_TTL as string, 10) || 3600,
  maxRetriesPerRequest: 3,
}));
