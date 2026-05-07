"use client";
import { useState, useMemo } from "react";
import { useSites, useCacheObjects } from "@/lib/hooks";
import { LoadingSpinner } from "@/components/ui/loading-spinner";

function formatBytes(bytes: number): string {
  if (bytes < 1024 ** 2) return `${(bytes / 1024).toFixed(1)} KB`;
  if (bytes < 1024 ** 3) return `${(bytes / 1024 ** 2).toFixed(1)} MB`;
  return `${(bytes / 1024 ** 3).toFixed(2)} GB`;
}

interface EvictionCandidate {
  id: string;
  name: string;
  size_bytes: number;
  cumulative_freed: number;
  content_type?: string;
}

interface SimResult {
  can_meet_target: boolean;
  freed_bytes: number;
  candidates: EvictionCandidate[];
}

function runSimulation(
  objects: { id: string; name: string; size_bytes: number; pinned: boolean; content_type?: string }[],
  targetBytes: number
): SimResult {
  // Exclude pinned objects
  const evictable = objects.filter(o => !o.pinned);
  // Largest first
  evictable.sort((a, b) => b.size_bytes - a.size_bytes);

  const candidates: EvictionCandidate[] = [];
  let cumulative = 0;

  for (const obj of evictable) {
    if (cumulative >= targetBytes) break;
    cumulative += obj.size_bytes;
    candidates.push({ ...obj, cumulative_freed: cumulative });
  }

  return {
    can_meet_target: cumulative >= targetBytes,
    freed_bytes: cumulative,
    candidates,
  };
}

export default function EvictionSimulatorPage() {
  const { data: sites, loading: sitesLoading } = useSites();
  const [selectedSiteId, setSelectedSiteId] = useState("");
  const [targetMB, setTargetMB] = useState("1024");
  const [result, setResult] = useState<SimResult | null>(null);
  const [hasRun, setHasRun] = useState(false);

  const { data: objects, loading: objectsLoading } = useCacheObjects(
    selectedSiteId ? { site_id: selectedSiteId, status: "ACTIVE" } : undefined
  );

  const loading = sitesLoading || objectsLoading;

  function handleRun() {
    if (!objects) return;
    const targetBytes = parseFloat(targetMB) * 1024 * 1024;
    if (isNaN(targetBytes) || targetBytes <= 0) return;
    setResult(runSimulation(objects, targetBytes));
    setHasRun(true);
  }

  const totalEvictable = useMemo(
    () => objects?.filter(o => !o.pinned).reduce((s, o) => s + o.size_bytes, 0) ?? 0,
    [objects]
  );

  const pinnedCount = useMemo(
    () => objects?.filter(o => o.pinned).length ?? 0,
    [objects]
  );

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-2">Eviction Simulator</h1>
      <p className="text-sm text-gray-500 mb-6">
        Simulate which objects would be evicted to free a target amount of space. Pinned objects are never evicted.
      </p>

      {/* Setup form */}
      <div className="bg-white rounded-lg shadow p-6 max-w-lg mb-6">
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Site</label>
            {sitesLoading ? <LoadingSpinner label="Loading sites…" /> : (
              <select
                value={selectedSiteId}
                onChange={e => { setSelectedSiteId(e.target.value); setResult(null); setHasRun(false); }}
                className="w-full border border-gray-300 rounded-md px-3 py-2 text-sm">
                <option value="">Select site…</option>
                {sites?.map(s => (
                  <option key={s.id} value={s.id}>{s.name}{s.location ? ` — ${s.location}` : ""}</option>
                ))}
              </select>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Target free space (MB)</label>
            <input
              type="number"
              value={targetMB}
              onChange={e => { setTargetMB(e.target.value); setResult(null); setHasRun(false); }}
              min="1"
              className="w-full border border-gray-300 rounded-md px-3 py-2 text-sm"
              placeholder="1024" />
            <p className="text-xs text-gray-400 mt-1">= {formatBytes(parseFloat(targetMB || "0") * 1024 * 1024)}</p>
          </div>

          {selectedSiteId && !objectsLoading && objects && (
            <div className="text-xs text-gray-500 bg-gray-50 rounded p-3 space-y-1">
              <p><span className="font-medium">{objects.length}</span> total objects for site</p>
              <p><span className="font-medium">{pinnedCount}</span> pinned (excluded from simulation)</p>
              <p><span className="font-medium">{formatBytes(totalEvictable)}</span> maximum evictable</p>
            </div>
          )}

          <button
            onClick={handleRun}
            disabled={!selectedSiteId || loading || !objects?.length}
            className="w-full bg-indigo-600 text-white px-4 py-2 rounded-md text-sm hover:bg-indigo-700 disabled:opacity-40 disabled:cursor-not-allowed">
            {loading ? "Loading…" : "Run Simulation"}
          </button>
        </div>
      </div>

      {/* Results */}
      {hasRun && result && (
        <div>
          {/* Summary banner */}
          <div className={`rounded-lg p-4 mb-4 border ${result.can_meet_target ? "bg-green-50 border-green-200" : "bg-red-50 border-red-200"}`}>
            <p className={`text-sm font-semibold ${result.can_meet_target ? "text-green-800" : "text-red-800"}`}>
              {result.can_meet_target
                ? `✓ Target met — ${formatBytes(result.freed_bytes)} can be freed by evicting ${result.candidates.length} object${result.candidates.length !== 1 ? "s" : ""}`
                : `✗ Cannot meet target — only ${formatBytes(result.freed_bytes)} available from ${result.candidates.length} evictable object${result.candidates.length !== 1 ? "s" : ""}`}
            </p>
            <p className="text-xs text-gray-500 mt-1">
              Strategy: largest non-pinned objects evicted first (greedy). This is a dry run — no data was deleted.
            </p>
          </div>

          {/* Candidates table */}
          {result.candidates.length > 0 && (
            <div className="bg-white rounded-lg shadow overflow-hidden">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    {["#", "Object Name", "Type", "Size", "Cumulative Freed"].map(h => (
                      <th key={h} className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">{h}</th>
                    ))}
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-100">
                  {result.candidates.map((c, i) => (
                    <tr key={c.id} className={i % 2 === 0 ? "bg-white" : "bg-gray-50"}>
                      <td className="px-6 py-4 text-sm text-gray-500">{i + 1}</td>
                      <td className="px-6 py-4 text-sm font-medium text-gray-900">{c.name}</td>
                      <td className="px-6 py-4 text-sm text-gray-500">{c.content_type ?? "—"}</td>
                      <td className="px-6 py-4 text-sm text-gray-700">{formatBytes(c.size_bytes)}</td>
                      <td className="px-6 py-4 text-sm">
                        <div className="flex items-center gap-2">
                          <div className="flex-1 bg-gray-200 rounded-full h-1.5">
                            <div
                              className="bg-indigo-500 h-1.5 rounded-full"
                              style={{ width: `${Math.min(100, (c.cumulative_freed / (parseFloat(targetMB) * 1024 * 1024)) * 100)}%` }}
                            />
                          </div>
                          <span className="text-xs text-gray-500 whitespace-nowrap">{formatBytes(c.cumulative_freed)}</span>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}

          {result.candidates.length === 0 && (
            <div className="bg-white rounded-lg shadow p-8 text-center text-gray-400 text-sm">
              No evictable objects found for this site.
            </div>
          )}
        </div>
      )}
    </div>
  );
}
