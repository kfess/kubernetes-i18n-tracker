import { startTransition, useEffect, useMemo, useRef, useState } from 'react';
import { Container, Stack, Text } from '@mantine/core';
import { useLocalStorage, useMediaQuery } from '@mantine/hooks';
import { ArticleCategorySelector } from '@/features/ArticleCategorySelector';
import { ArticleListControl } from '@/features/ArticleListControl';
import { useFetchTranslationArticles } from '@/features/hooks/useFetchTranslationArticles';
import { type LanguageCode, type LanguageCodeWithAll } from '@/features/language/languageCodes';
import { LocalizationPolicyWarning } from '@/features/LocalizationPolicyWarning';
import { MobileTranslationStatusMatrix } from '@/features/MobileTranslationStatusMatrix';
import { type ArticleCategory, type TranslationStatus } from '@/features/translations';
import { TranslationStatusMatrix } from '@/features/TranslationStatusMatrix';
import {
  IssueStatus,
  PrStatus,
  sortModes,
  type SortDirection,
  type SortMode,
} from '@/features/types';
import { useQueryParams, useUpdateQueryParams } from '@/hooks/useQueryParams';
import { getDeploymentInfo } from '@/utils/deploy';

export function HomePage() {
  const isMobile = useMediaQuery('(max-width: 768px)');
  const { deployedAt, gitCommit } = getDeploymentInfo();

  const [selectedArticleCategory, setSelectedArticleCategory] = useQueryParams<ArticleCategory>(
    'category',
    'docsConcept'
  );
  const translationArticles = useFetchTranslationArticles(selectedArticleCategory);

  // Pagination state
  const [activePage, setActivePage] = useQueryParams<number>('page', 1, String, (value) => {
    const parsed = parseInt(value, 10);
    return isNaN(parsed) || parsed < 1 ? 1 : parsed;
  });
  const [itemsPerPage, setItemsPerPage] = useQueryParams<string>('itemsPerPage', '30');

  // Filter states
  const [languageFilter, setLanguageFilter] = useQueryParams<LanguageCodeWithAll>(
    'lang',
    'all',
    String,
    (value) => (value as LanguageCodeWithAll) || 'all'
  );
  const [statusFilter, setStatusFilter] = useQueryParams<TranslationStatus | 'all'>(
    'status',
    'all',
    String,
    (value) => (value as TranslationStatus) || 'all'
  );
  const [issueFilter, setIssueFilter] = useQueryParams<IssueStatus>(
    'issue',
    'all',
    String,
    (value) => (value as IssueStatus) || 'all'
  );
  const [prFilter, setPrFilter] = useQueryParams<PrStatus>(
    'pr',
    'all',
    String,
    (value) => (value as PrStatus) || 'all'
  );
  const [searchQuery, setSearchQuery] = useQueryParams<string>(
    'search',
    '',
    String,
    (value) => value || ''
  );

  const [debouncedSearchQuery, setDebouncedSearchQuery] = useState<string>(searchQuery);

  // Sort states
  const [sortMode] = useQueryParams<SortMode>('sortMode', 'default', String, (value) =>
    sortModes.includes(value as SortMode) ? (value as SortMode) : 'default'
  );
  const [sortDirection] = useQueryParams<SortDirection>('sortDirection', 'desc', String, (value) =>
    value === 'asc' || value === 'desc' ? value : 'desc'
  );

  // Selected languages from localStorage
  const [selectedLanguages] = useLocalStorage<LanguageCode[]>({
    key: 'selected-languages',
  });

  // Track if this is the initial mount
  const isInitialMount = useRef(true);
  useEffect(() => {
    if (isInitialMount.current) {
      isInitialMount.current = false;
      return;
    }
    setActivePage(1);
  }, [
    selectedArticleCategory,
    statusFilter,
    languageFilter,
    issueFilter,
    prFilter,
    debouncedSearchQuery,
    itemsPerPage,
    sortMode,
    sortDirection,
  ]);

  const getFilteredArticles = () => {
    let filtered = translationArticles;

    if (sortMode === 'views') {
      filtered = [...filtered].sort((a, b) => {
        const aViews = a.translations.en.views || 0;
        const bViews = b.translations.en.views || 0;
        return sortDirection === 'asc' ? aViews - bViews : bViews - aViews;
      });
    } else if (sortMode === 'newUsers') {
      filtered = [...filtered].sort((a, b) => {
        const aUsers = a.translations.en.newUsers || 0;
        const bUsers = b.translations.en.newUsers || 0;
        return sortDirection === 'asc' ? aUsers - bUsers : bUsers - aUsers;
      });
    } else if (sortMode === 'updatedAt') {
      filtered = [...filtered].sort((a, b) => {
        const aDate = new Date(a.translations.en.englishLatestDate);
        const bDate = new Date(b.translations.en.englishLatestDate);
        return sortDirection === 'asc'
          ? aDate.getTime() - bDate.getTime()
          : bDate.getTime() - aDate.getTime();
      });
    } else if (sortMode === 'averageSessionDuration') {
      filtered = [...filtered].sort((a, b) => {
        const aDuration = a.translations.en.averageSessionDuration || 0;
        const bDuration = b.translations.en.averageSessionDuration || 0;
        return sortDirection === 'asc' ? aDuration - bDuration : bDuration - aDuration;
      });
    } else {
      // eslint-disable-next-line no-lonely-if
      if (sortDirection === 'asc') {
        filtered = [...filtered].reverse();
      }
    }

    if (statusFilter !== 'all') {
      filtered = filtered.filter((article) => {
        if (languageFilter === 'all') {
          return Object.values(article.translations).some(
            (translation) => translation.status === statusFilter
          );
        }

        const translation = article.translations[languageFilter];
        return translation && translation.status === statusFilter;
      });
    }

    if (languageFilter !== 'all' && statusFilter === 'all') {
      filtered = filtered.filter((article) => {
        const translation = article.translations[languageFilter];
        return translation !== undefined;
      });
    }

    if (issueFilter !== 'all') {
      filtered = filtered.filter((article) => {
        if (languageFilter === 'all') {
          const hasIssues = Object.values(article.translations).some(
            (translation) => translation.issues.length > 0
          );
          return issueFilter === 'withIssues' ? hasIssues : !hasIssues;
        }
        const hasIssues = article.translations[languageFilter]?.issues.length > 0;
        return issueFilter === 'withIssues' ? hasIssues : !hasIssues;
      });
    }

    if (prFilter !== 'all') {
      filtered = filtered.filter((article) => {
        if (languageFilter === 'all') {
          const hasPR = Object.values(article.translations).some(
            (translation) => translation.prs.length > 0
          );
          return prFilter === 'withPr' ? hasPR : !hasPR;
        }
        const hasPR = article.translations[languageFilter]?.prs.length > 0;
        return prFilter === 'withPr' ? hasPR : !hasPR;
      });
    }

    if (debouncedSearchQuery) {
      const lowerQuery = debouncedSearchQuery.toLowerCase();
      filtered = filtered.filter((article) =>
        article.englishPath.toLowerCase().includes(lowerQuery)
      );
    }

    return filtered;
  };

  const filteredArticles = useMemo(
    () => getFilteredArticles(),
    [
      translationArticles,
      statusFilter,
      languageFilter,
      issueFilter,
      prFilter,
      debouncedSearchQuery,
      sortMode,
      sortDirection,
    ]
  );

  const startIndex = (activePage - 1) * parseInt(itemsPerPage, 10);
  const endIndex = Math.min(startIndex + parseInt(itemsPerPage, 10), filteredArticles.length);
  const currentArticles = filteredArticles.slice(startIndex, endIndex);

  const onArticleCategoryChange = (category: ArticleCategory) => {
    startTransition(() => {
      setSelectedArticleCategory(category);
    });
  };

  const updateQueryParams = useUpdateQueryParams();
  const handleSortChange = (mode: SortMode, direction: SortDirection) => {
    updateQueryParams({
      sortMode: mode,
      sortDirection: direction,
    });
  };

  return (
    <Container fluid px={{ base: '0', sm: 'md' }} mt={{ base: 'xs', sm: 'md' }}>
      <LocalizationPolicyWarning selectedLanguages={selectedLanguages} />
      <Stack gap="sm">
        <ArticleCategorySelector
          articleCategory={selectedArticleCategory}
          onArticleCategoryChange={onArticleCategoryChange}
        />
        <ArticleListControl
          selectedArticleCategory={selectedArticleCategory}
          articles={translationArticles}
          filteredArticles={filteredArticles}
          activePage={activePage}
          setActivePage={setActivePage}
          itemsPerPage={itemsPerPage}
          setItemsPerPage={setItemsPerPage}
          statusFilter={statusFilter}
          setStatusFilter={setStatusFilter}
          languageFilter={languageFilter}
          setLanguageFilter={setLanguageFilter}
          issueFilter={issueFilter}
          setIssueFilter={setIssueFilter}
          prFilter={prFilter}
          setPrFilter={setPrFilter}
          searchQuery={searchQuery}
          setSearchQuery={setSearchQuery}
          debouncedSearchQuery={debouncedSearchQuery}
          setDebouncedSearchQuery={setDebouncedSearchQuery}
          sortMode={sortMode}
          sortDirection={sortDirection}
          onSortChange={handleSortChange}
          selectedLanguages={selectedLanguages || []}
          startIndex={startIndex}
          endIndex={endIndex}
        />
        {!isMobile ? (
          <TranslationStatusMatrix
            articles={currentArticles}
            languageFilter={languageFilter}
            selectedLanguages={selectedLanguages || []}
            selectedArticleCategory={selectedArticleCategory}
          />
        ) : (
          <MobileTranslationStatusMatrix
            articles={currentArticles}
            selectedLanguages={selectedLanguages || []}
            selectedArticleCategory={selectedArticleCategory}
          />
        )}
        {deployedAt && (
          <Text size="xs" c="dimmed" ta="right">
            Last Updated at {deployedAt} ({gitCommit})
          </Text>
        )}
      </Stack>
    </Container>
  );
}
