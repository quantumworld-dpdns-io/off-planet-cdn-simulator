"use client";
import Link from "next/link";
import { useNodes, useCacheObjects } from "@/lib/hooks";
import { LoadingSpinner } from "@/components/ui/loading-spinner";

const categories = [
  { href: "/app/medical", label: "Medical References", badge: "P0", bg: "bg-red-50", border: "border-red-200", badgeColor: "bg-red-100 text-red-800" },
  { href: "/app/engineering", label: "Engineering Manuals", badge: "P1", bg: "bg-orange-50", border: "border-orange-200", badgeColor: "bg-orange-100 text-orange-800" },
  { href: "/app/manuals", label: "Manuals", badge: "P1", bg: "bg-orange-50", border: "border-orange-200", badgeColor: "bg-orange-100 text-orange-800" },
  { href: "/app/education", label: "Education", badge: "P2", bg: "bg-yellow-50", border: "border-yellow-200", badgeColor: "bg-yellow-100 text-yellow-800" },
  { href: "/app/entertainment", label: "Entertainment", badge: "P4", bg: "bg-gray-50", border: "border-gray-200", badgeColor: "bg-gray-100 text-gray-700" },
];

function formatBytes(bytes: number): string {
  if (bytes < 1024 ** 2) return `${(bytes / 1024).toFixed(1)} KB`;
  if (bytes < 1024 ** 3) return `${(bytes / 1024 ** 2).toFixed(1)} MB`;
  return `${(bytes / 1024 ** 3).toFixed(2)} GB`;
}

export default function UserHomePage() {
  const { data: nodes, loading: nodesLoading } = useNodes();
  const { data: objects, loading: objLoading } = useCacheObjects({ status: "ACTIVE" });

  const onlineCount = nodes?.filter(n => n.status === "ONLINE").length ?? 0;
  const totalNodes = nodes?.length ?? 0;
  const isConnected = onlineCount > 0;

  const recent = objects?.slice(0, 6) ?? [];

  return (
    <div>
      {/* Connection status banner */}
      {!nodesLoading && (
        <div className={`rounded-lg p-4 mb-6 border ${isConnected ? "bg-green-50 border-green-200" : "bg-yellow-50 border-yellow-200"}`}>
          <p className={`text-sm font-medium ${isConnected ? "text-green-800" : "text-yellow-800"}`}>
            {isConnected ? `Connected — ${onlineCount}/${totalNodes} nodes online` : "Offline mode — serving from local cache"}
          </p>
          {nodes?.some(n => n.last_seen) && (
            <p className="text-xs text-gray-500 mt-1">
              Last sync: {new Date(nodes!.filter(n => n.last_seen).sort((a, b) => new Date(b.last_seen!).getTime() - new Date(a.last_seen!).getTime())[0].last_seen!).toLocaleString()}
            </p>
          )}
        </div>
      )}

      <h1 className="text-2xl font-bold text-gray-900 mb-6">Off-Planet CDN</h1>

      {/* Category grid */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-10">
        {categories.map(cat => (
          <Link key={cat.href} href={cat.href}
            className={`${cat.bg} border ${cat.border} rounded-lg p-6 hover:shadow-md transition-shadow`}>
            <span className={`inline-block text-xs font-medium px-2 py-0.5 rounded-full ${cat.badgeColor} mb-3`}>{cat.badge}</span>
            <h2 className="font-semibold text-gray-900">{cat.label}</h2>
          </Link>
        ))}
      </div>

      {/* Recently cached */}
      <h2 className="text-base font-semibold text-gray-800 mb-3">Recently Cached</h2>
      {objLoading ? <LoadingSpinner /> : (
        <div className="bg-white rounded-lg shadow divide-y divide-gray-100">
          {recent.length ? recent.map(obj => (
            <div key={obj.id} className="px-6 py-4 flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-900">{obj.name}</p>
                <p className="text-xs text-gray-400 mt-0.5">{obj.content_type ?? "unknown type"} · {formatBytes(obj.size_bytes)}</p>
              </div>
              {obj.source_url && (
                <a href={obj.source_url} target="_blank" rel="noreferrer"
                  className="text-xs text-indigo-600 hover:underline ml-4 shrink-0">Open ↗</a>
              )}
            </div>
          )) : (
            <p className="px-6 py-8 text-sm text-gray-400 text-center">No cached content yet.</p>
          )}
        </div>
      )}
    </div>
  );
}
