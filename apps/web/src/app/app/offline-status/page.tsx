"use client";
import { useNodes, useBandwidthWindows } from "@/lib/hooks";
import { CacheFillBar } from "@/components/ui/cache-fill-bar";
import { LoadingSpinner } from "@/components/ui/loading-spinner";

export default function OfflineStatusPage() {
  const { data: nodes, loading: nodesLoading } = useNodes();
  const { data: windows, loading: windowsLoading } = useBandwidthWindows();

  const totalUsed = nodes?.reduce((s, n) => s + n.cache_used_bytes, 0) ?? 0;
  const totalMax = nodes?.reduce((s, n) => s + n.cache_max_bytes, 0) ?? 0;
  const onlineNodes = nodes?.filter(n => n.status === "ONLINE").length ?? 0;

  const lastSync = nodes
    ?.filter(n => n.last_seen)
    .sort((a, b) => new Date(b.last_seen!).getTime() - new Date(a.last_seen!).getTime())[0]?.last_seen;

  const nextWindow = windows
    ?.filter(w => new Date(w.window_start) > new Date())
    .sort((a, b) => new Date(a.window_start).getTime() - new Date(b.window_start).getTime())[0];

  if (nodesLoading || windowsLoading) return <LoadingSpinner />;

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Offline Status</h1>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
        <div className="bg-white rounded-lg shadow p-6">
          <p className="text-sm text-gray-500">Nodes Online</p>
          <p className="text-3xl font-bold text-gray-900 mt-1">{onlineNodes} / {nodes?.length ?? 0}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-6">
          <p className="text-sm text-gray-500">Last Sync</p>
          <p className="text-lg font-semibold text-gray-900 mt-1">{lastSync ? new Date(lastSync).toLocaleString() : "Never"}</p>
        </div>
        <div className="bg-white rounded-lg shadow p-6">
          <p className="text-sm text-gray-500">Next Contact Window</p>
          <p className="text-lg font-semibold text-gray-900 mt-1">
            {nextWindow ? new Date(nextWindow.window_start).toLocaleString() : "None scheduled"}
          </p>
          {nextWindow && (
            <p className="text-xs text-gray-400 mt-1">{nextWindow.label ?? ""} · {(nextWindow.bandwidth_bps / 1_000_000).toFixed(1)} Mbps</p>
          )}
        </div>
      </div>

      <h2 className="text-base font-semibold text-gray-700 mb-3">Cache Storage</h2>
      <div className="bg-white rounded-lg shadow p-6">
        <div className="mb-4">
          <p className="text-sm text-gray-500 mb-2">Aggregate fill across all nodes</p>
          <CacheFillBar usedBytes={totalUsed} maxBytes={totalMax} />
        </div>
        <div className="divide-y divide-gray-100 mt-4">
          {nodes?.map(n => (
            <div key={n.id} className="py-3 flex items-center gap-4">
              <div className="w-32 shrink-0">
                <p className="text-sm font-medium text-gray-700 truncate">{n.name}</p>
                <span className={`text-xs font-medium ${n.status === "ONLINE" ? "text-green-600" : "text-red-500"}`}>{n.status}</span>
              </div>
              <div className="flex-1">
                <CacheFillBar usedBytes={n.cache_used_bytes} maxBytes={n.cache_max_bytes} />
              </div>
            </div>
          ))}
          {!nodes?.length && <p className="text-sm text-gray-400 py-4 text-center">No nodes registered.</p>}
        </div>
      </div>
    </div>
  );
}
