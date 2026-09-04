import { readFileSync, existsSync } from 'node:fs';
import { join } from 'node:path';

describe('openapi.json health contract', () => {
  const path = join(__dirname, '..', 'openapi.json');

  it('documents /health and HealthResponse without domain fields', () => {
    expect(existsSync(path)).toBe(true);
    const doc = JSON.parse(readFileSync(path, 'utf8')) as {
      paths?: Record<string, unknown>;
      components?: {
        schemas?: {
          HealthResponse?: { additionalProperties?: boolean; properties?: Record<string, unknown> };
        };
      };
    };
    expect(doc.paths?.['/health']).toBeDefined();
    const schema = doc.components?.schemas?.HealthResponse;
    expect(schema).toBeDefined();
    expect(schema?.additionalProperties).toBe(false);
    expect(Object.keys(schema?.properties ?? {})).toEqual(['status']);
    expect(JSON.stringify(doc)).not.toMatch(/organization|Establishment|transaction/i);
  });
});
