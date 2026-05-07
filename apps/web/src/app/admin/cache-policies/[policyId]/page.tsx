"use client";
import { use, useState, useEffect } from "react";
import Link from "next/link";
import { api } from "@/lib/api-client";
import type { CachePolicy } from "@/types";
import { LoadingSpinner } from "@/components/ui/loading-spinner";

export default function PolicyDetailPage({ params }: { params: Promise<{ policyId: string }> }) {
  const { policyId } = use(params);
  const [policy, setPolicy] = useState<CachePolicy | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    api.policies.get(policyId)
      .then(setPolicy)
      .catch((e: unknown) => setError(e instanceof Error ? e.message : String(e)))
      .finally(() => setLoading(false));
  }, [policyId]);

  if (loading) return <LoadingSpinner />;
  if (error) return <p className="text-red-500 text-sm">{error}</p>;
  if (!policy) return <p className="text-gray-500 text-sm">Policy not found.</p>;

  return (
    <div>
      <Link href="/admin/cache-policies" className="text-indigo-600 text-sm hover:underline">← Cache Policies</Link>
      <h1 className="text-2xl font-bold text-gray-900 mt-2 mb-6">{policy.name}</h1>
      <div className="bg-white rounded-lg shadow overflow-hidden">
        <dl className="divide-y divide-gray-100">
          {([
            ["ID", <span key="id" className="font-mono text-xs">{policy.id}</span>],
            ["Site ID", <span key="si" className="font-mono text-xs">{policy.site_id}</span>],
            ["Description", policy.description ?? "—"],
            ["Enabled", policy.enabled ? "Yes" : "No"],
            ["Created", new Date(policy.created_at).toLocaleString()],
          ] as [string, React.ReactNode][]).map(([label, value]) => (
            <div key={label} className="px-6 py-4 flex items-center">
              <dt className="text-sm font-medium text-gray-500 w-48">{label}</dt>
              <dd className="text-sm text-gray-900 flex-1">{value}</dd>
            </div>
          ))}
        </dl>
      </div>
    </div>
  );
}
