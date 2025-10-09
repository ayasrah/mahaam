import { NestFactory } from '@nestjs/core';
import { DocumentBuilder, SwaggerModule } from '@nestjs/swagger';
import Config from './infra/config';

import { AppModule } from './infra/factory';
import { HealthService } from './infra/monitor/health.service';
import Log from './infra/log';
import { Starter } from './infra/starter';

async function bootstrap() {
  process.loadEnvFile();
  Log.init();

  const app = await NestFactory.create(AppModule);

  // Initialize the application using Starter service
  const starter = app.get<Starter>('Starter');
  await starter.start();

  app.setGlobalPrefix(Config.baseUrl);

  // Swagger configuration
  const config = new DocumentBuilder()
    .setTitle('Mahaam API')
    .setDescription('The Mahaam API documentation')
    .setVersion('1.0')
    .addBearerAuth()
    .build();

  const document = SwaggerModule.createDocument(app, config);
  SwaggerModule.setup('mahaam-docs', app, document);

  await app.listen(Config.httpPort);

  const gracefulShutdown = async (signal: string) => {
    console.log(`Received ${signal}, shutting down gracefully...`);
    const healthService = app.get<HealthService>('HealthService');
    await healthService.serverStopped();
    await app.close();
    process.exit(0);
  };

  process.on('SIGINT', () => gracefulShutdown('SIGINT'));
  process.on('SIGTERM', () => gracefulShutdown('SIGTERM'));
}
bootstrap();
