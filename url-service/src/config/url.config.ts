import { registerAs } from '@nestjs/config';

export default registerAs('url', () => ({
  defaultTtlDays: parseInt(process.env.URL_DEFAULT_TTL_DAYS || '7', 10),
  shortCodeLength: parseInt(process.env.URL_SHORT_CODE_LENGTH || '6', 10),
  maxRetries: parseInt(process.env.URL_MAX_RETRIES || '10', 10),
  baseUrl: process.env.BASE_URL || 'http://localhost:3010',
}));
