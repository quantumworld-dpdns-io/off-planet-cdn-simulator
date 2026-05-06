export default function PreloadJobsPage() {
  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Preload Jobs</h1>
      <div className="flex gap-2 mb-4 text-xs">
        {["PENDING","RUNNING","DONE","FAILED","CANCELLED"].map(s => (
          <span key={s} className="px-2 py-1 bg-gray-100 rounded">{s}</span>
        ))}
      </div>
      <p className="text-gray-400 text-sm">Phase 1</p>
    </div>
  );
}
