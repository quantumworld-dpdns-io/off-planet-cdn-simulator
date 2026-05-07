"use client";
import Link from "next/link";
import { usePolicies } from "@/lib/hooks";
import { LoadingSpinner } from "@/components/ui/loading-spinner";

export default function CachePoliciesPage() {
  const { data: policies, loading, error } = usePolicies();

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Cache Policies</h1>
      {error && <p className="text-red-500 text-sm mb-4">{error}</p>}
      <div className="bg-white rounded-lg shadow overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>{["Name","Description","Enabled","Created"].map(h => (
              <th key={h} className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">{h}</th>
            ))}</tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {loading ? (
              <tr><td colSpan={4}><LoadingSpinner /></td></tr>
            ) : policies?.length ? policies.map(p => (
              <tr key={p.id} className="hover:bg-gray-50">
                <td className="px-6 py-4 text-sm font-medium">
                  <Link href={`/admin/cache-policies/${p.id}`} className="text-indigo-600 hover:underline">{p.name}</Link>
                </td>
                <td className="px-6 py-4 text-sm text-gray-500">{p.description ?? "—"}</td>
                <td className="px-6 py-4 text-sm">
                  <span className={`inline-flex px-2 py-0.5 rounded text-xs font-medium ${
                    p.enabled ? "bg-green-100 text-green-800" : "bg-gray-100 text-gray-500"
                  }`}>{p.enabled ? "Enabled" : "Disabled"}</span>
                </td>
                <td className="px-6 py-4 text-sm text-gray-500">{new Date(p.created_at).toLocaleDateString()}</td>
              </tr>
            )) : (
              <tr><td colSpan={4} className="px-6 py-8 text-center text-gray-400 text-sm">No policies configured.</td></tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
