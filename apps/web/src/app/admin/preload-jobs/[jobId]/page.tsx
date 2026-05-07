"use client";
import { use, useState, useEffect } from "react";
import Link from "next/link";
import { api } from "@/lib/api-client";
import type { PreloadJob } from "@/types";
import { JobStatusBadge } from "@/components/ui/job-status-badge";
import { LoadingSpinner } from "@/components/ui/loading-spinner";

export default function JobDetailPage({ params }: { params: Promise<{ jobId: string }> }) {
  const { jobId } = use(params);
  const [job, setJob] = useState<PreloadJob | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [cancelling, setCancelling] = useState(false);

  function load() {
    setLoading(true);
    api.preloadJobs.get(jobId)
      .then(setJob)
      .catch((e: unknown) => setError(e instanceof Error ? e.message : String(e)))
      .finally(() => setLoading(false));
  }

  useEffect(() => { load(); }, [jobId]);

  async function handleCancel() {
    setCancelling(true);
    try { await api.preloadJobs.cancel(jobId); load(); }
    catch (e: unknown) { setError(e instanceof Error ? e.message : String(e)); }
    finally { setCancelling(false); }
  }

  if (loading) return <LoadingSpinner />;
  if (error) return <p className="text-red-500 text-sm">{error}</p>;
  if (!job) return <p className="text-gray-500 text-sm">Job not found.</p>;

  const canCancel = job.status === "PENDING" || job.status === "RUNNING";

  return (
    <div>
      <Link href="/admin/preload-jobs" className="text-indigo-600 text-sm hover:underline">← Preload Jobs</Link>
      <div className="flex items-center justify-between mt-2 mb-6">
        <h1 className="text-2xl font-bold text-gray-900">{job.name}</h1>
        {canCancel && (
          <button onClick={handleCancel} disabled={cancelling}
            className="bg-red-600 text-white px-4 py-2 rounded text-sm hover:bg-red-700 disabled:opacity-50">
            {cancelling ? "Cancelling…" : "Cancel Job"}
          </button>
        )}
      </div>
      <div className="bg-white rounded-lg shadow overflow-hidden">
        <dl className="divide-y divide-gray-100">
          {([
            ["ID", <span key="id" className="font-mono text-xs">{job.id}</span>],
            ["Status", <JobStatusBadge key="s" status={job.status} />],
            ["Site ID", <span key="si" className="font-mono text-xs">{job.site_id}</span>],
            ["Bandwidth Budget", job.bandwidth_budget_bytes != null
              ? `${(job.bandwidth_budget_bytes / 1024 / 1024).toFixed(0)} MB`
              : "Unlimited"],
            ["Started", job.started_at ? new Date(job.started_at).toLocaleString() : "—"],
            ["Completed", job.completed_at ? new Date(job.completed_at).toLocaleString() : "—"],
            ["Created", new Date(job.created_at).toLocaleString()],
          ] as [string, React.ReactNode][]).map(([label, value]) => (
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
