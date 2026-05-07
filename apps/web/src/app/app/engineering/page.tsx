"use client";
import { useCacheObjects } from "@/lib/hooks";
import { LoadingSpinner } from "@/components/ui/loading-spinner";

function formatBytes(bytes: number): string {
  if (bytes < 1024 ** 2) return `${(bytes / 1024).toFixed(1)} KB`;
  if (bytes < 1024 ** 3) return `${(bytes / 1024 ** 2).toFixed(1)} MB`;
  return `${(bytes / 1024 ** 3).toFixed(2)} GB`;
}

export default function EngineeringPage() {
  const { data: objects, loading, error } = useCacheObjects({ priority: "P1", status: "ACTIVE" });

  return (
    <div>
      <span className="inline-block text-xs font-medium px-2 py-0.5 rounded-full bg-orange-100 text-orange-800 mb-4">P1 — Strongly Retained</span>
      <h1 className="text-2xl font-bold text-gray-900 mb-2">Engineering Manuals</h1>
      <p className="text-gray-500 text-sm mb-6">Rover maintenance, habitat systems, and mission runbooks. Kept locally even under bandwidth pressure.</p>
      {error && <p className="text-red-500 text-sm mb-4">{error}</p>}
      {loading ? <LoadingSpinner /> : (
        <div className="space-y-3">
          {objects?.length ? objects.map(obj => (
            <div key={obj.id} className="bg-white rounded-lg shadow p-5 flex items-start justify-between">
              <div className="flex-1">
                <div className="flex items-center gap-2 mb-1">
                  <h2 className="text-sm font-semibold text-gray-900">{obj.name}</h2>
                  {obj.pinned && <span className="text-xs bg-orange-100 text-orange-700 px-1.5 py-0.5 rounded font-medium">Pinned</span>}
                </div>
                <p className="text-xs text-gray-400">{obj.content_type ?? "document"} · {formatBytes(obj.size_bytes)}</p>
                {obj.tags.length > 0 && (
                  <div className="flex gap-1 mt-2 flex-wrap">
                    {obj.tags.map(t => <span key={t} className="text-xs bg-gray-100 text-gray-500 px-1.5 py-0.5 rounded">{t}</span>)}
                  </div>
                )}
              </div>
              {obj.source_url && (
                <a href={obj.source_url} target="_blank" rel="noreferrer"
                  className="ml-4 text-xs text-indigo-600 hover:underline shrink-0">Open ↗</a>
              )}
            </div>
          )) : (
            <div className="bg-white rounded-lg shadow p-8 text-center text-gray-400 text-sm">No P1 engineering content cached.</div>
          )}
        </div>
      )}
    </div>
  );
}
