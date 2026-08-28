import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { HealthModule } from './health/health.module';

@Module({
  imports: [
    TypeOrmModule.forRoot({
      type: 'postgres',
      host: process.env.DATABASE_HOST ?? 'localhost',
      port: Number(process.env.DATABASE_PORT ?? 5432),
      username: process.env.DATABASE_USER ?? 'expenses',
      password: process.env.DATABASE_PASSWORD ?? 'expenses',
      database: process.env.DATABASE_NAME ?? 'expenses',
      entities: [],
      synchronize: false,
      retryAttempts: 10,
      retryDelay: 3000,
    }),
    HealthModule,
  ],
})
export class AppModule {}
