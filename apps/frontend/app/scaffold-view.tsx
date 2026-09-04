export type HealthStatus = 'ok' | 'error';

export function ScaffoldView({ healthStatus }: { healthStatus: HealthStatus }) {
  const label = healthStatus === 'ok' ? 'ok' : 'error (unhealthy)';
  return (
    <main>
      <h1>Expenses Analyzer scaffold</h1>
      <p data-testid="health-status">Health: {label}</p>
    </main>
  );
}
