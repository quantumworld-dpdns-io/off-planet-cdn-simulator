export default function OfflineStatusPage() {
  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Offline Status</h1>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {[
          { label: "Cache Fill %", value: "—" },
          { label: "Last Sync", value: "—" },
          { label: "Next Contact Window", value: "—" },
        ].map(card => (
          <div key={card.label} className="bg-white rounded-lg shadow p-6">
            <p className="text-sm text-gray-500">{card.label}</p>
            <p className="text-2xl font-bold text-gray-900 mt-1">{card.value}</p>
          </div>
        ))}
      </div>
    </div>
  );
}
