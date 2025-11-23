import { IconRefresh } from '@tabler/icons-react';
import { useNavigate } from 'react-router-dom';
import {
  ActionIcon,
  Card,
  Group,
  Pagination,
  rem,
  Select,
  Stack,
  Text,
  TextInput,
} from '@mantine/core';
import { useDebouncedCallback, useMediaQuery } from '@mantine/hooks';
import {
  getSortedLangCodes,
  type LanguageCode,
  type LanguageCodeWithAll,
} from '@/features/language/languageCodes';
import { SortMenu } from '@/features/SortMenu';
import {
  type ArticleCategory,
  type ArticleTranslation,
  type TranslationStatus,
} from '@/features/translations';
import {
  type IssueStatus,
  type PrStatus,
  type SortDirection,
  type SortMode,
} from '@/features/types';

interface Props {
  selectedArticleCategory: ArticleCategory;
  articles: ArticleTranslation[];
  filteredArticles: ArticleTranslation[];
  activePage: number;
  setActivePage: (page: number) => void;
  itemsPerPage: string;
  setItemsPerPage: (items: string) => void;
  statusFilter: TranslationStatus | 'all';
  setStatusFilter: (status: TranslationStatus | 'all') => void;
  languageFilter: LanguageCodeWithAll;
  setLanguageFilter: (lang: LanguageCodeWithAll) => void;
  issueFilter: IssueStatus;
  setIssueFilter: (issue: IssueStatus) => void;
  prFilter: PrStatus;
  setPrFilter: (pr: PrStatus) => void;
  searchQuery: string;
  setSearchQuery: (query: string) => void;
  debouncedSearchQuery: string;
  setDebouncedSearchQuery: (query: string) => void;
  sortMode: SortMode;
  setSortMode: (mode: SortMode) => void;
  sortDirection: SortDirection;
  setSortDirection: (direction: SortDirection) => void;
  selectedLanguages: LanguageCode[];
  startIndex: number;
  endIndex: number;
}

export const ArticleListControl = ({
  selectedArticleCategory,
  articles,
  filteredArticles,
  activePage,
  setActivePage,
  itemsPerPage,
  setItemsPerPage,
  statusFilter,
  setStatusFilter,
  languageFilter,
  setLanguageFilter,
  issueFilter,
  setIssueFilter,
  prFilter,
  setPrFilter,
  searchQuery,
  setSearchQuery,
  debouncedSearchQuery,
  setDebouncedSearchQuery,
  sortMode,
  setSortMode,
  sortDirection,
  setSortDirection,
  selectedLanguages,
  startIndex,
  endIndex,
}: Props) => {
  const navigate = useNavigate();

  const isMobile = useMediaQuery('(max-width: 768px)');

  const debouncedSearch = useDebouncedCallback((query: string) => {
    setDebouncedSearchQuery(query);
  }, 500);

  const resetFilters = () => {
    debouncedSearch('');
    navigate({ pathname: `/`, search: `?category=${selectedArticleCategory}` });
  };

  const statusOptions = [
    { value: 'all', label: 'All Status' },
    { value: 'up_to_date', label: 'Up to date' },
    { value: 'possibly_outdated', label: 'Possibly outdated' },
    { value: 'outdated', label: 'Outdated' },
    { value: 'not_translated', label: 'Not translated' },
  ];

  const sortedLangCodes = getSortedLangCodes(selectedLanguages);
  const languageOptions = [
    { value: 'all', label: 'All Languages' },
    ...sortedLangCodes.map((code) => ({ value: code.value, label: code.label })),
  ];

  const totalPages = Math.ceil(filteredArticles.length / parseInt(itemsPerPage, 10));

  return (
    <Card withBorder radius="md" p="md">
      <Stack gap="md">
        <Group gap="xs" wrap="wrap" align="end">
          <Select
            label="Language"
            c="dimmed"
            size="sm"
            value={languageFilter}
            onChange={(value) => {
              setLanguageFilter((value || 'all') as LanguageCodeWithAll);
            }}
            data={languageOptions}
            w={180}
          />
          <Select
            label="Translation Status"
            c="dimmed"
            size="sm"
            value={statusFilter}
            onChange={(value) => {
              setStatusFilter((value as TranslationStatus | 'all') || 'all');
            }}
            data={statusOptions}
            w={180}
          />
          <Select
            label="Issue"
            c="dimmed"
            size="sm"
            value={issueFilter}
            onChange={(value) => {
              setIssueFilter((value as IssueStatus) || 'all');
            }}
            data={[
              { value: 'all', label: 'All Status' },
              { value: 'withIssues', label: 'With Issues' },
              { value: 'withoutIssues', label: 'No Issues' },
            ]}
            w={180}
          />
          <Select
            label="Pull Request"
            c="dimmed"
            size="sm"
            value={prFilter}
            onChange={(value) => {
              setPrFilter((value as PrStatus) || 'all');
            }}
            data={[
              { value: 'all', label: 'All Status' },
              { value: 'withPr', label: 'Pull Request' },
              { value: 'withoutPr', label: 'No Pull Request' },
            ]}
            w={180}
          />
          <TextInput
            label="Search"
            c="dimmed"
            value={searchQuery}
            onChange={(e) => {
              const newValue = e.target.value;
              setSearchQuery(newValue);
              debouncedSearch(newValue);
            }}
            placeholder="Article Path"
            style={{ flexGrow: 1 }}
            maw={rem(180)}
          />
          <SortMenu
            sortMode={sortMode}
            setSortMode={setSortMode}
            sortDirection={sortDirection}
            setSortDirection={setSortDirection}
          />
          <ActionIcon
            variant="light"
            color="orange"
            size="lg"
            onClick={resetFilters}
            title="Reset filters"
          >
            <IconRefresh size={17} />
          </ActionIcon>
        </Group>

        {/* Pagination Controls */}
        <Group justify="space-between" wrap="wrap">
          <Text size="sm" c="dimmed">
            Showing {startIndex + 1}-{endIndex} of {filteredArticles.length} articles
            {(statusFilter !== 'all' || languageFilter !== 'all' || debouncedSearchQuery) && (
              <Text span c="blue">
                {' '}
                (filtered from {articles.length} total)
              </Text>
            )}
          </Text>
          <Group gap="md">
            <Group gap="xs" align="center">
              <Select
                size="sm"
                value={itemsPerPage}
                onChange={(value) => {
                  setItemsPerPage(value as string);
                }}
                data={['30', '50', '100']}
                allowDeselect={false}
                w={80}
              />
            </Group>
            <Pagination
              value={activePage}
              onChange={setActivePage}
              total={totalPages}
              size="sm"
              withEdges={!isMobile}
              radius="xs"
            />
          </Group>
        </Group>
      </Stack>
    </Card>
  );
};
