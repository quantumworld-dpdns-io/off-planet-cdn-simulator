export type PriorityClassLevel = "P0" | "P1" | "P2" | "P3" | "P4" | "P5";
export type NodeStatus = "ONLINE" | "OFFLINE" | "DEGRADED" | "UNKNOWN";
export type JobStatus = "PENDING" | "RUNNING" | "DONE" | "FAILED" | "CANCELLED";

export interface Site {
  id: string;
  org_id: string;
  name: string;
  location?: string;
  description?: string;
  created_at: string;
  updated_at: string;
}

export interface Node {
  id: string;
  org_id: string;
  site_id: string;
  name: string;
  status: NodeStatus;
  cache_dir: string;
  cache_max_bytes: number;
  cache_used_bytes: number;
  last_seen?: string;
  created_at: string;
  updated_at: string;
}

export interface CacheObject {
  id: string;
  org_id: string;
  site_id: string;
  priority_class_id: string;
  name: string;
  content_type?: string;
  source_url?: string;
  content_hash?: string;
  size_bytes: number;
  pinned: boolean;
  status: "ACTIVE" | "DEPRECATED" | "DELETED";
  tags: string[];
  metadata: Record<string, unknown>;
  created_at: string;
  updated_at: string;
}

export interface CachePolicy {
  id: string;
  org_id: string;
  site_id: string;
  name: string;
  description?: string;
  enabled: boolean;
  created_at: string;
  updated_at: string;
}

export interface PreloadJob {
  id: string;
  org_id: string;
  site_id: string;
  name: string;
  status: JobStatus;
  bandwidth_budget_bytes?: number;
  started_at?: string;
  completed_at?: string;
  created_at: string;
  updated_at: string;
}

export interface AuditLog {
  id: string;
  org_id: string;
  actor_id?: string;
  action: string;
  resource_type: string;
  resource_id?: string;
  created_at: string;
}

export interface BandwidthWindow {
  id: string;
  org_id: string;
  site_id: string;
  label?: string;
  window_start: string;
  window_end: string;
  bandwidth_bps: number;
  reliability_score: number;
  created_at: string;
}
