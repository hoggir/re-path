import { NestFactory } from '@nestjs/core';
import { ValidationPipe, VersioningType } from '@nestjs/common';
import { SwaggerModule, DocumentBuilder } from '@nestjs/swagger';
import { ConfigService } from '@nestjs/config';
import helmet from 'helmet';
import compression from 'compression';
import { AppModule } from './app.module';

async function bootstrap() {
  const app = await NestFactory.create(AppModule, {
    logger: ['error', 'warn', 'log', 'debug', 'verbose'],
  });

  const configService = app.get(ConfigService);

  // Global prefix untuk semua routes
  app.setGlobalPrefix('api');

  // API Versioning
  app.enableVersioning({
    type: VersioningType.URI,
    defaultVersion: '1',
  });

  // Security - Helmet
  app.use(helmet());

  // CORS Configuration
  app.enableCors({
    origin: configService.get('CORS_ORIGIN', '*'),
    credentials: true,
    methods: ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'OPTIONS'],
    allowedHeaders: ['Content-Type', 'Authorization', 'Accept'],
  });

  // Compression middleware
  app.use(compression());

  // Global Validation Pipe
  app.useGlobalPipes(
    new ValidationPipe({
      whitelist: true, // Strip properties yang tidak ada di DTO
      forbidNonWhitelisted: true, // Throw error jika ada property tidak dikenal
      transform: true, // Auto-transform payloads ke DTO instances
      transformOptions: {
        enableImplicitConversion: true,
      },
      disableErrorMessages: 
        configService.get('NODE_ENV') === 'production', // Hide error details di production
    }),
  );

  // Swagger Documentation
  if (configService.get('NODE_ENV') !== 'production') {
    const config = new DocumentBuilder()
      .setTitle('Re P4th API Documentation')
      .setDescription('REST API Documentation')
      .setVersion('1.0')
      .addBearerAuth(
        {
          type: 'http',
          scheme: 'bearer',
          bearerFormat: 'JWT',
          name: 'JWT',
          description: 'Enter JWT token',
          in: 'header',
        },
        'JWT-auth',
      )
      .build();

    const document = SwaggerModule.createDocument(app, config);
    SwaggerModule.setup('api/docs', app, document, {
      swaggerOptions: {
        persistAuthorization: true,
      },
    });
  }

  // Graceful shutdown
  app.enableShutdownHooks();

  const port = configService.get('PORT', 3000);
  // const port = 3000;
  await app.listen(port);

  console.log(`🚀 Application is running on: http://localhost:${port}/api`);
  console.log(`📚 Swagger documentation: http://localhost:${port}/api/docs`);
}

bootstrap().catch((err) => {
  console.error('❌ Error starting application:', err);
  process.exit(1);
});