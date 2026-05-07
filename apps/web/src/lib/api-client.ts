import type { Site, Node, CacheObject, CachePolicy, PreloadJob, AuditLog, BandwidthWindow } from "@/types";

export interface CacheHitPoint { hour: string; hits: number; misses: number }
export interface PriorityBucket { level: string; count: number; total_bytes: number }
export interface NodeFill { node_id: string; node_name: string; used_bytes: number; max_bytes: number }

const BASE_URL = process.env.NEXT_PUBLIC_CONTROL_API_URL ?? "http://localhost:8080";
const DEV_ORG_ID = process.env.NEXT_PUBLIC_DEV_ORG_ID ?? "00000000-0000-0000-0000-000000000001";

function toQuery(params?: Record<string, string | number | undefined | null>): string {
  if (!params) return "";
  const q = Object.entries(params)
    .filter(([, v]) => v != null && v !== "")
    .map(([k, v]) => `${encodeURIComponent(k)}=${encodeURIComponent(String(v))}`)
    .join("&");
  return q ? `?${q}` : "";
}

async function apiFetch<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      "X-Org-ID": DEV_ORG_ID,
      ...(init?.headers ?? {}),
    },
  });
  if (!res.ok) {
    const text = await res.text().catch(() => "");
    throw new Error(`API ${res.status} ${path}: ${text}`);
  }
  if (res.status === 204) return undefined as T;
  return res.json() as Promise<T>;
}

export const api = {
  health: () => apiFetch<{ status: string }>("/v1/health"),

  sites: {
    list: () => apiFetch<{ sites: Site[] }>("/v1/sites"),
    get: (id: string) => apiFetch<Site>(`/v1/sites/${id}`),
    create: (body: { name: string; location?: string; description?: string }) =>
      apiFetch<Site>("/v1/sites", { method: "POST", body: JSON.stringify(body) }),
  },

  nodes: {
    list: (params?: { site_id?: string }) =>
      apiFetch<{ nodes: Node[] }>(`/v1/nodes${toQuery(params)}`),
    get: (id: string) => apiFetch<Node>(`/v1/nodes/${id}`),
  },

  cacheObjects: {
    list: (params?: { site_id?: string; priority?: string; status?: string }) =>
      apiFetch<{ objects: CacheObject[] }>(`/v1/objects${toQuery(params)}`),
    get: (id: string) => apiFetch<CacheObject>(`/v1/objects/${id}`),
    create: (body: Partial<CacheObject>) =>
      apiFetch<CacheObject>("/v1/objects", { method: "POST", body: JSON.stringify(body) }),
    pin: (id: string) => apiFetch<void>(`/v1/objects/${id}/pin`, { method: "POST" }),
    unpin: (id: string) => apiFetch<void>(`/v1/objects/${id}/unpin`, { method: "POST" }),
  },

  policies: {
    list: (params?: { site_id?: string }) =>
      apiFetch<{ policies: CachePolicy[] }>(`/v1/policies${toQuery(params)}`),
    get: (id: string) => apiFetch<CachePolicy>(`/v1/policies/${id}`),
    create: (body: Partial<CachePolicy>) =>
      apiFetch<CachePolicy>("/v1/policies", { method: "POST", body: JSON.stringify(body) }),
    update: (id: string, body: Partial<CachePolicy>) =>
      apiFetch<CachePolicy>(`/v1/policies/${id}`, { method: "PUT", body: JSON.stringify(body) }),
  },

  preloadJobs: {
    list: (params?: { site_id?: string; status?: string }) =>
      apiFetch<{ jobs: PreloadJob[] }>(`/v1/preload-jobs${toQuery(params)}`),
    get: (id: string) => apiFetch<PreloadJob>(`/v1/preload-jobs/${id}`),
    create: (body: { site_id: string; name: string; bandwidth_budget_bytes?: number }) =>
      apiFetch<PreloadJob>("/v1/preload-jobs", { method: "POST", body: JSON.stringify(body) }),
    cancel: (id: string) =>
      apiFetch<void>(`/v1/preload-jobs/${id}/cancel`, { method: "POST" }),
  },

  auditLogs: {
    list: (params?: { resource_type?: string; limit?: number }) =>
      apiFetch<{ logs: AuditLog[] }>(`/v1/audit-logs${toQuery(params)}`),
  },

  bandwidthWindows: {
    list: (params?: { site_id?: string }) =>
      apiFetch<{ windows: BandwidthWindow[] }>(`/v1/bandwidth-windows${toQuery(params)}`),
  },

  analytics: {
    cacheHits: () => apiFetch<{ points: CacheHitPoint[] }>("/v1/analytics/cache-hits"),
    priorityDistribution: () => apiFetch<{ buckets: PriorityBucket[] }>("/v1/analytics/priority-distribution"),
    nodeFill: () => apiFetch<{ nodes: NodeFill[] }>("/v1/analytics/node-fill"),
  },
};
