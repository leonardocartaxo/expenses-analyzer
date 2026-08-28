import assert from 'node:assert/strict';
import test from 'node:test';
import { existsSync, readFileSync } from 'node:fs';
import { join } from 'node:path';

const root = join(import.meta.dirname, '..');
const examples = ['.env.example', '.env.backend.example', '.env.frontend.example'];

test('env examples exist with dummy local credentials only', () => {
  for (const name of examples) {
    assert.equal(existsSync(join(root, name)), true, `${name} must exist`);
  }

  const gitignore = readFileSync(join(root, '.gitignore'), 'utf8');
  assert.match(gitignore, /^\.env$/m);
  assert.match(gitignore, /^\.env\.\*$/m);
  assert.match(gitignore, /^!\.env\.example$/m);
  assert.match(gitignore, /^!\.env\.\*\.example$/m);

  const backend = readFileSync(join(root, '.env.backend.example'), 'utf8');
  assert.match(backend, /DATABASE_HOST=localhost/);
  assert.match(backend, /PORT=3001/);

  const frontend = readFileSync(join(root, '.env.frontend.example'), 'utf8');
  assert.match(frontend, /API_BASE_URL=http:\/\/localhost:3001/);

  const combined = examples.map((name) => readFileSync(join(root, name), 'utf8')).join('\n');
  assert.doesNotMatch(combined, /AKIA[A-Z0-9]{16}/);
  assert.doesNotMatch(combined, /aws_secret_access_key/i);
  assert.doesNotMatch(combined, /rds\.amazonaws\.com/i);
  assert.match(combined, /expenses/);
});
