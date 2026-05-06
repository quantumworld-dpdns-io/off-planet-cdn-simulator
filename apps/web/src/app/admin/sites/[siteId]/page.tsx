export default function SiteDetailPage({ params }: { params: { siteId: string } }) {
  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-2">Site Detail</h1>
      <p className="text-gray-500 text-sm mb-6">ID: {params.siteId}</p>
      <div className="grid grid-cols-2 gap-6">
        <div className="bg-white rounded-lg shadow p-6"><h2 className="font-semibold mb-4">Nodes</h2><p className="text-gray-400 text-sm">Phase 1</p></div>
        <div className="bg-white rounded-lg shadow p-6"><h2 className="font-semibold mb-4">Cache Summary</h2><p className="text-gray-400 text-sm">Phase 1</p></div>
      </div>
    </div>
  );
}
