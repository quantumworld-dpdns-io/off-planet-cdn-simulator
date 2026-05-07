"use client";
import { use, useState, useEffect } from "react";
import Link from "next/link";
import { api } from "@/lib/api-client";
import type { CacheObject } from "@/types";
import { LoadingSpinner } from "@/components/ui/loading-spinner";

function formatBytes(bytes: number): string {
  if (bytes < 1024 ** 2) return `${(bytes / 1024).toFixed(1)} KB`;
  if (bytes < 1024 ** 3) return `${(bytes / 1024 ** 2).toFixed(1)} MB`;
  return `${(bytes / 1024 ** 3).toFixed(2)} GB`;
}

export default function ObjectDetailPage({ params }: { params: Promise<{ objectId: string }> }) {
  const { objectId } = use(params);
  const [obj, setObj] = useState<CacheObject | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [toggling, setToggling] = useState(false);

  function load() {
    setLoading(true);
    api.cacheObjects.get(objectId)
      .then(setObj)
      .catch((e: unknown) => setError(e instanceof Error ? e.message : String(e)))
      .finally(() => setLoading(false));
  }

  useEffect(() => { load(); }, [objectId]);

  async function togglePin() {
    if (!obj) return;
    setToggling(true);
    try {
      if (obj.pinned) await api.cacheObjects.unpin(objectId);
      else await api.cacheObjects.pin(objectId);
      load();
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : String(e));
    } finally {
      setToggling(false);
    }
  }

  if (loading) return <LoadingSpinner />;
  if (error) return <p className="text-red-500 text-sm">{error}</p>;
  if (!obj) return <p className="text-gray-500 text-sm">Object not found.</p>;

  return (
    <div>
      <Link href="/admin/content-catalog" className="text-indigo-600 text-sm hover:underline">← Content Catalog</Link>
      <div className="flex items-center justify-between mt-2 mb-6">
        <h1 className="text-2xl font-bold text-gray-900">{obj.name}</h1>
        <button onClick={togglePin} disabled={toggling}
          className={`px-4 py-2 rounded text-sm font-medium disabled:opacity-50 ${
            obj.pinned ? "bg-yellow-100 text-yellow-800 hover:bg-yellow-200" : "bg-indigo-600 text-white hover:bg-indigo-700"
          }`}>
          {toggling ? "…" : obj.pinned ? "Unpin" : "Pin"}
        </button>
      </div>

      <div className="bg-white rounded-lg shadow overflow-hidden">
        <dl className="divide-y divide-gray-100">
          {([
            ["ID", <span key="id" className="font-mono text-xs">{obj.id}</span>],
            ["Status", obj.status],
            ["Content Type", obj.content_type ?? "—"],
            ["Size", formatBytes(obj.size_bytes)],
            ["Pinned", obj.pinned ? "Yes" : "No"],
            ["Source URL", obj.source_url ? (
              <a key="url" href={obj.source_url} className="text-indigo-600 hover:underline break-all">{obj.source_url}</a>
            ) : "—"],
            ["Priority ID", <span key="p" className="font-mono text-xs">{obj.priority_class_id}</span>],
            ["Tags", obj.tags.join(", ") || "—"],
            ["Created", new Date(obj.created_at).toLocaleString()],
          ] as [string, React.ReactNode][]).map(([label, value]) => (
            <div key={label} className="px-6 py-4 flex items-start">
              <dt className="text-sm font-medium text-gray-500 w-48 shrink-0">{label}</dt>
              <dd className="text-sm text-gray-900 flex-1">{value}</dd>
            </div>
          ))}
        </dl>
      </div>
    </div>
  );
}
