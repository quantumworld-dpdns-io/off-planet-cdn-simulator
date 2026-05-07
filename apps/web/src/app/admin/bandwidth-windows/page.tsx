"use client";
import { useBandwidthWindows } from "@/lib/hooks";
import { LoadingSpinner } from "@/components/ui/loading-spinner";

function fmtBps(bps: number): string {
  if (bps < 1000) return `${bps} bps`;
  if (bps < 1_000_000) return `${(bps / 1000).toFixed(1)} Kbps`;
  if (bps < 1_000_000_000) return `${(bps / 1_000_000).toFixed(1)} Mbps`;
  return `${(bps / 1_000_000_000).toFixed(2)} Gbps`;
}

export default function BandwidthWindowsPage() {
  const { data: windows, loading, error } = useBandwidthWindows();

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Bandwidth Windows</h1>
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
