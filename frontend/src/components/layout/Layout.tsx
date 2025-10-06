import React from 'react';

interface LayoutProps {
  sidebar: React.ReactNode;
  header: React.ReactNode;
  footer: React.ReactNode;
  children: React.ReactNode;
  sidebarOpen: boolean;
  onSidebarClose: () => void;
}

/**
 * Application shell shared between light/dark themes.
 * Keeps sidebar responsive while maintaining a consistent background surface.
 */
const Layout: React.FC<LayoutProps> = ({
  sidebar,
  header,
  footer,
  children,
  sidebarOpen,
  onSidebarClose,
}) => {
  return (
    <div className="min-h-screen bg-slate-100 dark:bg-[#010409] text-slate-800 dark:text-[#e0f7fa]">
      {/* mobile overlay */}
      {sidebarOpen && (
        <button
          type="button"
          className="fixed inset-0 z-30 bg-slate-900/50 backdrop-blur-sm lg:hidden"
          onClick={onSidebarClose}
          aria-label="Close navigation sidebar"
        />
      )}

      <div className="relative flex min-h-screen">
        <aside
          className={`fixed inset-y-0 left-0 z-40 w-80 transform border-r border-slate-200/80 dark:border-slate-800/60 bg-white/90 dark:bg-[#0a0f14]/95 backdrop-blur-xl transition-transform duration-300 ease-in-out lg:static lg:translate-x-0 ${sidebarOpen ? 'translate-x-0 shadow-xl shadow-slate-900/20' : '-translate-x-full'
            }`}
        >
          {sidebar}
        </aside>

        <div className="flex flex-1 flex-col lg:ml-5">
          <header className="sticky top-0 z-20 backdrop-blur-lg bg-white/80 dark:bg-[#010409]/80 border-b border-slate-200/80 dark:border-slate-800/60">
            {header}
          </header>

          <main className="flex-1 px-4 py-6 sm:px-6 lg:px-10">
            {children}
          </main>

          <footer className="border-t border-slate-200/80 dark:border-slate-800/60 bg-white/50 dark:bg-[#010409]/80">
            {footer}
          </footer>
        </div>
      </div>
    </div>
  );
};

export default Layout;
