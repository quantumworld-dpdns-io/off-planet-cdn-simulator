"use client";
import { useState } from "react";
import { useCacheObjects } from "@/lib/hooks";
import { LoadingSpinner } from "@/components/ui/loading-spinner";

function formatBytes(bytes: number): string {
  if (bytes < 1024 ** 2) return `${(bytes / 1024).toFixed(1)} KB`;
  if (bytes < 1024 ** 3) return `${(bytes / 1024 ** 2).toFixed(1)} MB`;
  return `${(bytes / 1024 ** 3).toFixed(2)} GB`;
}

export default function SearchPage() {
  const [query, setQuery] = useState("");
  const { data: objects, loading } = useCacheObjects({ status: "ACTIVE" });

  const results = query.trim()
    ? objects?.filter(o =>
        o.name.toLowerCase().includes(query.toLowerCase()) ||
        (o.content_type ?? "").toLowerCase().includes(query.toLowerCase()) ||
        o.tags.some(t => t.toLowerCase().includes(query.toLowerCase()))
      ) ?? []
    : [];

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Search</h1>
      <input
        value={query}
        onChange={e => setQuery(e.target.value)}
        placeholder="Search cached content by name, type, or tag…"
        className="w-full border border-gray-300 rounded-lg px-4 py-3 mb-6 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
        autoFocus
      />
      {loading ? <LoadingSpinner /> : (
        <>
          {query.trim() && (
            <p className="text-xs text-gray-400 mb-3">{results.length} result{results.length !== 1 ? "s" : ""} for "{query}"</p>
          )}
          <div className="bg-white rounded-lg shadow divide-y divide-gray-100">
            {results.length ? results.map(obj => (
              <div key={obj.id} className="px-6 py-4 flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-gray-900">{obj.name}</p>
                  <p className="text-xs text-gray-400 mt-0.5">
                    {obj.content_type ?? "unknown"} · {formatBytes(obj.size_bytes)}
                    {obj.pinned && <span className="ml-2 text-indigo-600 font-medium">Pinned</span>}
                  </p>
                  {obj.tags.length > 0 && (
                    <div className="flex gap-1 mt-1 flex-wrap">
                      {obj.tags.map(t => (
                        <span key={t} className="text-xs bg-gray-100 text-gray-500 px-1.5 py-0.5 rounded">{t}</span>
                      ))}
                    </div>
                  )}
                </div>
                {obj.source_url && (
                  <a href={obj.source_url} target="_blank" rel="noreferrer"
                    className="text-xs text-indigo-600 hover:underline ml-4 shrink-0">Open ↗</a>
                )}
              </div>
            )) : query.trim() ? (
              <p className="px-6 py-8 text-sm text-gray-400 text-center">No results for "{query}".</p>
            ) : (
              <p className="px-6 py-8 text-sm text-gray-400 text-center">Start typing to search.</p>
            )}
          </div>
        </>
      )}
    </div>
  );
}
