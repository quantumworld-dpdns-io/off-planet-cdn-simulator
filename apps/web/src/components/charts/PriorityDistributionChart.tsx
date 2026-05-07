"use client";

import {
  PieChart,
  Pie,
  Cell,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import type { PriorityBucket } from "@/lib/api-client";

interface Props {
  buckets: PriorityBucket[];
}

const COLOR_MAP: Record<string, string> = {
  P0: "#dc2626",
  P1: "#ea580c",
  P2: "#2563eb",
  P3: "#7c3aed",
  P4: "#16a34a",
  P5: "#6b7280",
  UNKNOWN: "#9ca3af",
};

function formatBytes(bytes: number): string {
  const GB = 1024 * 1024 * 1024;
  const MB = 1024 * 1024;
  const KB = 1024;
  if (bytes >= GB) return `${(bytes / GB).toFixed(1)} GB`;
  if (bytes >= MB) return `${(bytes / MB).toFixed(1)} MB`;
  if (bytes >= KB) return `${(bytes / KB).toFixed(1)} KB`;
  return `${bytes} B`;
}

interface TooltipPayloadEntry {
  name: string;
  value: number;
  payload: PriorityBucket;
}

interface CustomTooltipProps {
  active?: boolean;
  payload?: TooltipPayloadEntry[];
}

function CustomTooltip({ active, payload }: CustomTooltipProps) {
  if (!active || !payload || payload.length === 0) return null;
  const entry = payload[0];
  const bucket = entry.payload;
  return (
    <div className="bg-white border border-gray-200 rounded shadow p-2 text-xs">
      <p className="font-semibold">{bucket.level}</p>
      <p>Count: {bucket.count}</p>
      <p>Total: {formatBytes(bucket.total_bytes)}</p>
    </div>
  );
}

export function PriorityDistributionChart({ buckets }: Props) {
  if (buckets.length === 0) {
    return (
      <p className="text-gray-400 text-sm text-center py-8">
        No active objects
      </p>
    );
  }

  return (
    <ResponsiveContainer width="100%" height={220}>
      <PieChart>
        <Pie
          data={buckets}
          dataKey="count"
          nameKey="level"
          outerRadius={80}
        >
          {buckets.map((bucket) => (
            <Cell
              key={bucket.level}
              fill={COLOR_MAP[bucket.level] ?? COLOR_MAP["UNKNOWN"]}
            />
          ))}
        </Pie>
        <Tooltip content={<CustomTooltip />} />
        <Legend />
      </PieChart>
    </ResponsiveContainer>
  );
}
