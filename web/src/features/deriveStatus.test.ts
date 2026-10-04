import { describe, expect, it } from 'vitest';
import sample from '@/__fixtures__/docs_concepts.sample.json';
import { applyDetectionMode, deriveStatus } from '@/features/deriveStatus';
import { type LanguageCode } from '@/features/language/languageCodes';
import {
  type ArticleTranslation,
  type TranslationInfo,
  type TranslationStatus,
  type TranslationStatusReport,
} from '@/features/translations';
import { detectionModes, type DetectionMode } from '@/features/types';

const articles = (sample as TranslationStatusReport).articles;

const statusOf = (list: ArticleTranslation[], page: string, lang: LanguageCode) => {
  const article = list.find((a) => a.englishPath.endsWith(`/${page}.md`));
  if (!article) {
    throw new Error(`no sample article named ${page}`);
  }
  return article.translations[lang].status;
};

// [page, language, combined, git, structure]
const cases: [string, LanguageCode, TranslationStatus, TranslationStatus, TranslationStatus][] = [
  // git up to date, structure none / moderate / strong
  ['1-git-up-to-date', 'ja', 'up_to_date', 'up_to_date', 'up_to_date'],
  ['1-git-up-to-date', 'ko', 'possibly_outdated', 'up_to_date', 'possibly_outdated'],
  ['1-git-up-to-date', 'es', 'possibly_outdated', 'up_to_date', 'outdated'],
  ['1-git-up-to-date', 'de', 'not_translated', 'not_translated', 'not_translated'],
  // git outdated, structure none / moderate / strong
  ['2-git-outdated', 'ja', 'outdated', 'outdated', 'up_to_date'],
  ['2-git-outdated', 'ko', 'outdated', 'outdated', 'possibly_outdated'],
  ['2-git-outdated', 'es', 'outdated', 'outdated', 'outdated'],
  // empty translation, and translations without a structural comparison
  ['3-edge-cases', 'ja', 'possibly_outdated', 'up_to_date', 'outdated'],
  ['3-edge-cases', 'es', 'up_to_date', 'up_to_date', 'up_to_date'],
  ['3-edge-cases', 'ko', 'outdated', 'outdated', 'outdated'],
  // English is never compared with itself
  ['3-edge-cases', 'en', 'up_to_date', 'up_to_date', 'up_to_date'],
];

describe('applyDetectionMode', () => {
  const modes: DetectionMode[] = ['combined', 'git', 'structure'];

  it.each(cases)('%s / %s', (page, lang, ...expected) => {
    const actual = modes.map((mode) => statusOf(applyDetectionMode(articles, mode), page, lang));
    expect(actual).toEqual(expected);
  });

  it('covers every detection mode', () => {
    expect([...detectionModes].sort()).toEqual([...modes].sort());
  });

  it('leaves the input untouched', () => {
    const before = JSON.stringify(articles);
    applyDetectionMode(articles, 'structure');
    expect(JSON.stringify(articles)).toBe(before);
  });

  it('returns the same array in the default mode', () => {
    expect(applyDetectionMode(articles, 'combined')).toBe(articles);
  });
});

describe('deriveStatus', () => {
  it('keeps the exported status for data without gitStatus', () => {
    const legacy = { status: 'possibly_outdated' } as TranslationInfo;
    expect(detectionModes.map((mode) => deriveStatus(legacy, mode))).toEqual([
      'possibly_outdated',
      'possibly_outdated',
      'possibly_outdated',
    ]);
  });
});
