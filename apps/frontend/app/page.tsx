import { getHealth } from '@expenses/api-client';
import { ScaffoldView } from './scaffold-view';

// Health must be checked per request (kind/host), not frozen at `next build`.
export const dynamic = 'force-dynamic';

export default async function Page() {
  let healthStatus: 'ok' | 'error' = 'error';
  try {
    const result = await getHealth();
    const httpStatus = result.status;
    const payload = result.data;
    healthStatus = httpStatus === 200 && payload?.status === 'ok' ? 'ok' : 'error';
  } catch {
    healthStatus = 'error';
  }
  return <ScaffoldView healthStatus={healthStatus} />;
}
