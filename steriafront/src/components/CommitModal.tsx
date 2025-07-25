// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: CommitModal.tsx
// Description: Modal for showing commit graph and mermaid diagram. Uses Tailwind CSS for styling.
import React, { useState } from 'react';
import CommitGraph, { type Commit } from './CommitGraph';

interface CommitModalProps {
  open: boolean;
  onClose: () => void;
  commits: Commit[];
}

const CommitModal: React.FC<CommitModalProps> = ({ open, onClose, commits }) => {
  const [tab, setTab] = useState<'graph' | 'mermaid'>('graph');

  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-40">
      <div className="bg-white rounded-lg p-8 max-w-2xl w-full relative">
        <button onClick={onClose} className="absolute top-2 right-2 text-pink-700 text-2xl">&times;</button>
        <div className="flex gap-4 mb-4">
          <button onClick={() => setTab('graph')} className={`px-4 py-2 rounded-t-lg ${tab==='graph' ? 'bg-pink-700 text-white' : 'bg-pink-100 text-pink-700'}`}>Vertical Graph</button>
          <button onClick={() => setTab('mermaid')} className={`px-4 py-2 rounded-t-lg ${tab==='mermaid' ? 'bg-pink-700 text-white' : 'bg-pink-100 text-pink-700'}`}>Mermaid Diagram</button>
        </div>
        {tab === 'graph' && <CommitGraph commits={commits} />}
        {tab === 'mermaid' && <div className="mermaid bg-pink-50 rounded p-4">Mermaid diagram placeholder</div>}
      </div>
    </div>
  );
};

export default CommitModal; 