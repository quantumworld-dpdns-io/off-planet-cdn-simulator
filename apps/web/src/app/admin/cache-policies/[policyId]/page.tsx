export default function PolicyDetailPage({ params }: { params: { policyId: string } }) {
  return <div><h1 className="text-2xl font-bold text-gray-900 mb-2">Policy Detail</h1><p className="text-gray-500 text-sm">ID: {params.policyId}</p></div>;
}
