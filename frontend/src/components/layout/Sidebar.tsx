import { LayoutDashboard, LucideIcon, MessageCircle, NotebookPen, Sparkles, Wand2, Workflow } from 'lucide-react';
import React from 'react';

export type SidebarSection = {
  id: string;
  label: string;
  description?: string;
  icon?: LucideIcon;
};

interface SidebarProps {
  sections: SidebarSection[];
  activeSection: string;
  onSectionChange: (section: string) => void;
  onClose?: () => void;
}

const defaultIcons: Record<string, LucideIcon> = {
  welcome: LayoutDashboard,
  prompt: Sparkles,
  chat: MessageCircle,
  summarizer: NotebookPen,
  code: Workflow,
  images: Wand2,
};

const Sidebar: React.FC<SidebarProps> = ({ sections, activeSection, onSectionChange, onClose }) => {
  return (
    <div className="flex h-full flex-col">
      <div className="flex items-center justify-between px-3 py-5 lg:hidden">
        <h2 className="text-base font-semibold tracking-wide text-slate-700 dark:text-[#e0f7fa]">
          Workspace
        </h2>
        <button
          type="button"
          className="rounded-full border border-slate-200/70 bg-white/40 p-2 text-slate-500 transition hover:bg-white/80 dark:border-slate-700/70 dark:bg-[#0a0f14]/70 dark:text-slate-300 dark:hover:border-slate-500"
          onClick={onClose}
          aria-label="Close navigation"
        >
          <span className="sr-only">Close navigation</span>
          ×
        </button>
      </div>

      <div className="hidden lg:block px-6 pt-8 pb-5">
        <p className="text-xs uppercase tracking-[0.4em] text-slate-400 dark:text-slate-500">Kubex Suite</p>
        <h1 className="mt-3 text-2xl font-orbitron font-semibold text-slate-900 dark:text-white">Grompt Hub</h1>
        <p className="mt-2 text-sm text-slate-500 dark:text-slate-400">
          Navigate through AI-assisted modules tailored for prompt engineering, ideation, and delivery.
        </p>
      </div>

      <nav className="flex-1 overflow-y-auto px-4 pb-6 lg:px-6">
        <ul className="space-y-2">
          {sections.map((section) => {
            const Icon = section.icon || defaultIcons[section.id] || LayoutDashboard;
            const isActive = activeSection === section.id;
            return (
              <li key={section.id}>
                <button
                  type="button"
                  onClick={() => {
                    onSectionChange(section.id);
                    if (onClose) onClose();
                  }}
                  className={`w-full rounded-xl border px-4 py-3 text-left transition-all duration-200 ${isActive
                      ? 'border-slate-900/80 bg-slate-900 text-white shadow-[0_20px_45px_-35px_rgba(15,23,42,0.8)] dark:border-[#00f0ff]/80 dark:bg-[#00f0ff]/10 dark:text-[#e0f7fa]'
                      : 'border-transparent bg-white/70 text-slate-600 hover:border-slate-300 hover:bg-white dark:bg-[#0a0f14]/40 dark:text-slate-300 dark:hover:border-slate-600'
                    }`}
                >
                  <div className="flex items-center gap-3">
                    <span className={`flex h-10 w-10 items-center justify-center rounded-lg border ${isActive
                        ? 'border-white/20 bg-white/10 text-white dark:border-[#00f0ff]/30 dark:bg-[#00f0ff]/20'
                        : 'border-slate-200 bg-white text-slate-600 dark:border-slate-700 dark:bg-[#0a0f14]/80'
                      }`}>
                      <Icon size={20} />
                    </span>
                    <div>
                      <p className={`text-sm font-semibold ${isActive ? 'text-inherit' : 'text-slate-700 dark:text-slate-200'}`}>
                        {section.label}
                      </p>
                      {section.description && (
                        <p className="text-xs text-slate-500 dark:text-slate-400">{section.description}</p>
                      )}
                    </div>
                  </div>
                </button>
              </li>
            );
          })}
        </ul>
      </nav>

      <div className="hidden border-t border-slate-200/70 px-6 py-5 text-xs text-slate-500 dark:border-slate-800/60 dark:text-slate-400 lg:block">
        <p>Build better prompts with Kubex governance.</p>
        <p className="mt-1">Lightweight, portable, and embeddable.</p>
      </div>
    </div>
  );
};

export default Sidebar;
