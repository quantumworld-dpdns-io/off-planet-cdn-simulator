import Link from "next/link";

const navItems = [
  { href: "/admin/dashboard", label: "Dashboard" },
  { href: "/admin/sites", label: "Sites" },
  { href: "/admin/nodes", label: "Nodes" },
  { href: "/admin/cache-policies", label: "Cache Policies" },
  { href: "/admin/content-catalog", label: "Content Catalog" },
  { href: "/admin/preload-jobs", label: "Preload Jobs" },
  { href: "/admin/eviction-simulator", label: "Eviction Simulator" },
  { href: "/admin/bandwidth-windows", label: "Bandwidth Windows" },
  { href: "/admin/package-mirrors", label: "Package Mirrors" },
  { href: "/admin/model-mirrors", label: "Model Mirrors" },
  { href: "/admin/incidents", label: "Incidents" },
  { href: "/admin/audit-logs", label: "Audit Logs" },
  { href: "/admin/settings", label: "Settings" },
];

export default function AdminLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex h-screen bg-gray-100">
      <nav className="w-64 bg-gray-900 text-white flex flex-col">
        <div className="p-4 border-b border-gray-700">
          <h1 className="text-lg font-bold">Off-Planet CDN</h1>
          <p className="text-xs text-gray-400">Admin Console</p>
        </div>
        <ul className="flex-1 overflow-y-auto py-4">
          {navItems.map(item => (
            <li key={item.href}>
              <Link href={item.href}
                className="block px-4 py-2 text-sm text-gray-300 hover:bg-gray-700 hover:text-white transition-colors">
                {item.label}
              </Link>
            </li>
          ))}
        </ul>
      </nav>
      <main className="flex-1 overflow-y-auto p-8">{children}</main>
    </div>
  );
}
