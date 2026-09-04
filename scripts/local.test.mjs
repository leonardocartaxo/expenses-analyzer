import assert from 'node:assert/strict';
import test from 'node:test';
import { existsSync, readFileSync } from 'node:fs';
import { join } from 'node:path';

const root = join(import.meta.dirname, '..');

test('local kind scripts and package.json local:* commands exist', () => {
  for (const name of ['up.sh', 'down.sh', 'status.sh', 'lib.sh']) {
    assert.equal(existsSync(join(root, 'scripts/local', name)), true, `scripts/local/${name} must exist`);
  }
  const lib = readFileSync(join(root, 'scripts/local/lib.sh'), 'utf8');
  assert.match(lib, /host\.docker\.internal/);
  assert.match(lib, /fix_kind_kubeconfig_for_dood/);
  assert.match(lib, /ensure_kind_kubectl/);
  assert.match(lib, /check_backend_health/);
  assert.match(lib, /kind-darwin-|kind-\$\{os\}/);
  const pkg = JSON.parse(readFileSync(join(root, 'package.json'), 'utf8'));
  for (const script of ['local:up', 'local:down', 'local:status']) {
    assert.equal(typeof pkg.scripts?.[script], 'string', `${script} must be in root package.json`);
  }
});
