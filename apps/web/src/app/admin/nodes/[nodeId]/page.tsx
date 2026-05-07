"use client";
import { use, useState, useEffect } from "react";
import Link from "next/link";
import { api } from "@/lib/api-client";
import type { Node } from "@/types";
import { StatusBadge } from "@/components/ui/status-badge";
import { CacheFillBar } from "@/components/ui/cache-fill-bar";
import { LoadingSpinner } from "@/components/ui/loading-spinner";

export default function NodeDetailPage({ params }: { params: Promise<{ nodeId: string }> }) {
  const { nodeId } = use(params);
  const [node, setNode] = useState<Node | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    api.nodes.get(nodeId)
      .then(setNode)
      .catch((e: unknown) => setError(e instanceof Error ? e.message : String(e)))
      .finally(() => setLoading(false));
  }, [nodeId]);

  if (loading) return <LoadingSpinner />;
  if (error) return <p className="text-red-500 text-sm">{error}</p>;
  if (!node) return <p className="text-gray-500 text-sm">Node not found.</p>;

  const rows: [string, React.ReactNode][] = [
    ["ID", <span key="id" className="font-mono text-xs">{node.id}</span>],
    ["Status", <StatusBadge key="s" status={node.status} />],
    ["Site ID", <span key="si" className="font-mono text-xs">{node.site_id}</span>],
    ["Cache Directory", node.cache_dir],
    ["Cache Fill", <CacheFillBar key="cf" usedBytes={node.cache_used_bytes} maxBytes={node.cache_max_bytes} />],
    ["Last Heartbeat", node.last_seen ? new Date(node.last_seen).toLocaleString() : "Never"],
    ["Created", new Date(node.created_at).toLocaleString()],
  ];

  return (
    <div>
      <Link href="/admin/nodes" className="text-indigo-600 text-sm hover:underline">← Nodes</Link>
      <h1 className="text-2xl font-bold text-gray-900 mt-2 mb-6">{node.name}</h1>
      <div className="bg-white rounded-lg shadow overflow-hidden">
        <dl className="divide-y divide-gray-100">
          {rows.map(([label, value]) => (
            <div key={label} className="px-6 py-4 flex items-center">
              <dt className="text-sm font-medium text-gray-500 w-48">{label}</dt>
              <dd className="text-sm text-gray-900 flex-1">{value}</dd>
            </div>
          ))}
        </dl>
      </div>
    </div>
  );
}
