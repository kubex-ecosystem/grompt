import React from 'react';
import { Routes, Route, Link } from 'react-router-dom';
import PromptCrafter from './PromptCrafter';
import AgentsDashboard from './components/AgentsDashboard';
import AgentForm from './components/AgentForm';

const App = () => (
  <div>
    <nav className="p-4 bg-gray-800 text-white flex gap-4">
      <Link to="/">Prompt Crafter</Link>
      <Link to="/agents">Agents</Link>
    </nav>
    <Routes>
      <Route path="/" element={<PromptCrafter />} />
      <Route path="/agents" element={<AgentsDashboard />} />
      <Route path="/agents/new" element={<AgentForm />} />
    </Routes>
  </div>
);

export default App;
