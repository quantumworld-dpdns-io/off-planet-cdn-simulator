"use client";
import { clsx } from "clsx";

type BadgeVariant = "p0"|"p1"|"p2"|"p3"|"p4"|"p5"|"online"|"offline"|"running"|"done"|"failed";

const variantClasses: Record<BadgeVariant, string> = {
  p0: "bg-red-100 text-red-800",
  p1: "bg-orange-100 text-orange-800",
  p2: "bg-yellow-100 text-yellow-800",
  p3: "bg-blue-100 text-blue-800",
  p4: "bg-gray-100 text-gray-700",
  p5: "bg-gray-50 text-gray-400",
  online: "bg-green-100 text-green-800",
  offline: "bg-red-100 text-red-800",
  running: "bg-blue-100 text-blue-800",
  done: "bg-green-100 text-green-800",
  failed: "bg-red-100 text-red-800",
};

export function Badge({ variant, children }: { variant: BadgeVariant; children: React.ReactNode }) {
  return (
    <span className={clsx("inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium", variantClasses[variant])}>
      {children}
    </span>
  );
}
