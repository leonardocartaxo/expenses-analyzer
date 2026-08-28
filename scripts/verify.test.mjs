import assert from 'node:assert/strict';
import test from 'node:test';
import { checkRequiredScripts, runVerify } from './verify.mjs';

test('verify runs lint then typecheck then test and stops on first failure', async () => {
  const calls = [];
  const result = await runVerify({
    packages: [{ name: 'pkg', scripts: { lint: 'x', typecheck: 'x', test: 'x' } }],
    exec: async (step) => {
      calls.push(step);
      if (step === 'lint') {
        return { ok: false, status: 1 };
      }
      return { ok: true, status: 0 };
    },
  });
  assert.equal(result.ok, false);
  assert.deepEqual(calls, ['lint']);
});

test('verify fails when a workspace package lacks lint, typecheck, or test', async () => {
  const result = await runVerify({
    packages: [{ name: 'broken', scripts: { lint: 'x' } }],
    exec: async () => ({ ok: true, status: 0 }),
  });
  assert.equal(result.ok, false);
  assert.deepEqual(checkRequiredScripts({ lint: 'x' }), ['typecheck', 'test']);
});
