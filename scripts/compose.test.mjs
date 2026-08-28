import assert from 'node:assert/strict';
import test from 'node:test';
import { existsSync, readFileSync } from 'node:fs';
import { join } from 'node:path';

const root = join(import.meta.dirname, '..');
const composePath = join(root, 'compose.yaml');

test('compose.yaml has a single postgres:18.4 service and no Nest/Next services', () => {
  assert.equal(existsSync(composePath), true, 'compose.yaml must exist');
  const text = readFileSync(composePath, 'utf8');
  assert.match(text, /image:\s*['"]?postgres:18\.4['"]?/);
  const serviceNames = [...text.matchAll(/^\s{2}([A-Za-z0-9_-]+):\s*$/gm)].map((m) => m[1]);
  assert.deepEqual(serviceNames, ['postgres']);
  assert.doesNotMatch(text, /nest|next|@expenses\/backend|@expenses\/frontend/i);
});
