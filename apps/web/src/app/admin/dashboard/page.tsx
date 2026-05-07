"use client";
import { useSites } from "@/lib/hooks";
import { useNodes } from "@/lib/hooks";
import { usePreloadJobs } from "@/lib/hooks";
import { LoadingSpinner } from "@/components/ui/loading-spinner";

export default function DashboardPage() {
  const { data: sites, loading: sitesLoading } = useSites();
  const { data: nodes, loading: nodesLoading } = useNodes();
  const { data: jobs, loading: jobsLoading } = usePreloadJobs();

  const loading = sitesLoading || nodesLoading || jobsLoading;

  const onlineNodes = nodes?.filter(n => n.status === "ONLINE").length ?? 0;
  const totalNodes = nodes?.length ?? 0;

  const totalUsed = nodes?.reduce((s, n) => s + n.cache_used_bytes, 0) ?? 0;
  const totalMax = nodes?.reduce((s, n) => s + n.cache_max_bytes, 0) ?? 1;
  const fillPct = totalMax > 0 ? Math.round((totalUsed / totalMax) * 100) : 0;

  const activeJobs = jobs?.filter(j => j.status === "PENDING" || j.status === "RUNNING").length ?? 0;

  const cards = [
    { label: "Nodes Online", value: loading ? "—" : `${onlineNodes} / ${totalNodes}` },
    { label: "Cache Fill", value: loading ? "—" : `${fillPct}%` },
    { label: "Active Jobs", value: loading ? "—" : String(activeJobs) },
    { label: "Sites", value: loading ? "—" : String(sites?.length ?? 0) },
  ];

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Admin Dashboard</h1>
      {loading ? (
        <LoadingSpinner />
      ) : (
        <>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
            {cards.map(card => (
              <div key={card.label} className="bg-white rounded-lg shadow p-6">
                <p className="text-sm text-gray-500">{card.label}</p>
                <p className="text-3xl font-bold text-gray-900 mt-1">{card.value}</p>
              </div>
            ))}
          </div>
          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-base font-semibold text-gray-800 mb-4">Node Health</h2>
            <div className="space-y-2">
              {nodes?.slice(0, 5).map(node => (
                <div key={node.id} className="flex items-center justify-between text-sm">
                  <span className="text-gray-700">{node.name}</span>
                  <span className={`font-medium ${node.status === "ONLINE" ? "text-green-600" : "text-red-500"}`}>
                    {node.status}
                  </span>
                </div>
              ))}
              {!nodes?.length && <p className="text-gray-400 text-sm">No nodes registered.</p>}
            </div>
          </div>
        </>
      )}
    </div>
  );
}
