"use client";
import { use } from "react";
import Link from "next/link";
import { useNodes } from "@/lib/hooks";
import { usePreloadJobs } from "@/lib/hooks";
import { StatusBadge } from "@/components/ui/status-badge";
import { CacheFillBar } from "@/components/ui/cache-fill-bar";
import { JobStatusBadge } from "@/components/ui/job-status-badge";
import { LoadingSpinner } from "@/components/ui/loading-spinner";

export default function SiteDetailPage({ params }: { params: Promise<{ siteId: string }> }) {
  const { siteId } = use(params);
  const { data: nodes, loading: nodesLoading } = useNodes(siteId);
  const { data: jobs, loading: jobsLoading } = usePreloadJobs({ site_id: siteId });

  return (
    <div>
      <div className="mb-6">
        <Link href="/admin/sites" className="text-indigo-600 text-sm hover:underline">← Sites</Link>
        <h1 className="text-2xl font-bold text-gray-900 mt-2">Site Detail</h1>
        <p className="text-xs text-gray-400 mt-1 font-mono">{siteId}</p>
      </div>

      <section className="mb-8">
        <h2 className="text-base font-semibold text-gray-700 mb-3">Nodes</h2>
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>{["Name","Status","Cache Fill","Last Heartbeat"].map(h => (
                <th key={h} className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">{h}</th>
              ))}</tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {nodesLoading ? (
                <tr><td colSpan={4}><LoadingSpinner /></td></tr>
              ) : nodes?.length ? nodes.map(n => (
                <tr key={n.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 text-sm font-medium">
                    <Link href={`/admin/nodes/${n.id}`} className="text-indigo-600 hover:underline">{n.name}</Link>
                  </td>
                  <td className="px-6 py-4 text-sm"><StatusBadge status={n.status} /></td>
                  <td className="px-6 py-4 text-sm w-48"><CacheFillBar usedBytes={n.cache_used_bytes} maxBytes={n.cache_max_bytes} /></td>
                  <td className="px-6 py-4 text-sm text-gray-500">{n.last_seen ? new Date(n.last_seen).toLocaleString() : "Never"}</td>
                </tr>
              )) : (
                <tr><td colSpan={4} className="px-6 py-8 text-center text-gray-400 text-sm">No nodes for this site.</td></tr>
              )}
            </tbody>
          </table>
        </div>
      </section>

      <section>
        <h2 className="text-base font-semibold text-gray-700 mb-3">Preload Jobs</h2>
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>{["Name","Status","Created"].map(h => (
                <th key={h} className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">{h}</th>
              ))}</tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {jobsLoading ? (
                <tr><td colSpan={3}><LoadingSpinner /></td></tr>
              ) : jobs?.length ? jobs.map(j => (
                <tr key={j.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 text-sm font-medium">
                    <Link href={`/admin/preload-jobs/${j.id}`} className="text-indigo-600 hover:underline">{j.name}</Link>
                  </td>
                  <td className="px-6 py-4 text-sm"><JobStatusBadge status={j.status} /></td>
                  <td className="px-6 py-4 text-sm text-gray-500">{new Date(j.created_at).toLocaleDateString()}</td>
                </tr>
              )) : (
                <tr><td colSpan={3} className="px-6 py-8 text-center text-gray-400 text-sm">No preload jobs for this site.</td></tr>
              )}
            </tbody>
          </table>
        </div>
      </section>
    </div>
  );
}
