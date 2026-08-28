import assert from 'node:assert/strict';
import test from 'node:test';
import { existsSync, readdirSync, readFileSync, statSync } from 'node:fs';
import { join } from 'node:path';

const root = join(import.meta.dirname, '..');
const domain = /organization|establishment|transaction|bill-manager|OrgAdmin/i;

function walk(dir) {
  const files = [];
  if (!existsSync(dir)) {
    return files;
  }
  for (const name of readdirSync(dir)) {
    const path = join(dir, name);
    if (statSync(path).isDirectory()) {
      files.push(...walk(path));
    } else if (/\.(ts|tsx|js|mjs)$/.test(name)) {
      files.push(path);
    }
  }
  return files;
}

test('backend and frontend scaffold have no organization/bill/user domain routes or screens', () => {
  const files = [...walk(join(root, 'apps/backend/src')), ...walk(join(root, 'apps/frontend/app'))];
  assert.ok(files.length > 0);
  for (const file of files) {
    const text = readFileSync(file, 'utf8');
    assert.doesNotMatch(text, domain, file);
    assert.doesNotMatch(text, /\/organizations\b|\/bills\b|\/users\b/);
  }
});

test('this slice has no CDK, Helm, Amplify apply, or pnpm wake/sleep', () => {
  const forbiddenNames = /cdk|helm|amplify|Chart\.yaml|pnpm-wake|wake\.sh|sleep\.sh/;
  function scan(dir) {
    if (!existsSync(dir)) {
      return;
    }
    for (const name of readdirSync(dir)) {
      const path = join(dir, name);
      assert.doesNotMatch(name, forbiddenNames, path);
      if (statSync(path).isDirectory() && name !== 'node_modules') {
        scan(path);
      }
    }
  }
  scan(join(root, 'deploy'));
  scan(join(root, '.github'));
  const workflow = readFileSync(join(root, '.github/workflows/verify.yml'), 'utf8');
  assert.doesNotMatch(workflow, /aws-actions|AWS_ACCESS_KEY|secrets\.AWS/i);
  const localScripts = ['up.sh', 'down.sh', 'status.sh']
    .map((name) => readFileSync(join(root, 'scripts/local', name), 'utf8'))
    .join('\n');
  assert.doesNotMatch(localScripts, /pnpm wake|pnpm sleep|aws /);
});
