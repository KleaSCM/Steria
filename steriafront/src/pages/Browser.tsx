// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: Browser.tsx
// Description: Main file browser page for Steria frontend, ported from Go web BrowserHandler. Uses Tailwind CSS for styling and modular components.
import React, { useEffect, useState } from 'react';
import FileList from '../components/FileList';
import UploadForm from '../components/UploadForm';

interface FileEntry {
  name: string;
  isDir: boolean;
  link: string;
}

const Browser: React.FC = () => {
  const [files, setFiles] = useState<FileEntry[]>([]);
  const [msg, setMsg] = useState('');

  useEffect(() => {
    // TODO: Replace with real API call
    setFiles([
      { name: 'Documents', isDir: true, link: 'Documents' },
      { name: 'test.txt', isDir: false, link: 'test.txt' },
    ]);
  }, []);

  return (
    <div className="min-h-screen bg-pink-50 flex flex-col items-center py-8">
      <div className="w-full max-w-3xl bg-white rounded-lg shadow-md p-8">
        <h1 className="text-3xl font-bold text-pink-700 mb-6 text-center">Steria File Browser</h1>
        <FileList files={files} />
        <UploadForm onUpload={() => setMsg('File uploaded!')} />
        {msg && <div className="text-green-600 mt-4 text-center">{msg}</div>}
      </div>
    </div>
  );
};

export default Browser; 