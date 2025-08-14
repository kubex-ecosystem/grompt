'use client';

import { Suspense } from 'react';
import AgentForm from '../../../components/AgentForm';

function AgentFormWithSuspense() {
  return (
    <Suspense fallback={<div className="flex items-center justify-center min-h-screen bg-gray-900 text-white">
      <div className="text-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-400 mx-auto mb-4"></div>
        <p>Loading...</p>
      </div>
    </div>}>
      <AgentForm />
    </Suspense>
  );
}

export default function NewAgentPage() {
  return <AgentFormWithSuspense />;
}
