import assert from 'node:assert/strict';
import test from 'node:test';
import { existsSync, readFileSync } from 'node:fs';
import { join } from 'node:path';

const root = join(import.meta.dirname, '..');
const path = join(root, '.devcontainer/devcontainer.json');

test('Dev Container is not Postgres and documents host.docker.internal plus app ports', () => {
  assert.equal(existsSync(path), true);
  const raw = readFileSync(path, 'utf8');
  const cfg = JSON.parse(raw);
  assert.notEqual(cfg.build?.dockerfile, undefined);
  assert.doesNotMatch(JSON.stringify(cfg.image ?? ''), /postgres/i);
  assert.doesNotMatch(raw, /"service"\s*:\s*"postgres"/);
  assert.match(raw, /host\.docker\.internal/);
  assert.match(raw, /\/var\/run\/docker\.sock/);
  assert.equal(cfg.containerEnv?.DATABASE_HOST, 'host.docker.internal');
  const ports = cfg.forwardPorts ?? [];
  for (const port of [3000, 3001]) {
    assert.equal(ports.includes(port), true, `forwardPorts must include ${port}`);
  }
  // Compose already publishes 5432 on the host. Forwarding 5432 from the Dev Container
  // conflicts with that publish under JetBrains/VS Code and breaks pg ("Connection terminated unexpectedly").
  assert.equal(ports.includes(5432), false, 'forwardPorts must not include 5432');
});
