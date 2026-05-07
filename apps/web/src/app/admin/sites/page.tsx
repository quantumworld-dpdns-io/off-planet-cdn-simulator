"use client";
import { useState } from "react";
import Link from "next/link";
import { useSites } from "@/lib/hooks";
import { api } from "@/lib/api-client";
import { LoadingSpinner } from "@/components/ui/loading-spinner";

export default function SitesPage() {
  const { data: sites, loading, error, refetch } = useSites();
  const [showForm, setShowForm] = useState(false);
  const [name, setName] = useState("");
  const [location, setLocation] = useState("");
  const [saving, setSaving] = useState(false);

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault();
    if (!name.trim()) return;
    setSaving(true);
    try {
      await api.sites.create({ name: name.trim(), location: location.trim() || undefined });
      setName(""); setLocation(""); setShowForm(false);
      refetch();
    } finally {
      setSaving(false);
    }
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Sites</h1>
        <button onClick={() => setShowForm(v => !v)}
          className="bg-indigo-600 text-white px-4 py-2 rounded-md text-sm hover:bg-indigo-700">
          {showForm ? "Cancel" : "Add Site"}
        </button>
      </div>

      {showForm && (
        <form onSubmit={handleCreate} className="bg-white rounded-lg shadow p-6 mb-6 flex gap-4 items-end">
          <div className="flex-1">
            <label className="block text-xs font-medium text-gray-600 mb-1">Name *</label>
            <input value={name} onChange={e => setName(e.target.value)} required
              className="w-full border border-gray-300 rounded px-3 py-2 text-sm" placeholder="Lunar Base Alpha" />
          </div>
          <div className="flex-1">
            <label className="block text-xs font-medium text-gray-600 mb-1">Location</label>
            <input value={location} onChange={e => setLocation(e.target.value)}
              className="w-full border border-gray-300 rounded px-3 py-2 text-sm" placeholder="Sea of Tranquility" />
          </div>
          <button type="submit" disabled={saving}
            className="bg-indigo-600 text-white px-4 py-2 rounded text-sm hover:bg-indigo-700 disabled:opacity-50">
            {saving ? "Saving…" : "Create"}
          </button>
        </form>
      )}

      {error && <p className="text-red-500 text-sm mb-4">{error}</p>}

      <div className="bg-white rounded-lg shadow overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              {["Name", "Location", "Created"].map(h => (
                <th key={h} className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">{h}</th>
              ))}
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {loading ? (
              <tr><td colSpan={3}><LoadingSpinner /></td></tr>
            ) : sites?.length ? (
              sites.map(site => (
                <tr key={site.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 text-sm font-medium">
                    <Link href={`/admin/sites/${site.id}`} className="text-indigo-600 hover:underline">{site.name}</Link>
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-500">{site.location ?? "—"}</td>
                  <td className="px-6 py-4 text-sm text-gray-500">{new Date(site.created_at).toLocaleDateString()}</td>
                </tr>
              ))
            ) : (
              <tr><td colSpan={3} className="px-6 py-8 text-center text-gray-400 text-sm">No sites yet. Add your first site above.</td></tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
