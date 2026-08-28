import assert from 'node:assert/strict';
import test from 'node:test';
import { existsSync, readFileSync, readdirSync } from 'node:fs';
import { join } from 'node:path';

const overlay = join(import.meta.dirname, '..', 'deploy/kustomize/overlays/kind');

function collectFiles(dir, acc = []) {
  for (const entry of readdirSync(dir, { withFileTypes: true })) {
    const path = join(dir, entry.name);
    if (entry.isDirectory()) {
      collectFiles(path, acc);
    } else {
      acc.push(path);
    }
  }
  return acc;
}

test('kind overlay includes Postgres, Nest, and frontend and does not use Helm', () => {
  assert.equal(existsSync(join(overlay, 'kustomization.yaml')), true);
  const files = collectFiles(overlay);
  const text = files.map((f) => readFileSync(f, 'utf8')).join('\n');
  assert.match(text, /postgres:18\.4/);
  assert.match(text, /expenses-analyzer-backend:local/);
  assert.match(text, /expenses-analyzer-frontend:local/);
  assert.doesNotMatch(text, /helm|HelmChart|apiVersion:\s*helm/i);
});
