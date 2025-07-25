// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: CommitGraph.tsx
// Description: Renders a vertical commit graph for the Steria file browser. Uses Tailwind CSS for styling.
import React from 'react';

export interface Commit {
  hash: string;
  author: string;
  timestamp: string;
  message: string;
  parent?: string;
}

interface CommitGraphProps {
  commits: Commit[];
}

const CommitGraph: React.FC<CommitGraphProps> = ({ commits }) => (
  <div className="commit-graph max-h-96 overflow-y-auto font-mono">
    {commits.map((commit) => (
      <div key={commit.hash} className="mb-6 border-l-4 border-pink-700 pl-4 relative">
        <div className="absolute -left-3 top-0 text-pink-700 text-xl">●</div>
        <div className="text-xs text-pink-900 font-bold">{commit.hash.substring(0,8)}</div>
        <div className="text-green-700">{commit.author}</div>
        <div className="text-gray-500 text-xs">{new Date(commit.timestamp).toLocaleString()}</div>
        <div className="font-bold text-pink-700">{commit.message}</div>
        {commit.parent && (
          <div className="text-xs text-gray-400">Parent: {commit.parent.substring(0,8)}</div>
        )}
      </div>
    ))}
  </div>
);

export default CommitGraph; 