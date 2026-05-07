import type { PriorityClassLevel } from "@/types";

const colors: Record<PriorityClassLevel, string> = {
  P0: "bg-red-100 text-red-800",
  P1: "bg-orange-100 text-orange-800",
  P2: "bg-yellow-100 text-yellow-800",
  P3: "bg-blue-100 text-blue-800",
  P4: "bg-gray-100 text-gray-600",
  P5: "bg-gray-50 text-gray-400",
};

const labels: Record<PriorityClassLevel, string> = {
  P0: "P0 Medical",
  P1: "P1 Engineering",
  P2: "P2 Education",
  P3: "P3 Packages",
  P4: "P4 Entertainment",
  P5: "P5 Stale",
};

export function PriorityBadge({ priority }: { priority: PriorityClassLevel }) {
  return (
    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${colors[priority] ?? colors.P5}`}>
      {labels[priority] ?? priority}
    </span>
  );
}
