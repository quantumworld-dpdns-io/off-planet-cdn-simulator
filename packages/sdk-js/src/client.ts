import type { Site, Node, CacheObject, PreloadJob } from "./types";

export class OffPlanetCdnClient {
  private baseUrl: string;
  private token: string;

  constructor(baseUrl: string, token: string) {
    this.baseUrl = baseUrl.replace(/\/$/, "");
    this.token = token;
  }

  private async fetch<T>(path: string, init?: RequestInit): Promise<T> {
    const res = await fetch(`${this.baseUrl}${path}`, {
      ...init,
      headers: { "Content-Type": "application/json", Authorization: `Bearer ${this.token}`, ...(init?.headers ?? {}) },
    });
    if (!res.ok) throw new Error(`HTTP ${res.status}: ${path}`);
    return res.json() as Promise<T>;
  }

  health() { return this.fetch<{ status: string }>("/v1/health"); }

  sites = {
    list: () => this.fetch<{ sites: Site[] }>("/v1/sites"),
    get: (id: string) => this.fetch<Site>(`/v1/sites/${id}`),
    create: (body: Omit<Site, "id"|"created_at"|"updated_at">) =>
      this.fetch<Site>("/v1/sites", { method: "POST", body: JSON.stringify(body) }),
  };

  nodes = {
    list: () => this.fetch<{ nodes: Node[] }>("/v1/nodes"),
    register: (body: Omit<Node, "id"|"created_at"|"updated_at">) =>
      this.fetch<Node>("/v1/nodes/register", { method: "POST", body: JSON.stringify(body) }),
    heartbeat: (nodeId: string, body: object) =>
      this.fetch(`/v1/nodes/${nodeId}/heartbeat`, { method: "POST", body: JSON.stringify(body) }),
  };

  cache = {
    list: (params?: { site_id?: string; priority_class?: string }) => {
      const qs = params ? "?" + new URLSearchParams(params as Record<string, string>).toString() : "";
      return this.fetch<{ objects: CacheObject[] }>(`/v1/cache/objects${qs}`);
    },
    get: (id: string) => this.fetch<CacheObject>(`/v1/cache/objects/${id}`),
    create: (body: object) => this.fetch<CacheObject>("/v1/cache/objects", { method: "POST", body: JSON.stringify(body) }),
    pin: (id: string) => this.fetch(`/v1/cache/objects/${id}/pin`, { method: "POST" }),
    unpin: (id: string) => this.fetch(`/v1/cache/objects/${id}/unpin`, { method: "POST" }),
  };

  preload = {
    create: (body: object) => this.fetch<PreloadJob>("/v1/preload/jobs", { method: "POST", body: JSON.stringify(body) }),
    list: () => this.fetch<{ jobs: PreloadJob[] }>("/v1/preload/jobs"),
    get: (id: string) => this.fetch<PreloadJob>(`/v1/preload/jobs/${id}`),
    cancel: (id: string) => this.fetch(`/v1/preload/jobs/${id}/cancel`, { method: "POST" }),
  };

  eviction = {
    simulate: (body: object) => this.fetch("/v1/eviction/simulate", { method: "POST", body: JSON.stringify(body) }),
    run: (body: object) => this.fetch("/v1/eviction/run", { method: "POST", body: JSON.stringify(body) }),
  };
}
