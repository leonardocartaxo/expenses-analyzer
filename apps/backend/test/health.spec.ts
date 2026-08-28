import { ServiceUnavailableException } from '@nestjs/common';
import { Test } from '@nestjs/testing';
import { HealthController } from '../src/health/health.controller';
import { DataSource } from 'typeorm';

describe('GET /health', () => {
  it('returns 200 { status: "ok" } when postgres ping succeeds', async () => {
    const moduleRef = await Test.createTestingModule({
      controllers: [HealthController],
      providers: [
        {
          provide: DataSource,
          useValue: { query: jest.fn().mockResolvedValue([{ '?column?': 1 }]) },
        },
      ],
    }).compile();
    const controller = moduleRef.get(HealthController);
    await expect(controller.check()).resolves.toEqual({ status: 'ok' });
  });

  it('returns 503 { status: "error" } when postgres is unreachable', async () => {
    const moduleRef = await Test.createTestingModule({
      controllers: [HealthController],
      providers: [
        {
          provide: DataSource,
          useValue: { query: jest.fn().mockRejectedValue(new Error('ECONNREFUSED')) },
        },
      ],
    }).compile();
    const controller = moduleRef.get(HealthController);
    try {
      await controller.check();
      fail('expected unhealthy');
    } catch (err) {
      expect(err).toBeInstanceOf(ServiceUnavailableException);
      const response = (err as ServiceUnavailableException).getResponse();
      expect(response).toEqual({ status: 'error' });
    }
  });

  it('does not include organization, bill, or user fields', async () => {
    const moduleRef = await Test.createTestingModule({
      controllers: [HealthController],
      providers: [{ provide: DataSource, useValue: { query: jest.fn().mockResolvedValue([]) } }],
    }).compile();
    const body = await moduleRef.get(HealthController).check();
    expect(Object.keys(body).sort()).toEqual(['status']);
    expect(JSON.stringify(body)).not.toMatch(/organization|bill|"user"/i);
  });
});
