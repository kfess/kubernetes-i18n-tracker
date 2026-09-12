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

export const useFetchTranslationArticles = (
  articleCategory: ArticleCategory
): ArticleTranslation[] => {
  return articles[articleCategory].articles;
};
