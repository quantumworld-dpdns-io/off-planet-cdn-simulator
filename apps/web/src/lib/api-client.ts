const BASE_URL = process.env.NEXT_PUBLIC_CONTROL_API_URL ?? "http://localhost:8080";

async function apiFetch<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    ...init,
    headers: { "Content-Type": "application/json", ...(init?.headers ?? {}) },
  });
  if (!res.ok) throw new Error(`API error ${res.status}: ${path}`);
  return res.json() as Promise<T>;
}

export const api = {
  health: () => apiFetch<{ status: string }>("/v1/health"),
  sites: { list: () => apiFetch<{ sites: unknown[] }>("/v1/sites") },
  nodes: { list: () => apiFetch<{ nodes: unknown[] }>("/v1/nodes") },
};
