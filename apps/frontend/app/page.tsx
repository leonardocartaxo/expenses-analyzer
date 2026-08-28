import { getHealth } from '@expenses/api-client';
import { ScaffoldView } from './scaffold-view';

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
