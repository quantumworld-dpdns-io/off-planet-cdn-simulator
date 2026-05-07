"use client";
import { useState, useCallback } from "react";
import { useMirrorSources, useMirrorArtifacts } from "@/lib/hooks";
import { api } from "@/lib/api-client";
import { LoadingSpinner } from "@/components/ui/loading-spinner";

const PACKAGE_REGISTRY_TYPES = ["npm", "pypi", "crates_io", "oci"] as const;

export default function PackageMirrorsPage() {
  const { data: allSources, loading: sourcesLoading, error: sourcesError, refetch: refetchSources } = useMirrorSources();
  const { data: artifacts, loading: artifactsLoading, error: artifactsError } = useMirrorArtifacts();

  const sources = allSources?.filter(s => PACKAGE_REGISTRY_TYPES.includes(s.registry_type as typeof PACKAGE_REGISTRY_TYPES[number])) ?? [];

  const [showForm, setShowForm] = useState(false);
  const [formRegistryType, setFormRegistryType] = useState<string>("npm");
  const [formUpstreamURL, setFormUpstreamURL] = useState("");
  const [formLabel, setFormLabel] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [formError, setFormError] = useState<string | null>(null);

  const sourceMap = Object.fromEntries((allSources ?? []).map(s => [s.id, s]));

  const handleSubmit = useCallback(async (e: React.FormEvent) => {
    e.preventDefault();
    setFormError(null);
    if (!formUpstreamURL.trim()) {
      setFormError("Upstream URL is required.");
      return;
    }
    setSubmitting(true);
    try {
      await api.mirrors.createSource({
        registry_type: formRegistryType,
        upstream_url: formUpstreamURL.trim(),
        label: formLabel.trim() || undefined,
      });
      setFormUpstreamURL("");
      setFormLabel("");
      setShowForm(false);
      refetchSources();
    } catch (err) {
      setFormError(err instanceof Error ? err.message : String(err));
    } finally {
      setSubmitting(false);
    }
  }, [formRegistryType, formUpstreamURL, formLabel, refetchSources]);

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Package Mirrors</h1>
        <button
          onClick={() => setShowForm(v => !v)}
          className="px-4 py-2 bg-indigo-600 text-white text-sm font-medium rounded-md hover:bg-indigo-700 transition-colors"
        >
          {showForm ? "Cancel" : "Add Source"}
        </button>
      </div>

      {showForm && (
        <div className="bg-white rounded-lg shadow p-6 mb-6">
          <h2 className="text-lg font-semibold text-gray-800 mb-4">Add Mirror Source</h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Registry Type</label>
              <select
                value={formRegistryType}
                onChange={e => setFormRegistryType(e.target.value)}
                className="w-full border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
              >
                {PACKAGE_REGISTRY_TYPES.map(t => (
                  <option key={t} value={t}>{t}</option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Upstream URL</label>
              <input
                type="url"
                value={formUpstreamURL}
                onChange={e => setFormUpstreamURL(e.target.value)}
                placeholder="https://registry.npmjs.org"
                required
                className="w-full border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Label (optional)</label>
              <input
                type="text"
                value={formLabel}
                onChange={e => setFormLabel(e.target.value)}
                placeholder="e.g. primary npm mirror"
                className="w-full border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
              />
            </div>
            {formError && <p className="text-red-500 text-sm">{formError}</p>}
            <button
              type="submit"
              disabled={submitting}
              className="px-4 py-2 bg-indigo-600 text-white text-sm font-medium rounded-md hover:bg-indigo-700 disabled:opacity-50 transition-colors"
            >
              {submitting ? "Creating..." : "Create Source"}
            </button>
          </form>
        </div>
      )}

      {sourcesError && <p className="text-red-500 text-sm mb-4">{sourcesError}</p>}

      <div className="bg-white rounded-lg shadow overflow-hidden mb-8">
        <div className="px-6 py-4 border-b border-gray-200">
          <h2 className="text-base font-semibold text-gray-800">Mirror Sources</h2>
          <p className="text-xs text-gray-500 mt-0.5">npm · PyPI · crates.io · OCI</p>
        </div>
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              {["Registry", "Label", "Upstream URL", "Enabled", "Created At"].map(h => (
                <th key={h} className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">{h}</th>
              ))}
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {sourcesLoading ? (
              <tr><td colSpan={5}><LoadingSpinner /></td></tr>
            ) : sources.length ? sources.map(src => (
              <tr key={src.id} className="hover:bg-gray-50">
                <td className="px-6 py-4 text-sm font-mono text-gray-700">{src.registry_type}</td>
                <td className="px-6 py-4 text-sm text-gray-500">{src.label || <span className="italic text-gray-300">—</span>}</td>
                <td className="px-6 py-4 text-sm text-gray-500 max-w-xs truncate">{src.upstream_url}</td>
                <td className="px-6 py-4 text-sm">
                  {src.enabled ? (
                    <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-green-100 text-green-800">enabled</span>
                  ) : (
                    <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-600">disabled</span>
                  )}
                </td>
                <td className="px-6 py-4 text-sm text-gray-500">{new Date(src.created_at).toLocaleString()}</td>
              </tr>
            )) : (
              <tr>
                <td colSpan={5} className="px-6 py-8 text-center text-gray-400 text-sm">No package mirror sources yet.</td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      {artifactsError && <p className="text-red-500 text-sm mb-4">{artifactsError}</p>}

      <div className="bg-white rounded-lg shadow overflow-hidden">
        <div className="px-6 py-4 border-b border-gray-200">
          <h2 className="text-base font-semibold text-gray-800">Recent Artifacts</h2>
        </div>
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              {["Name", "Version", "Registry", "Size", "Synced At"].map(h => (
                <th key={h} className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">{h}</th>
              ))}
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {artifactsLoading ? (
              <tr><td colSpan={5}><LoadingSpinner /></td></tr>
            ) : artifacts?.length ? artifacts.map(a => (
              <tr key={a.id} className="hover:bg-gray-50">
                <td className="px-6 py-4 text-sm font-mono text-gray-700">{a.name}</td>
                <td className="px-6 py-4 text-sm text-gray-500">{a.version}</td>
                <td className="px-6 py-4 text-sm text-gray-500">{sourceMap[a.source_id]?.registry_type ?? a.source_id.slice(0, 8) + "…"}</td>
                <td className="px-6 py-4 text-sm text-gray-500">
                  {a.size_bytes ? `${(a.size_bytes / 1024 / 1024).toFixed(2)} MB` : "—"}
                </td>
                <td className="px-6 py-4 text-sm text-gray-500">
                  {a.synced_at ? new Date(a.synced_at).toLocaleString() : "—"}
                </td>
              </tr>
            )) : (
              <tr>
                <td colSpan={5} className="px-6 py-8 text-center text-gray-400 text-sm">No artifacts mirrored yet.</td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
