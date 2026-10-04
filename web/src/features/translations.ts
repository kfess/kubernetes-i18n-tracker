import { type LanguageCode } from '@/features/language/languageCodes';

export type TranslationStatus = 'up_to_date' | 'outdated' | 'not_translated' | 'possibly_outdated';

export type StructuralSignal = 'none' | 'moderate' | 'strong';
export type Gap = {
  lineRatio: number;
  bodyWordRatio: number;
  missingH2: number;
  missingH3: number;
  missingCodeBlocks: number;
  missingAnchors: number;
  missingNewVersions: number;
  missingFeatureState: number;
  missingApiOrKind: number;
};
export type Structure = {
  signal: StructuralSignal;
  gap: Gap;
};

type Severity = 'current' | 'minor' | 'moderate' | 'significant' | 'critical';

interface PullRequest {
  number: number;
  title: string;
  url: string;
  files: string[];
}

interface Issue {
  number: number;
  title: string;
  url: string;
}

export interface TranslationInfo {
  // Overall translation status. Rewritten by the selected detection mode when loaded.
  status: TranslationStatus;

  // Status from git history alone. Missing in data exported before detection modes.
  gitStatus?: TranslationStatus;

  // Structural comparison with the English page.
  // Null when nothing was compared (English itself, no translation, unreadable file).
  structure?: Structure | null;

  severity: Severity;
  daysBehind: number;
  commitsBehind: number;
  totalChangeLines: number;
  targetLatestDate: string | null;
  englishLatestDate: string;
  translationUrl: string | null;
  views: number;
  newUsers: number;
  averageSessionDuration: number;
  issues: Issue[];
  prs: PullRequest[];
  englishLatestCommitHash: string | null;
  refEnglishCommitHash: string | null;
  refEnglishCommitDate: string | null;
}

export interface ArticleTranslation {
  englishPath: string;
  englishUrl: string | null;
  translations: Record<LanguageCode, TranslationInfo>;
}

export interface TranslationStatusReport {
  lastUpdated: string;
  articles: ArticleTranslation[];
}

export const articleCategories = [
  { value: 'docsConcept', label: 'Docs / Concept' },
  { value: 'docsTask', label: 'Docs / Task' },
  { value: 'docsSetup', label: 'Docs / Setup' },
  { value: 'docsReference', label: 'Docs / Reference' },
  { value: 'docsTutorial', label: 'Docs / Tutorial' },
  { value: 'docsContribute', label: 'Docs / Contribution' },
  { value: 'docsHome', label: 'Docs / Home' },
  { value: 'blog', label: 'Blog' },
  { value: 'community', label: 'Community' },
  { value: 'examples', label: 'Examples' },
  { value: 'includes', label: 'Includes' },
  { value: 'release', label: 'Release' },
  { value: 'partner', label: 'Partner' },
  { value: 'training', label: 'Training' },
  { value: 'career', label: 'Career' },
] as const;

export type ArticleCategory = (typeof articleCategories)[number]['value'];

export type Diff = {
  [key: string]: {
    englishPath: string;
    language: string;
    refEnglishCommitHash: string;
    englishLatestCommitHash: string;
    diff: string;
  };
};
