import React from 'react';

interface CardProps {
  title?: string;
  description?: string;
  action?: React.ReactNode;
  children: React.ReactNode;
  className?: string;
}

const Card: React.FC<CardProps> = ({ title, description, action, children, className }) => {
  return (
    <section
      className={`rounded-2xl border border-slate-200/80 bg-white/80 p-6 shadow-[0_28px_60px_-45px_rgba(15,23,42,0.45)] transition dark:border-slate-800/60 dark:bg-[#0a0f14]/70 ${
        className ?? ''
      }`}
    >
      {(title || description || action) && (
        <header className="mb-4 flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between">
          <div>
            {title && <h3 className="text-lg font-semibold text-slate-900 dark:text-[#e0f7fa]">{title}</h3>}
            {description && <p className="text-sm text-slate-500 dark:text-slate-400">{description}</p>}
          </div>
          {action && <div className="flex-shrink-0">{action}</div>}
        </header>
      )}
      <div>{children}</div>
    </section>
  );
};

export default Card;
