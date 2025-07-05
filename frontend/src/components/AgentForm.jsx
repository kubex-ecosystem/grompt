import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

const AgentForm = () => {
  const [title, setTitle] = useState('');
  const [role, setRole] = useState('');
  const [skills, setSkills] = useState('');
  const [restrictions, setRestrictions] = useState('');
  const [promptExample, setPromptExample] = useState('');
  const navigate = useNavigate();

  const handleSubmit = e => {
    e.preventDefault();
    fetch('/api/agents', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        Title: title,
        Role: role,
        Skills: skills.split(',').map(s => s.trim()).filter(Boolean),
        Restrictions: restrictions.split(',').map(s => s.trim()).filter(Boolean),
        PromptExample: promptExample
      })
    }).then(res => {
      if (res.ok) navigate('/agents');
    });
  };

  return (
    <div className="p-4">
      <h1 className="text-2xl mb-4">New Agent</h1>
      <form onSubmit={handleSubmit} className="flex flex-col gap-2 max-w-md">
        <input value={title} onChange={e => setTitle(e.target.value)} placeholder="Title" className="border p-1" required />
        <input value={role} onChange={e => setRole(e.target.value)} placeholder="Role" className="border p-1" />
        <input value={skills} onChange={e => setSkills(e.target.value)} placeholder="Skills (comma separated)" className="border p-1" />
        <input value={restrictions} onChange={e => setRestrictions(e.target.value)} placeholder="Restrictions (comma separated)" className="border p-1" />
        <textarea value={promptExample} onChange={e => setPromptExample(e.target.value)} placeholder="Prompt Example" className="border p-1" />
        <button type="submit" className="bg-blue-600 text-white px-2 py-1">Save</button>
      </form>
    </div>
  );
};

export default AgentForm;
