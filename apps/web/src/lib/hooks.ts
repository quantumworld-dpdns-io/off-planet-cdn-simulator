"use client";
import { useState, useEffect, useCallback } from "react";
import { api } from "./api-client";
import type { Site, Node, CacheObject, CachePolicy, PreloadJob, AuditLog, BandwidthWindow } from "@/types";

type HookResult<T> = { data: T | null; loading: boolean; error: string | null; refetch: () => void };

function useApiData<T>(fetcher: () => Promise<T>): HookResult<T> {
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetch = useCallback(() => {
    setLoading(true);
    setError(null);
    fetcher()
      .then(setData)
      .catch((e: unknown) => setError(e instanceof Error ? e.message : String(e)))
      .finally(() => setLoading(false));
  }, [fetcher]);

  useEffect(() => { fetch(); }, [fetch]);

  return { data, loading, error, refetch: fetch };
}

export function useSites() {
  const fetcher = useCallback(() => api.sites.list().then(r => r.sites), []);
  return useApiData<Site[]>(fetcher);
}

export function useNodes(siteId?: string) {
  const fetcher = useCallback(
    () => api.nodes.list(siteId ? { site_id: siteId } : undefined).then(r => r.nodes),
    [siteId]
  );
  return useApiData<Node[]>(fetcher);
}

export function useCacheObjects(params?: { site_id?: string; priority?: string; status?: string }) {
  const key = JSON.stringify(params);
  const fetcher = useCallback(
    () => api.cacheObjects.list(params).then(r => r.objects),
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [key]
  );
  return useApiData<CacheObject[]>(fetcher);
}

export function usePolicies(siteId?: string) {
  const fetcher = useCallback(
    () => api.policies.list(siteId ? { site_id: siteId } : undefined).then(r => r.policies),
    [siteId]
  );
  return useApiData<CachePolicy[]>(fetcher);
}

export function usePreloadJobs(params?: { site_id?: string; status?: string }) {
  const key = JSON.stringify(params);
  const fetcher = useCallback(
    () => api.preloadJobs.list(params).then(r => r.jobs),
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [key]
  );
  return useApiData<PreloadJob[]>(fetcher);
}

export function useAuditLogs(params?: { resource_type?: string; limit?: number }) {
  const key = JSON.stringify(params);
  const fetcher = useCallback(
    () => api.auditLogs.list(params).then(r => r.logs),
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [key]
  );
  return useApiData<AuditLog[]>(fetcher);
}

export function useBandwidthWindows(siteId?: string) {
  const fetcher = useCallback(
    () => api.bandwidthWindows.list(siteId ? { site_id: siteId } : undefined).then(r => r.windows),
    [siteId]
  );
  return useApiData<BandwidthWindow[]>(fetcher);
}
