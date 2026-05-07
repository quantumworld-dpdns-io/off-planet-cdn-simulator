"use client";
import { usePreloadJobs } from "@/lib/hooks";
import { JobStatusBadge } from "@/components/ui/job-status-badge";
import { LoadingSpinner } from "@/components/ui/loading-spinner";

export default function DownloadsPage() {
  const { data: jobs, loading, error } = usePreloadJobs();

  const active = jobs?.filter(j => j.status === "PENDING" || j.status === "RUNNING") ?? [];
  const completed = jobs?.filter(j => j.status === "DONE" || j.status === "FAILED" || j.status === "CANCELLED") ?? [];

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Downloads</h1>
      {error && <p className="text-red-500 text-sm mb-4">{error}</p>}
      {loading ? <LoadingSpinner /> : (
        <>
          <section className="mb-8">
            <h2 className="text-base font-semibold text-gray-700 mb-3">In Progress</h2>
            <div className="bg-white rounded-lg shadow divide-y divide-gray-100">
              {active.length ? active.map(job => (
                <div key={job.id} className="px-6 py-4 flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-gray-900">{job.name}</p>
                    {job.bandwidth_budget_bytes != null && (
                      <p className="text-xs text-gray-400 mt-0.5">Budget: {(job.bandwidth_budget_bytes / 1024 / 1024).toFixed(0)} MB</p>
                    )}
                  </div>
                  <JobStatusBadge status={job.status} />
                </div>
              )) : (
                <p className="px-6 py-8 text-sm text-gray-400 text-center">No active downloads.</p>
              )}
            </div>
          </section>

          <section>
            <h2 className="text-base font-semibold text-gray-700 mb-3">Completed</h2>
            <div className="bg-white rounded-lg shadow divide-y divide-gray-100">
              {completed.length ? completed.map(job => (
                <div key={job.id} className="px-6 py-4 flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-gray-900">{job.name}</p>
                    <p className="text-xs text-gray-400 mt-0.5">
                      {job.completed_at ? `Completed ${new Date(job.completed_at).toLocaleString()}` : new Date(job.created_at).toLocaleDateString()}
                    </p>
                  </div>
                  <JobStatusBadge status={job.status} />
                </div>
              )) : (
                <p className="px-6 py-8 text-sm text-gray-400 text-center">No completed downloads.</p>
              )}
            </div>
          </section>
        </>
      )}
    </div>
  );
}
