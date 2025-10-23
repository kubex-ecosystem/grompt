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
      className={`rounded-2xl border border-[#e2e8f0] bg-white p-6 shadow-soft-card transition dark:border-[#0b1220] dark:bg-[#0a1523]/85 ${
        className ?? ''
      }`}
    >
      {(title || description || action) && (
        <header className="mb-4 flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between">
          <div>
            {title && <h3 className="text-lg font-semibold text-[#111827] dark:text-[#e5f2f2]">{title}</h3>}
            {description && <p className="text-sm text-[#475569] dark:text-[#94a3b8]">{description}</p>}
          </div>
          {action && <div className="flex-shrink-0">{action}</div>}
        </header>
      )}
      <div>{children}</div>
    </section>
  );
};

export default Card;
