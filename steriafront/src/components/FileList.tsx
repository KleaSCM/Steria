// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: FileList.tsx
// Description: Renders a list of files and directories for the Steria file browser. Uses Tailwind CSS for styling.
import React from 'react';

export interface FileEntry {
  name: string;
  isDir: boolean;
  link: string;
}

interface FileListProps {
  files: FileEntry[];
}

const FileList: React.FC<FileListProps> = ({ files }) => (
  <ul className="mb-6">
    {files.map((entry) => (
      <li key={entry.link} className="mb-2">
        {entry.isDir ? (
          <a href={`?path=${encodeURIComponent(entry.link)}`} className="flex items-center text-pink-700 hover:underline">
            <span className="mr-2">📁</span>
            <span>{entry.name}</span>
          </a>
        ) : (
          <a href={`/download?path=${encodeURIComponent(entry.link)}`} className="flex items-center text-pink-900 hover:underline">
            <span className="mr-2">📄</span>
            <span>{entry.name}</span>
          </a>
        )}
      </li>
    ))}
  </ul>
);

export default FileList; 