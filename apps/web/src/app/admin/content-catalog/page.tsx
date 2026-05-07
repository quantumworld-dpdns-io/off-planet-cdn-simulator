"use client";
import { useState } from "react";
import Link from "next/link";
import { useCacheObjects } from "@/lib/hooks";
import { LoadingSpinner } from "@/components/ui/loading-spinner";

function formatBytes(bytes: number): string {
  if (bytes < 1024 ** 2) return `${(bytes / 1024).toFixed(1)} KB`;
  if (bytes < 1024 ** 3) return `${(bytes / 1024 ** 2).toFixed(1)} MB`;
  return `${(bytes / 1024 ** 3).toFixed(2)} GB`;
}

export default function ContentCatalogPage() {
  const [statusFilter, setStatusFilter] = useState("");
  const { data: objects, loading, error } = useCacheObjects(
    statusFilter ? { status: statusFilter } : undefined
  );

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Content Catalog</h1>
        <select value={statusFilter} onChange={e => setStatusFilter(e.target.value)}
          className="border border-gray-300 rounded px-3 py-1.5 text-sm">
          <option value="">All statuses</option>
          <option value="ACTIVE">Active</option>
          <option value="DEPRECATED">Deprecated</option>
          <option value="DELETED">Deleted</option>
        </select>
      </div>
      {error && <p className="text-red-500 text-sm mb-4">{error}</p>}
      <div className="bg-white rounded-lg shadow overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>{["Name","Content Type","Size","Pinned","Status"].map(h => (
              <th key={h} className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">{h}</th>
            ))}</tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {loading ? (
              <tr><td colSpan={5}><LoadingSpinner /></td></tr>
            ) : objects?.length ? objects.map(obj => (
              <tr key={obj.id} className="hover:bg-gray-50">
                <td className="px-6 py-4 text-sm font-medium">
                  <Link href={`/admin/content-catalog/${obj.id}`} className="text-indigo-600 hover:underline">{obj.name}</Link>
                </td>
                <td className="px-6 py-4 text-sm text-gray-500">{obj.content_type ?? "—"}</td>
                <td className="px-6 py-4 text-sm text-gray-500">{formatBytes(obj.size_bytes)}</td>
                <td className="px-6 py-4 text-sm">{obj.pinned ? (
                  <span className="text-indigo-700 font-medium">Pinned</span>
                ) : <span className="text-gray-400">—</span>}</td>
                <td className="px-6 py-4 text-sm">
                  <span className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium ${
                    obj.status === "ACTIVE" ? "bg-green-100 text-green-800" :
                    obj.status === "DEPRECATED" ? "bg-yellow-100 text-yellow-700" :
                    "bg-gray-100 text-gray-500"
                  }`}>{obj.status}</span>
                </td>
              </tr>
            )) : (
              <tr><td colSpan={5} className="px-6 py-8 text-center text-gray-400 text-sm">No objects found.</td></tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
