export default function JobDetailPage({ params }: { params: { jobId: string } }) {
  return <div><h1 className="text-2xl font-bold text-gray-900 mb-2">Job Detail</h1><p className="text-gray-500 text-sm">ID: {params.jobId}</p></div>;
}
