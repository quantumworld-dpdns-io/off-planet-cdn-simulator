"use client";
import { useState } from "react";
import { useSites, useBandwidthWindows } from "@/lib/hooks";
import { LoadingSpinner } from "@/components/ui/loading-spinner";

const DEV_ORG_ID = process.env.NEXT_PUBLIC_DEV_ORG_ID ?? "00000000-0000-0000-0000-000000000001";
const BASE_URL = process.env.NEXT_PUBLIC_CONTROL_API_URL ?? "http://localhost:8080";

async function createWindow(body: object) {
  const res = await fetch(`${BASE_URL}/v1/bandwidth-windows`, {
    method: "POST",
    headers: { "Content-Type": "application/json", "X-Org-ID": DEV_ORG_ID },
    body: JSON.stringify(body),
  });
  if (!res.ok) throw new Error(`API error ${res.status}`);
  return res.json();
}

function fmtBps(bps: number): string {
  if (bps < 1000) return `${bps} bps`;
  if (bps < 1_000_000) return `${(bps / 1000).toFixed(1)} Kbps`;
  if (bps < 1_000_000_000) return `${(bps / 1_000_000).toFixed(1)} Mbps`;
  return `${(bps / 1_000_000_000).toFixed(2)} Gbps`;
}

export default function BandwidthWindowsPage() {
  const { data: sites } = useSites();
  const { data: windows, loading, error, refetch } = useBandwidthWindows();

  const [showForm, setShowForm] = useState(false);
  const [saving, setSaving] = useState(false);
  const [formError, setFormError] = useState<string | null>(null);

  const [siteId, setSiteId] = useState("");
  const [label, setLabel] = useState("");
  const [windowStart, setWindowStart] = useState("");
  const [windowEnd, setWindowEnd] = useState("");
  const [bandwidthMbps, setBandwidthMbps] = useState("10");
  const [reliability, setReliability] = useState("0.95");

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault();
    setFormError(null);
    if (!windowStart || !windowEnd || !bandwidthMbps) return;
    setSaving(true);
    try {
      await createWindow({
        site_id: siteId || undefined,
        label: label || undefined,
        window_start: new Date(windowStart).toISOString(),
        window_end: new Date(windowEnd).toISOString(),
        bandwidth_bps: Math.round(parseFloat(bandwidthMbps) * 1_000_000),
        reliability_score: parseFloat(reliability) || 1.0,
      });
      setShowForm(false);
      setLabel(""); setWindowStart(""); setWindowEnd(""); setBandwidthMbps("10"); setReliability("0.95");
      refetch();
    } catch (e: unknown) {
      setFormError(e instanceof Error ? e.message : "Failed to create window");
    } finally {
      setSaving(false);
    }
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Bandwidth Windows</h1>
        <button onClick={() => { setShowForm(v => !v); setFormError(null); }}
          className="bg-indigo-600 text-white px-4 py-2 rounded-md text-sm hover:bg-indigo-700">
          {showForm ? "Cancel" : "Add Window"}
        </button>
      </div>

      {showForm && (
        <form onSubmit={handleCreate} className="bg-white rounded-lg shadow p-6 mb-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
            <div>
              <label className="block text-xs font-medium text-gray-600 mb-1">Site (optional)</label>
              <select value={siteId} onChange={e => setSiteId(e.target.value)}
                className="w-full border border-gray-300 rounded px-3 py-2 text-sm">
                <option value="">All sites</option>
                {sites?.map(s => <option key={s.id} value={s.id}>{s.name}</option>)}
              </select>
            </div>
            <div>
              <label className="block text-xs font-medium text-gray-600 mb-1">Label</label>
              <input value={label} onChange={e => setLabel(e.target.value)}
                className="w-full border border-gray-300 rounded px-3 py-2 text-sm" placeholder="Lunar downlink pass" />
            </div>
            <div>
              <label className="block text-xs font-medium text-gray-600 mb-1">Window Start *</label>
              <input type="datetime-local" required value={windowStart} onChange={e => setWindowStart(e.target.value)}
                className="w-full border border-gray-300 rounded px-3 py-2 text-sm" />
            </div>
            <div>
              <label className="block text-xs font-medium text-gray-600 mb-1">Window End *</label>
              <input type="datetime-local" required value={windowEnd} onChange={e => setWindowEnd(e.target.value)}
                className="w-full border border-gray-300 rounded px-3 py-2 text-sm" />
            </div>
            <div>
              <label className="block text-xs font-medium text-gray-600 mb-1">Bandwidth (Mbps) *</label>
              <input type="number" required min="0.001" step="any" value={bandwidthMbps} onChange={e => setBandwidthMbps(e.target.value)}
                className="w-full border border-gray-300 rounded px-3 py-2 text-sm" placeholder="10" />
            </div>
            <div>
              <label className="block text-xs font-medium text-gray-600 mb-1">Reliability (0–1)</label>
              <input type="number" min="0" max="1" step="0.01" value={reliability} onChange={e => setReliability(e.target.value)}
                className="w-full border border-gray-300 rounded px-3 py-2 text-sm" />
            </div>
          </div>
          {formError && <p className="text-red-500 text-sm mb-3">{formError}</p>}
          <button type="submit" disabled={saving}
            className="bg-indigo-600 text-white px-4 py-2 rounded text-sm hover:bg-indigo-700 disabled:opacity-50">
            {saving ? "Saving…" : "Create Window"}
          </button>
        </form>
      )}

      {error && <p className="text-red-500 text-sm mb-4">{error}</p>}
      <div className="bg-white rounded-lg shadow overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>{["Label","Start","End","Bandwidth","Reliability"].map(h => (
              <th key={h} className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">{h}</th>
            ))}</tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {loading ? (
              <tr><td colSpan={5}><LoadingSpinner /></td></tr>
            ) : windows?.length ? windows.map(w => (
              <tr key={w.id} className="hover:bg-gray-50">
                <td className="px-6 py-4 text-sm font-medium text-gray-900">{w.label ?? "—"}</td>
                <td className="px-6 py-4 text-sm text-gray-500">{new Date(w.window_start).toLocaleString()}</td>
                <td className="px-6 py-4 text-sm text-gray-500">{new Date(w.window_end).toLocaleString()}</td>
                <td className="px-6 py-4 text-sm text-gray-700">{fmtBps(w.bandwidth_bps)}</td>
                <td className="px-6 py-4 text-sm text-gray-700">{(w.reliability_score * 100).toFixed(0)}%</td>
              </tr>
            )) : (
              <tr><td colSpan={5} className="px-6 py-8 text-center text-gray-400 text-sm">No bandwidth windows configured.</td></tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
