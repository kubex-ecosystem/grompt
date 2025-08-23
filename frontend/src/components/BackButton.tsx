import { ArrowLeft } from 'lucide-react';
import { useRouter } from 'next/navigation';
import { useTranslation } from 'react-i18next';
import { logger } from '../utils/logger';

interface BackButtonProps {
  to?: string;
  currentTheme: any;
  label?: string;
}

const BackButton = ({ to = '/', currentTheme, label }: BackButtonProps) => {
  const { t } = useTranslation();
  const router = useRouter();

  const handleBack = () => {
    const currentPath = window.location.pathname;

    if (to === 'back') {
      logger.logNavigation(currentPath, 'history-back', 'back');
      router.back();
    } else {
      logger.logNavigation(currentPath, to, 'push');
      router.push(to);
    }
  };

  return (
    <button
      title={label || t('common.back')}
      onClick={handleBack}
      className={`flex items-center gap-2 px-3 py-2 rounded-lg ${currentTheme.buttonSecondary} transition-colors hover:scale-105`}
    >
      <ArrowLeft size={16} />
      <span className="hidden sm:inline text-sm">{label || t('common.back')}</span>
    </button>
  );
};

export default BackButton;
