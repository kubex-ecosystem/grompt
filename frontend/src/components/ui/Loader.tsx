import { Loader2 } from 'lucide-react';
import React from 'react';

interface LoaderProps {
  label?: string;
  className?: string;
}

const Loader: React.FC<LoaderProps> = ({ label, className }) => (
  <div className={`flex items-center gap-3 text-slate-500 dark:text-slate-300 ${className ?? ''}`}>
    <Loader2 className="h-5 w-5 animate-spin" />
    {label && <span className="text-sm font-medium">{label}</span>}
  </div>
);

export default Loader;
