export default function NodeDetailPage({ params }: { params: { nodeId: string } }) {
  return <div><h1 className="text-2xl font-bold text-gray-900 mb-2">Node Detail</h1><p className="text-gray-500 text-sm">ID: {params.nodeId}</p><p className="text-gray-400 text-sm mt-4">Phase 1</p></div>;
}
