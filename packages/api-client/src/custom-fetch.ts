const DEFAULT_BASE = 'http://localhost:3001';

export async function customFetch<T>(url: string, options?: RequestInit): Promise<T> {
  const base = (process.env.API_BASE_URL || DEFAULT_BASE).replace(/\/$/, '');
  const resolved = url.startsWith('http') ? url : `${base}${url.startsWith('/') ? url : `/${url}`}`;
  const res = await fetch(resolved, { ...options });
  const body = [204, 205, 304].includes(res.status) ? null : await res.text();
  const data = body ? JSON.parse(body) : {};
  return { data, status: res.status, headers: res.headers } as T;
}
