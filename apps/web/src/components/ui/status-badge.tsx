import type { NodeStatus } from "@/types";

const colors: Record<NodeStatus, string> = {
  ONLINE: "bg-green-100 text-green-800",
  OFFLINE: "bg-red-100 text-red-800",
  DEGRADED: "bg-yellow-100 text-yellow-800",
  UNKNOWN: "bg-gray-100 text-gray-600",
};

export function StatusBadge({ status }: { status: NodeStatus }) {
  return (
    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${colors[status] ?? colors.UNKNOWN}`}>
      {status}
    </span>
  );
}
