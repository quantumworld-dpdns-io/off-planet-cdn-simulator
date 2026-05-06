export default function DashboardPage() {
  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Admin Dashboard</h1>
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
        {[
          { label: "Nodes Online", value: "—" },
          { label: "Cache Fill %", value: "—" },
          { label: "Active Jobs", value: "—" },
          { label: "Recent Alerts", value: "—" },
        ].map(card => (
          <div key={card.label} className="bg-white rounded-lg shadow p-6">
            <p className="text-sm text-gray-500">{card.label}</p>
            <p className="text-3xl font-bold text-gray-900 mt-1">{card.value}</p>
          </div>
        ))}
      </div>
      <p className="text-gray-500 text-sm">System overview — live data available in Phase 1.</p>
    </div>
  );
}
