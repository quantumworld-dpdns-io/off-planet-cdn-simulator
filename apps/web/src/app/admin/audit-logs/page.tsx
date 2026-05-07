"use client";
import { useAuditLogs } from "@/lib/hooks";
import { LoadingSpinner } from "@/components/ui/loading-spinner";

export default function AuditLogsPage() {
  const { data: logs, loading, error } = useAuditLogs({ limit: 100 });

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Audit Logs</h1>
      {error && <p className="text-red-500 text-sm mb-4">{error}</p>}
      <div className="bg-white rounded-lg shadow overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>{["Action","Resource","Actor","Timestamp"].map(h => (
              <th key={h} className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">{h}</th>
            ))}</tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {loading ? (
              <tr><td colSpan={4}><LoadingSpinner /></td></tr>
            ) : logs?.length ? logs.map(log => (
              <tr key={log.id} className="hover:bg-gray-50">
                <td className="px-6 py-4 text-sm font-mono text-gray-700">{log.action}</td>
                <td className="px-6 py-4 text-sm text-gray-500">
                  {log.resource_type}{log.resource_id ? ` · ${log.resource_id.slice(0, 8)}…` : ""}
                </td>
                <td className="px-6 py-4 text-sm font-mono text-gray-500">
                  {log.actor_id ? log.actor_id.slice(0, 8) + "…" : "system"}
                </td>
                <td className="px-6 py-4 text-sm text-gray-500">{new Date(log.created_at).toLocaleString()}</td>
              </tr>
            )) : (
              <tr><td colSpan={4} className="px-6 py-8 text-center text-gray-400 text-sm">No audit logs yet.</td></tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
