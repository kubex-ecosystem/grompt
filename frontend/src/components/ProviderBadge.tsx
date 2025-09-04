import { getProviderMeta } from '../lib/providersTheme';

type Props = {
  provider?: string;
  size?: 'sm' | 'md';
  showLabel?: boolean;
  className?: string;
};

export default function ProviderBadge({ provider, size = 'sm', showLabel = true, className = '' }: Props) {
  const meta = getProviderMeta(provider);
  const dim = size === 'sm' ? 8 : 10;
  const gap = size === 'sm' ? 'gap-1' : 'gap-2';
  const textCls = size === 'sm' ? 'text-[10px]' : 'text-xs';
  return (
    <span className={`inline-flex items-center ${gap} ${className}`} title={meta.label}>
      <span style={{ backgroundColor: meta.color, width: dim, height: dim }} className="inline-block rounded-full" />
      {showLabel && <span className={`${textCls}`} style={{ color: meta.color }}>{meta.label}</span>}
    </span>
  );
}

