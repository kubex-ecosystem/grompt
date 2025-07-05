import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';

const AgentsDashboard = () => {
  const [agents, setAgents] = useState([]);

  useEffect(() => {
    fetch('/api/agents')
      .then(res => res.json())
      .then(setAgents)
      .catch(err => console.error(err));
  }, []);

  return (
    <div className="p-4">
      <h1 className="text-2xl mb-4">Agents</h1>
      <Link to="/agents/new" className="text-blue-600">Create Agent</Link>
      <ul className="mt-4 list-disc list-inside">
        {agents.map(a => (
          <li key={a.id}>{a.title}</li>
        ))}
      </ul>
    </div>
  );
};

export default AgentsDashboard;
