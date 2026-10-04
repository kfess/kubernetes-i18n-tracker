import { useMemo } from 'react';
import blogArticles from '@/data/output/matrix/blog.json';
import careerArticles from '@/data/output/matrix/careers.json';
import communityArticles from '@/data/output/matrix/community.json';
import docsConceptsArticles from '@/data/output/matrix/docs_concepts.json';
import docsContributionArticles from '@/data/output/matrix/docs_contribute.json';
import docsHomeArticles from '@/data/output/matrix/docs_home.json';
import docsReferenceArticles from '@/data/output/matrix/docs_reference.json';
import docsSetupArticles from '@/data/output/matrix/docs_setup.json';
import docsTasksArticles from '@/data/output/matrix/docs_tasks.json';
import docsTutorialsArticles from '@/data/output/matrix/docs_tutorials.json';
import examplesArticles from '@/data/output/matrix/examples.json';
import includesArticles from '@/data/output/matrix/includes.json';
import partnerArticles from '@/data/output/matrix/partners.json';
import releaseArticles from '@/data/output/matrix/releases.json';
import trainingArticles from '@/data/output/matrix/training.json';
import { applyDetectionMode } from '@/features/deriveStatus';
import { useDetectionMode } from '@/features/hooks/useDetectionMode';
import {
  ArticleCategory,
  ArticleTranslation,
  TranslationStatusReport,
} from '@/features/translations';

const articles = {
  blog: blogArticles as TranslationStatusReport,
  community: communityArticles as TranslationStatusReport,
  examples: examplesArticles as TranslationStatusReport,
  docsConcept: docsConceptsArticles as TranslationStatusReport,
  docsContribute: docsContributionArticles as TranslationStatusReport,
  docsHome: docsHomeArticles as TranslationStatusReport,
  docsTask: docsTasksArticles as TranslationStatusReport,
  docsReference: docsReferenceArticles as TranslationStatusReport,
  docsSetup: docsSetupArticles as TranslationStatusReport,
  docsTutorial: docsTutorialsArticles as TranslationStatusReport,
  includes: includesArticles as TranslationStatusReport,
  partner: partnerArticles as TranslationStatusReport,
  release: releaseArticles as TranslationStatusReport,
  training: trainingArticles as TranslationStatusReport,
  career: careerArticles as TranslationStatusReport,
};

// Both hooks return the data with each translation's status decided by the
// selected detection mode, so consumers can keep reading `status` directly.

export const useFetchTranslationArticles = (
  articleCategory: ArticleCategory
): ArticleTranslation[] => {
  const [detectionMode] = useDetectionMode();

  return useMemo(
    () => applyDetectionMode(articles[articleCategory].articles, detectionMode),
    [articleCategory, detectionMode]
  );
};

export const useFetchAllTranslationArticles = (): Record<
  ArticleCategory,
  TranslationStatusReport
> => {
  const [detectionMode] = useDetectionMode();

  return useMemo(
    () =>
      Object.fromEntries(
        Object.entries(articles).map(([category, report]) => [
          category,
          { ...report, articles: applyDetectionMode(report.articles, detectionMode) },
        ])
      ) as Record<ArticleCategory, TranslationStatusReport>,
    [detectionMode]
  );
};
