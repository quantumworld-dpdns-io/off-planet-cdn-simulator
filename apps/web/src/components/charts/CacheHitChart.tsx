"use client";

import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import type { CacheHitPoint } from "@/lib/api-client";

interface Props {
  points: CacheHitPoint[];
}

function formatHour(hour: string): string {
  return new Date(hour).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
}

export function CacheHitChart({ points }: Props) {
  if (points.length === 0) {
    return (
      <p className="text-gray-400 text-sm text-center py-8">
        No cache events in the last 24h
      </p>
    );
  }

  const data = points.map((p) => ({
    ...p,
    label: formatHour(p.hour),
  }));

  return (
    <ResponsiveContainer width="100%" height={200}>
      <AreaChart data={data} margin={{ top: 4, right: 8, left: 0, bottom: 0 }}>
        <XAxis dataKey="label" tick={{ fontSize: 11 }} />
        <YAxis tick={{ fontSize: 11 }} />
        <Tooltip />
        <Legend />
        <Area
          type="monotone"
          dataKey="hits"
          fill="#16a34a"
          stroke="#15803d"
          stackId="1"
        />
        <Area
          type="monotone"
          dataKey="misses"
          fill="#dc2626"
          stroke="#b91c1c"
          stackId="2"
        />
      </AreaChart>
    </ResponsiveContainer>
  );
}
