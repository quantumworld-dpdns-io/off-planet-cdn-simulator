export default function NodesPage() {
  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Edge Nodes</h1>
      <div className="bg-white rounded-lg shadow">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>{["Name","Status","Cache Fill","Last Heartbeat"].map(h => (
              <th key={h} className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">{h}</th>
            ))}</tr>
          </thead>
          <tbody><tr><td colSpan={4} className="px-6 py-8 text-center text-gray-400 text-sm">Loading in Phase 1</td></tr></tbody>
        </table>
      </div>
    </div>
  );
}
