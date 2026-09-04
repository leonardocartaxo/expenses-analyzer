import { writeFileSync } from 'node:fs';
import { join } from 'node:path';
import { Test } from '@nestjs/testing';
import { DocumentBuilder, SwaggerModule } from '@nestjs/swagger';
import { DataSource } from 'typeorm';
import { HealthController } from './health/health.controller';

async function exportOpenApi() {
  const moduleRef = await Test.createTestingModule({
    controllers: [HealthController],
    providers: [{ provide: DataSource, useValue: { query: async () => [{ ok: 1 }] } }],
  }).compile();
  const app = moduleRef.createNestApplication();
  await app.init();
  const config = new DocumentBuilder()
    .setTitle('Expenses Analyzer API')
    .setDescription('Bootstrap health/scaffold API. Nest @nestjs/swagger is the runtime source.')
    .setVersion('0.0.0')
    .addServer('http://localhost:3001', 'Host and Dev Container Nest')
    .build();
  const document = SwaggerModule.createDocument(app, config);
  const schemas = document.components?.schemas ?? {};
  if ('HealthResponseDto' in schemas) {
    schemas.HealthResponse = schemas.HealthResponseDto;
    delete schemas.HealthResponseDto;
  }
  const health = schemas.HealthResponse as { additionalProperties?: boolean } | undefined;
  if (health) {
    health.additionalProperties = false;
  }
  const json = JSON.stringify(document, null, 2).replaceAll('HealthResponseDto', 'HealthResponse');
  const out = join(__dirname, '../../../packages/api-client/openapi.json');
  writeFileSync(out, `${json}\n`);
  await app.close();
}

void exportOpenApi();
