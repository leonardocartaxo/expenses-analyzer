import { spawnSync } from 'node:child_process';
import { readFileSync, readdirSync, existsSync } from 'node:fs';
import { join, resolve } from 'node:path';
import { pathToFileURL } from 'node:url';

export const STEPS = ['lint', 'typecheck', 'test'];
export const REQUIRED_SCRIPTS = ['lint', 'typecheck', 'test'];

export function checkRequiredScripts(scripts = {}) {
  return REQUIRED_SCRIPTS.filter((name) => typeof scripts[name] !== 'string' || scripts[name].length === 0);
}

export function workspacePackageDirs(root) {
  const dirs = [];
  for (const group of ['apps', 'packages']) {
    const base = join(root, group);
    if (!existsSync(base)) {
      continue;
    }
    for (const name of readdirSync(base)) {
      const pkgDir = join(base, name);
      if (existsSync(join(pkgDir, 'package.json'))) {
        dirs.push(pkgDir);
      }
    }
  }
  return dirs;
}

export function loadWorkspacePackages(root) {
  return workspacePackageDirs(root).map((dir) => {
    const pkg = JSON.parse(readFileSync(join(dir, 'package.json'), 'utf8'));
    return { name: pkg.name ?? dir, scripts: pkg.scripts ?? {}, dir };
  });
}

export async function runVerify({
  exec,
  packages,
  root = process.cwd(),
} = {}) {
  const pkgs = packages ?? loadWorkspacePackages(root);
  for (const pkg of pkgs) {
    const missing = checkRequiredScripts(pkg.scripts);
    if (missing.length > 0) {
      return {
        ok: false,
        status: 1,
        error: `${pkg.name} is missing scripts: ${missing.join(', ')}`,
      };
    }
  }

  for (const step of STEPS) {
    const result = typeof exec === 'function' ? await exec(step) : defaultExec(step);
    if (!result.ok) {
      return { ok: false, status: result.status ?? 1 };
    }
  }
  return { ok: true, status: 0 };
}

function defaultExec(step) {
  const result = spawnSync('pnpm', [step], {
    stdio: 'inherit',
    cwd: process.cwd(),
    env: process.env,
  });
  const status = result.status ?? 1;
  return { ok: status === 0, status };
}

const invokedDirectly =
  Boolean(process.argv[1]) && import.meta.url === pathToFileURL(resolve(process.argv[1])).href;

if (invokedDirectly) {
  const { ok, status, error } = await runVerify();
  if (error) {
    console.error(error);
  }
  process.exit(ok ? 0 : status);
}
