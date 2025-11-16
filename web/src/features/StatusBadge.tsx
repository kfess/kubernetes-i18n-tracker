import { Tooltip } from '@mantine/core';
import { type TranslationStatus } from '@/features/translations';

export const StatusBadge = ({ status }: { status: TranslationStatus }) => {
  const getStatusConfig = (status: TranslationStatus) => {
    switch (status) {
      case 'up_to_date':
        return { emoji: '✅', label: 'Up to date' };
      case 'possibly_outdated':
        return { emoji: '✅ ⚠️', label: "Possibly outdated (Document's header structures differ)" };
      case 'outdated':
        return { emoji: '⚠️', label: 'Outdated' };
      case 'not_translated':
        return { emoji: '—', label: 'Not translated' };
      default:
        return { emoji: '-', label: 'Unknown' };
    }
  };

  const statusConfig = getStatusConfig(status);

  return (
    <Tooltip label={statusConfig.label} withArrow>
      <span>{statusConfig.emoji}</span>
    </Tooltip>
  );
};
