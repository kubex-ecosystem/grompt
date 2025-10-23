import { Bot, LayoutDashboard, LucideIcon, MessageCircle, NotebookPen, Sparkles, Wand2, Workflow } from 'lucide-react';
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
  collapsed?: boolean;
}

const defaultIcons: Record<string, LucideIcon> = {
  welcome: LayoutDashboard,
  prompt: Sparkles,
  agents: Bot,
  chat: MessageCircle,
  summarizer: NotebookPen,
  code: Workflow,
  images: Wand2,
};

const Sidebar: React.FC<SidebarProps> = ({ sections, activeSection, onSectionChange, onClose, collapsed = false }) => {
  return (
    <div className="flex h-full flex-col">
      <div className="flex items-center justify-between px-3 py-5 lg:hidden">
        <h2 className="text-base font-semibold tracking-wide text-[#0f172a] dark:text-[#e5f2f2]">
          Workspace
        </h2>
        <button
          type="button"
          className="rounded-full border border-[#e2e8f0] bg-white/60 p-2 text-[#475569] transition hover:bg-white dark:border-[#13263a] dark:bg-[#0a1523]/70 dark:text-[#94a3b8]"
          onClick={onClose}
          aria-label="Close navigation"
        >
          <span className="sr-only">Close navigation</span>
          ×
        </button>
      </div>

      <div className={`hidden px-6 pt-8 pb-5 ${collapsed ? 'lg:hidden' : 'lg:block'}`}>
        <p className="text-xs uppercase tracking-[0.4em] text-[#94a3b8] dark:text-[#475569]">Kubex Suite</p>
        <h1 className="mt-3 text-2xl font-orbitron font-semibold text-[#111827] dark:text-[#f5f3ff]">Grompt Hub</h1>
        <p className="mt-2 text-sm text-[#475569] dark:text-[#94a3b8]">
          Navigate through AI-assisted modules tailored for prompt engineering, ideation, and delivery.
        </p>
      </div>

      <div className={`hidden items-center justify-center gap-3 px-3 pt-6 pb-4 ${collapsed ? 'lg:flex' : 'lg:hidden'}`}>
        <div className="flex h-12 w-12 items-center justify-center rounded-2xl border border-[#e2e8f0] bg-white text-[#06b6d4] shadow-sm dark:border-[#13263a] dark:bg-[#0f172a]/70 dark:text-[#38cde4]">
          <span className="font-orbitron text-base font-semibold">G</span>
        </div>
      </div>

      <nav className={`flex-1 overflow-y-auto px-4 pb-6 ${collapsed ? 'lg:hidden' : 'lg:px-6'}`}>
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
                      ? 'border-[#06b6d4]/70 bg-[#ecfeff] text-[#0f172a] shadow-soft-card dark:border-[#06b6d4]/70 dark:bg-[#0a1523]/70 dark:text-[#e5f2f2]'
                      : 'border-transparent bg-white/75 text-[#475569] hover:border-[#bae6fd] hover:bg-white dark:bg-[#0a1523]/40 dark:text-[#94a3b8] dark:hover:border-[#13263a]'
                    }`}
                >
                  <div className="flex items-center gap-3">
                    <span className={`flex h-10 w-10 items-center justify-center rounded-lg border ${isActive
                        ? 'border-transparent bg-white text-[#06b6d4] dark:bg-[#0a1523]/80 dark:text-[#06b6d4]'
                        : 'border-[#e2e8f0] bg-white text-[#475569] dark:border-[#13263a] dark:bg-[#0a1523]/80'
                      }`}>
                      <Icon size={20} />
                    </span>
                    <div>
                      <p className={`text-sm font-semibold ${isActive ? 'text-inherit' : 'text-[#1f2937] dark:text-[#e5f2f2]'}`}>
                        {section.label}
                      </p>
                      {section.description && (
                        <p className="text-xs text-[#64748b] dark:text-[#94a3b8]">{section.description}</p>
                      )}
                    </div>
                  </div>
                </button>
              </li>
            );
          })}
        </ul>
      </nav>

      <nav className={`hidden flex-1 overflow-y-auto pb-6 ${collapsed ? 'lg:flex' : 'lg:hidden'}`}>
        <ul className="flex w-full flex-col items-center gap-3">
          {sections.map((section) => {
            const Icon = section.icon || defaultIcons[section.id] || LayoutDashboard;
            const isActive = activeSection === section.id;
            return (
              <li key={section.id} className="w-full">
                <button
                  type="button"
                  title={section.label}
                  aria-label={section.label}
                  onClick={() => {
                    onSectionChange(section.id);
                    if (onClose) onClose();
                  }}
                  className={`mx-auto flex h-12 w-12 items-center justify-center rounded-2xl border text-[#475569] transition-all duration-200 ${isActive
                      ? 'border-[#06b6d4] bg-[#ecfeff] text-[#06b6d4] shadow-soft-card dark:border-[#06b6d4] dark:bg-[#0a1523]/70 dark:text-[#38cde4]'
                      : 'border-transparent bg-white/80 hover:border-[#bae6fd] hover:bg-white dark:border-transparent dark:bg-[#0a1523]/40 dark:text-[#94a3b8] dark:hover:border-[#13263a]'
                    }`}
                >
                  <Icon size={20} />
                </button>
              </li>
            );
          })}
        </ul>
      </nav>

      <div className={`hidden border-t border-[#e2e8f0] px-6 py-5 text-xs text-[#64748b] dark:border-[#13263a] dark:text-[#94a3b8] ${collapsed ? 'lg:hidden' : 'lg:block'}`}>
        <p>Build better prompts with Kubex governance.</p>
        <p className="mt-1">Lightweight, portable, and embeddable.</p>
      </div>
    </div>
  );
};

export default Sidebar;
