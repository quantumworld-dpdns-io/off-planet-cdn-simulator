export default function ContentCatalogPage() {
  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Content Catalog</h1>
      <div className="flex gap-3 mb-4">
        <select className="border border-gray-300 rounded-md px-3 py-2 text-sm"><option>All Priority Classes</option></select>
        <input placeholder="Filter by tag…" className="border border-gray-300 rounded-md px-3 py-2 text-sm" />
      </div>
      <p className="text-gray-400 text-sm">Phase 1</p>
    </div>
  );
}
