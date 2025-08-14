'use client';

import { useSearchParams } from 'next/navigation';
import { Suspense } from 'react';
import AgentsDashboard from '../../components/AgentsDashboard';
import AgentView from '../../components/AgentView';

function AgentsContent() {
  const searchParams = useSearchParams();
  const viewId = searchParams.get('view');

  if (viewId) {
    return <AgentView />;
  }

  return <AgentsDashboard />;
}

export default function AgentsPage() {
  return (
    <Suspense fallback={<div className="flex items-center justify-center min-h-screen bg-gray-900 text-white">
      <div className="text-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-400 mx-auto mb-4"></div>
        <p>Loading...</p>
      </div>
    </div>}>
      <AgentsContent />
    </Suspense>
  );
}
