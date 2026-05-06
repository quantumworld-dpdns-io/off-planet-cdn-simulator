export default function ObjectDetailPage({ params }: { params: { objectId: string } }) {
  return <div><h1 className="text-2xl font-bold text-gray-900 mb-2">Object Detail</h1><p className="text-gray-500 text-sm">ID: {params.objectId}</p></div>;
}
