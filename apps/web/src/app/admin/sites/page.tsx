export default function SitesPage() {
  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Sites</h1>
        <button className="bg-brand-600 text-white px-4 py-2 rounded-md text-sm hover:bg-brand-700">Add Site</button>
      </div>
      <div className="bg-white rounded-lg shadow">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>{["Name","Location","Nodes","Created"].map(h => (
              <th key={h} className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">{h}</th>
            ))}</tr>
          </thead>
          <tbody><tr><td colSpan={4} className="px-6 py-8 text-center text-gray-400 text-sm">Loading in Phase 1</td></tr></tbody>
        </table>
      </div>
    </div>
  );
}
