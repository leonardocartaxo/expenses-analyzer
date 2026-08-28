import { readFileSync, existsSync } from 'node:fs';
import { join } from 'node:path';
import React from 'react';
import { render, screen } from '@testing-library/react';
import { ScaffoldView } from '../app/scaffold-view';

describe('scaffold page', () => {
  it('imports @expenses/api-client and does not hardcode a backend fetch URL', () => {
    const pagePath = join(__dirname, '..', 'app', 'page.tsx');
    expect(existsSync(pagePath)).toBe(true);
    const src = readFileSync(pagePath, 'utf8');
    expect(src).toContain('@expenses/api-client');
    expect(src).not.toMatch(/fetch\s*\(\s*['"`]https?:\/\//);
  });

  it('still renders when health is unhealthy', () => {
    render(React.createElement(ScaffoldView, { healthStatus: 'error' }));
    expect(screen.getByTestId('health-status').textContent).toMatch(/error|unhealthy/i);
  });
});
