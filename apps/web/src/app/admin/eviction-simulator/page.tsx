export default function EvictionSimulatorPage() {
  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Eviction Simulator</h1>
      <div className="bg-white rounded-lg shadow p-6 max-w-lg">
        <div className="space-y-4">
          <div><label className="block text-sm font-medium text-gray-700 mb-1">Site</label>
            <select className="w-full border border-gray-300 rounded-md px-3 py-2"><option>Select site…</option></select></div>
          <div><label className="block text-sm font-medium text-gray-700 mb-1">Target free bytes</label>
            <input type="number" placeholder="1073741824" className="w-full border border-gray-300 rounded-md px-3 py-2" /></div>
          <div className="flex items-center gap-2">
            <input type="checkbox" id="dryrun" defaultChecked />
            <label htmlFor="dryrun" className="text-sm text-gray-700">Dry run (simulation only)</label></div>
          <button className="bg-brand-600 text-white px-4 py-2 rounded-md text-sm hover:bg-brand-700">Run Simulation</button>
        </div>
      </div>
      <p className="text-gray-400 text-sm mt-4">Candidate list — Phase 2</p>
    </div>
  );
}
