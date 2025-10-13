import { registerAs } from '@nestjs/config';

export default registerAs('mongodb', () => {
  const uri = process.env.MONGODB_URI || 'mongodb://localhost:27017/repath';
  const dbName = process.env.MONGODB_DATABASE || 'repath';

  return {
    uri,
    options: {
      dbName,
      retryWrites: true,
      w: 'majority',
      maxPoolSize:
        parseInt(process.env.MONGODB_MAX_POOL_SIZE as string, 10) || 10,
      minPoolSize:
        parseInt(process.env.MONGODB_MIN_POOL_SIZE as string, 10) || 2,
      socketTimeoutMS:
        parseInt(process.env.MONGODB_SOCKET_TIMEOUT as string, 10) || 45000,
      serverSelectionTimeoutMS:
        parseInt(process.env.MONGODB_SERVER_SELECTION_TIMEOUT as string, 10) ||
        5000,
      family: 4,
    },
  };
});
