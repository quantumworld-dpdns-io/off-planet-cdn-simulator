import type { JobStatus } from "@/types";

const colors: Record<JobStatus, string> = {
  PENDING: "bg-gray-100 text-gray-600",
  RUNNING: "bg-blue-100 text-blue-800",
  DONE: "bg-green-100 text-green-800",
  FAILED: "bg-red-100 text-red-800",
  CANCELLED: "bg-yellow-100 text-yellow-700",
};

export function JobStatusBadge({ status }: { status: JobStatus }) {
  return (
    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${colors[status] ?? colors.PENDING}`}>
      {status}
    </span>
  );
}
