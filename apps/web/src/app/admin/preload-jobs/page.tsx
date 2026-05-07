"use client";
import { useState } from "react";
import Link from "next/link";
import { usePreloadJobs } from "@/lib/hooks";
import { JobStatusBadge } from "@/components/ui/job-status-badge";
import { LoadingSpinner } from "@/components/ui/loading-spinner";
import type { JobStatus } from "@/types";

export default function PreloadJobsPage() {
  const [statusFilter, setStatusFilter] = useState<JobStatus | "">("");
  const { data: jobs, loading, error } = usePreloadJobs(
    statusFilter ? { status: statusFilter } : undefined
  );

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Preload Jobs</h1>
        <select value={statusFilter} onChange={e => setStatusFilter(e.target.value as JobStatus | "")}
          className="border border-gray-300 rounded px-3 py-1.5 text-sm">
          <option value="">All statuses</option>
          {(["PENDING","RUNNING","DONE","FAILED","CANCELLED"] as JobStatus[]).map(s => (
            <option key={s} value={s}>{s}</option>
          ))}
        </select>
      </div>
      {error && <p className="text-red-500 text-sm mb-4">{error}</p>}
      <div className="bg-white rounded-lg shadow overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>{["Name","Status","Bandwidth Budget","Created"].map(h => (
              <th key={h} className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">{h}</th>
            ))}</tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {loading ? (
              <tr><td colSpan={4}><LoadingSpinner /></td></tr>
            ) : jobs?.length ? jobs.map(j => (
              <tr key={j.id} className="hover:bg-gray-50">
                <td className="px-6 py-4 text-sm font-medium">
                  <Link href={`/admin/preload-jobs/${j.id}`} className="text-indigo-600 hover:underline">{j.name}</Link>
                </td>
                <td className="px-6 py-4 text-sm"><JobStatusBadge status={j.status} /></td>
                <td className="px-6 py-4 text-sm text-gray-500">
                  {j.bandwidth_budget_bytes != null ? `${(j.bandwidth_budget_bytes / 1024 / 1024).toFixed(0)} MB` : "Unlimited"}
                </td>
                <td className="px-6 py-4 text-sm text-gray-500">{new Date(j.created_at).toLocaleDateString()}</td>
              </tr>
            )) : (
              <tr><td colSpan={4} className="px-6 py-8 text-center text-gray-400 text-sm">No preload jobs found.</td></tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
