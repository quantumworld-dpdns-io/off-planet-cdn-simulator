"use client";

import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import type { NodeFill } from "@/lib/api-client";

interface Props {
  nodes: NodeFill[];
}

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
  payload: { node_name: string; used: number; free: number; max_bytes: number };
}

interface CustomTooltipProps {
  active?: boolean;
  payload?: TooltipPayloadEntry[];
  label?: string;
}

function CustomTooltip({ active, payload }: CustomTooltipProps) {
  if (!active || !payload || payload.length === 0) return null;
  const d = payload[0].payload;
  const fillPct =
    d.max_bytes > 0 ? Math.round((d.used / d.max_bytes) * 100) : 0;
  return (
    <div className="bg-white border border-gray-200 rounded shadow p-2 text-xs">
      <p className="font-semibold">{d.node_name}</p>
      <p>Used: {formatBytes(d.used)}</p>
      <p>Fill: {fillPct}%</p>
    </div>
  );
}

export function NodeFillChart({ nodes }: Props) {
  if (nodes.length === 0) {
    return (
      <p className="text-gray-400 text-sm text-center py-8">No online nodes</p>
    );
  }

  const data = nodes.map((n) => ({
    node_name: n.node_name,
    label:
      n.node_name.length > 12 ? n.node_name.slice(0, 12) + "…" : n.node_name,
    used: n.used_bytes,
    free: n.max_bytes - n.used_bytes,
    max_bytes: n.max_bytes,
  }));

  return (
    <ResponsiveContainer width="100%" height={200}>
      <BarChart data={data} margin={{ top: 4, right: 8, left: 0, bottom: 0 }}>
        <XAxis dataKey="label" tick={{ fontSize: 11 }} />
        <YAxis tickFormatter={formatBytes} tick={{ fontSize: 11 }} width={60} />
        <Tooltip content={<CustomTooltip />} />
        <Legend />
        <Bar dataKey="used" fill="#2563eb" stackId="a" name="used" />
        <Bar dataKey="free" fill="#e5e7eb" stackId="a" name="free" />
      </BarChart>
    </ResponsiveContainer>
  );
}
