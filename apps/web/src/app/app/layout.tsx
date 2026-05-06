import Link from "next/link";

const navItems = [
  { href: "/app/home", label: "Home" },
  { href: "/app/search", label: "Search" },
  { href: "/app/manuals", label: "Manuals" },
  { href: "/app/medical", label: "Medical" },
  { href: "/app/engineering", label: "Engineering" },
  { href: "/app/education", label: "Education" },
  { href: "/app/entertainment", label: "Entertainment" },
  { href: "/app/packages", label: "Packages" },
  { href: "/app/models", label: "Models" },
  { href: "/app/downloads", label: "Downloads" },
  { href: "/app/offline-status", label: "Offline Status" },
];

export default function AppLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-gray-900 text-white">
        <div className="max-w-7xl mx-auto px-4 py-3 flex gap-6 overflow-x-auto">
          <span className="font-bold text-brand-300 shrink-0">Off-Planet CDN</span>
          {navItems.map(item => (
            <Link key={item.href} href={item.href} className="text-sm text-gray-300 hover:text-white whitespace-nowrap">{item.label}</Link>
          ))}
        </div>
      </nav>
      <main className="max-w-7xl mx-auto px-4 py-8">{children}</main>
    </div>
  );
}
