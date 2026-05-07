"use client";
import Link from "next/link";
import { useNodes } from "@/lib/hooks";
import { StatusBadge } from "@/components/ui/status-badge";
import { CacheFillBar } from "@/components/ui/cache-fill-bar";
import { LoadingSpinner } from "@/components/ui/loading-spinner";

export default function NodesPage() {
  const { data: nodes, loading, error } = useNodes();

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Edge Nodes</h1>
      {error && <p className="text-red-500 text-sm mb-4">{error}</p>}
      <div className="bg-white rounded-lg shadow overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>{["Name","Status","Cache Fill","Last Heartbeat"].map(h => (
              <th key={h} className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">{h}</th>
            ))}</tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {loading ? (
              <tr><td colSpan={4}><LoadingSpinner /></td></tr>
            ) : nodes?.length ? nodes.map(n => (
              <tr key={n.id} className="hover:bg-gray-50">
                <td className="px-6 py-4 text-sm font-medium">
                  <Link href={`/admin/nodes/${n.id}`} className="text-indigo-600 hover:underline">{n.name}</Link>
                </td>
                <td className="px-6 py-4 text-sm"><StatusBadge status={n.status} /></td>
                <td className="px-6 py-4 text-sm w-56"><CacheFillBar usedBytes={n.cache_used_bytes} maxBytes={n.cache_max_bytes} /></td>
                <td className="px-6 py-4 text-sm text-gray-500">{n.last_seen ? new Date(n.last_seen).toLocaleString() : "Never"}</td>
              </tr>
            )) : (
              <tr><td colSpan={4} className="px-6 py-8 text-center text-gray-400 text-sm">No nodes registered yet.</td></tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
