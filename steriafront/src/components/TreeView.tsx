// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: TreeView.tsx
// Description: Collapsible directory tree for Steria file browser. Uses Tailwind CSS for styling.
import React, { useState } from 'react';

export interface TreeNode {
  name: string;
  children?: TreeNode[];
}

interface TreeViewProps {
  tree: TreeNode;
  basePath?: string;
}

const TreeView: React.FC<TreeViewProps> = ({ tree, basePath = '' }) => {
  const [expanded, setExpanded] = useState<{ [key: string]: boolean }>({});

  const toggle = (path: string) => {
    setExpanded((prev) => ({ ...prev, [path]: !prev[path] }));
  };

  const renderNode = (node: TreeNode, path: string) => (
    <div key={path} className="ml-2">
      <div className="flex items-center cursor-pointer" onClick={() => node.children && toggle(path)}>
        <span className="mr-1">{node.children ? '📁' : '📄'}</span>
        <a href={`?path=${encodeURIComponent(path)}`} className="text-pink-700 hover:underline">{node.name}</a>
      </div>
      {node.children && expanded[path] && (
        <div className="ml-4">
          {node.children.map((child) => renderNode(child, path + '/' + child.name))}
        </div>
      )}
    </div>
  );

  return <div>{renderNode(tree, basePath + '/' + tree.name)}</div>;
};

export default TreeView; 