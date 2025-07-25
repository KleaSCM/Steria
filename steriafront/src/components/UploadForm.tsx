// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: UploadForm.tsx
// Description: File upload form for Steria file browser. Uses Tailwind CSS for styling.
import React, { useRef } from 'react';

interface UploadFormProps {
  onUpload: () => void;
}

const UploadForm: React.FC<UploadFormProps> = ({ onUpload }) => {
  const fileInput = useRef<HTMLInputElement>(null);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // TODO: Implement real upload logic
    if (fileInput.current && fileInput.current.files && fileInput.current.files.length > 0) {
      // Simulate upload
      setTimeout(() => onUpload(), 500);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="flex items-center space-x-4 mt-4">
      <input
        type="file"
        ref={fileInput}
        className="border border-pink-200 rounded px-2 py-1 focus:outline-none focus:ring-2 focus:ring-pink-400"
      />
      <button
        type="submit"
        className="bg-pink-600 text-white px-4 py-2 rounded hover:bg-pink-700 transition"
      >
        Upload
      </button>
    </form>
  );
};

export default UploadForm; 