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
  const port = configService.get('PORT', 3000);

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
      disableErrorMessages: configService.get('NODE_ENV') === 'production', // Hide error details di production
    }),
  );

  // Swagger Documentation
  if (configService.get('NODE_ENV') !== 'production') {
    const config = new DocumentBuilder()
      .setTitle('Re:Path API Documentation')
      .setDescription(
        'Complete REST API documentation for Re:Path application. This API provides authentication, user management, and other core functionalities.',
      )
      .setVersion('1.0')
      .setContact('Re:Path Support', 'https://repath.com', 'support@repath.com')
      .setLicense('MIT', 'https://opensource.org/licenses/MIT')
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
      .addTag(
        'Authentication',
        'Authentication endpoints for login, register, and token management',
      )
      .addTag('Users', 'User management endpoints')
      .addTag('Health', 'Health check and monitoring endpoints')
      .addTag('URL Shortener', 'Shortened URL management endpoints')
      .addServer('http://localhost:8080', 'Nginx Gateway')
      .addServer(`http://localhost:${port}`, 'Local Development')
      // .addServer('https://api-staging.repath.com', 'Staging')
      // .addServer('https://api.repath.com', 'Production')
      .build();

    const document = SwaggerModule.createDocument(app, config, {
      deepScanRoutes: true,
    });

    SwaggerModule.setup('api/docs', app, document, {
      swaggerOptions: {
        persistAuthorization: true,
        docExpansion: 'none',
        filter: true,
        showRequestDuration: true,
        syntaxHighlight: {
          theme: 'monokai',
        },
      },
      customSiteTitle: 'Re:Path API Docs',
      customfavIcon: 'https://repath.com/favicon.ico',
      customCss: '.swagger-ui .topbar { display: none }',
    });

    console.log(
      `📚 Swagger documentation available at: http://localhost:${port}/api/docs`,
    );
  }

  // Graceful shutdown
  app.enableShutdownHooks();

  await app.listen(port);

  console.log(`🚀 Application is running on: http://localhost:${port}/api`);
}

bootstrap().catch((err) => {
  console.error('❌ Error starting application:', err);
  process.exit(1);
});
