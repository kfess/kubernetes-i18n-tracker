import {
  type ArticleTranslation,
  type StructuralSignal,
  type TranslationInfo,
  type TranslationStatus,
} from '@/features/translations';
import { defaultDetectionMode, type DetectionMode } from '@/features/types';

const statusBySignal: Record<StructuralSignal, TranslationStatus> = {
  strong: 'outdated',
  moderate: 'possibly_outdated',
  none: 'up_to_date',
};

/**
 * Decides the status of one translation under the given detection mode.
 *
 * - combined: the status exported by the tracker (git decides, structure can
 *   flag an up-to-date page as possibly outdated).
 * - git: git history alone.
 * - structure: the structural signal alone. Without a comparison there is
 *   nothing to judge by, so it falls back to git; that also keeps untranslated
 *   pages untranslated.
 */
export const deriveStatus = (info: TranslationInfo, mode: DetectionMode): TranslationStatus => {
  // Data exported before detection modes only carries the combined status.
  if (info.gitStatus === undefined) {
    return info.status;
  }

  switch (mode) {
    case 'git':
      return info.gitStatus;
    case 'structure':
      return info.structure ? statusBySignal[info.structure.signal] : info.gitStatus;
    default:
      return info.status;
  }
};

/**
 * Returns the articles with every translation's status decided by the mode.
 * The input is not modified; in the default mode it is returned as is.
 */
export const applyDetectionMode = (
  articles: ArticleTranslation[],
  mode: DetectionMode
): ArticleTranslation[] => {
  if (mode === defaultDetectionMode) {
    return articles;
  }

  return articles.map((article) => ({
    ...article,
    translations: Object.fromEntries(
      Object.entries(article.translations).map(([langCode, info]) => [
        langCode,
        { ...info, status: deriveStatus(info, mode) },
      ])
    ) as ArticleTranslation['translations'],
  }));
};
