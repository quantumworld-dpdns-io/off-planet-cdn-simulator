import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Off-Planet CDN",
  description: "Priority-aware CDN for Moon/Mars habitats",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
