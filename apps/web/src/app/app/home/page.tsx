"use client";
export default function UserHomePage() {
  return (
    <div>
      <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4 mb-6">
        <p className="text-yellow-800 text-sm font-medium">Offline status: Checking connection…</p>
        <p className="text-yellow-600 text-xs mt-1">Last sync: —</p>
      </div>
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Welcome to Off-Planet CDN</h1>
      <input placeholder="Search cached content…" className="w-full border border-gray-300 rounded-lg px-4 py-3 mb-8 text-sm" />
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {[
          { label: "Medical References", badge: "P0", color: "red", href: "/app/medical" },
          { label: "Engineering Manuals", badge: "P1", color: "orange", href: "/app/engineering" },
          { label: "Manuals", badge: "P1", color: "orange", href: "/app/manuals" },
        ].map(card => (
          <a key={card.href} href={card.href}
            className="bg-white rounded-lg shadow p-6 hover:shadow-md transition-shadow">
            <span className={`inline-block text-xs font-medium px-2 py-0.5 rounded-full bg-${card.color}-100 text-${card.color}-800 mb-3`}>{card.badge}</span>
            <h2 className="font-semibold text-gray-900">{card.label}</h2>
          </a>
        ))}
      </div>
    </div>
  );
}
