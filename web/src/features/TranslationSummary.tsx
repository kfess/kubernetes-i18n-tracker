import { useFetchAllTranslationArticles } from '@/features/hooks/useFetchTranslationArticles';
import { type LanguageCode } from '@/features/language/languageCodes';
import {
  articleCategories,
  type ArticleCategory,
  type TranslationStatus,
} from '@/features/translations';
import { Anchor, Box, Group, rem, Table, Text, Tooltip } from '@mantine/core';
import { useMemo } from 'react';
import { Link } from 'react-router-dom';

const statuses: { value: TranslationStatus; label: string; color: string }[] = [
  { value: 'up_to_date', label: 'Up to date', color: 'var(--mantine-color-green-8)' },
  {
    value: 'possibly_outdated',
    label: 'Possibly outdated',
    color: 'var(--mantine-color-yellow-7)',
  },
  { value: 'outdated', label: 'Outdated', color: 'var(--mantine-color-orange-9)' },
  { value: 'not_translated', label: 'Not translated', color: 'var(--mantine-color-gray-4)' },
];

const Swatch = ({ color }: { color: string }) => (
  <Box w={10} h={10} bg={color} style={{ borderRadius: 2, flexShrink: 0 }} />
);

type StatusCounts = Record<TranslationStatus, number>;

interface SummaryRow {
  // undefined for the row that sums up all categories
  category?: ArticleCategory;
  label: string;
  counts: StatusCounts;
  total: number;
}

const emptyCounts = (): StatusCounts => ({
  up_to_date: 0,
  possibly_outdated: 0,
  outdated: 0,
  not_translated: 0,
});

const formatPercent = (count: number, total: number): string => {
  if (total === 0) {
    return '0%';
  }
  const percent = (count / total) * 100;
  // Avoid showing "0%" or "100%" for values that are only close to them
  if (percent > 0 && percent < 1) {
    return '<1%';
  }
  if (percent > 99 && percent < 100) {
    return '>99%';
  }
  return `${Math.round(percent)}%`;
};

const buildListUrl = (category: ArticleCategory, langCode: LanguageCode, status?: string) => {
  const params = new URLSearchParams({ category, lang: langCode });
  if (status) {
    params.set('status', status);
  }
  return `/?${params.toString()}`;
};

interface StatusBarProps {
  row: SummaryRow;
  langCode: LanguageCode;
}

const StatusBar = ({ row, langCode }: StatusBarProps) => {
  return (
    <Group gap={2} wrap="nowrap" h={16} style={{ borderRadius: 4, overflow: 'hidden' }}>
      {statuses.map((status) => {
        const count = row.counts[status.value];
        if (count === 0) {
          return null;
        }

        const label = `${status.label}: ${count.toLocaleString()} (${formatPercent(count, row.total)})`;
        const style = {
          flex: `${count} 1 0`,
          minWidth: 4,
          height: '100%',
          backgroundColor: status.color,
        };

        return (
          <Tooltip key={status.value} label={label} withArrow>
            {row.category ? (
              <Link
                to={buildListUrl(row.category, langCode, status.value)}
                style={style}
                aria-label={label}
              />
            ) : (
              <span style={style} />
            )}
          </Tooltip>
        );
      })}
    </Group>
  );
};

interface Props {
  langCode: LanguageCode;
}

export const TranslationSummary = ({ langCode }: Props) => {
  const allArticles = useFetchAllTranslationArticles();

  const rows = useMemo(() => {
    const totalCounts = emptyCounts();
    const categoryRows: SummaryRow[] = [];

    articleCategories.forEach((category) => {
      const counts = emptyCounts();
      allArticles[category.value].articles.forEach((article) => {
        const status = article.translations[langCode]?.status;
        if (status && status in counts) {
          counts[status]++;
          totalCounts[status]++;
        }
      });

      const total = Object.values(counts).reduce((sum, count) => sum + count, 0);
      if (total > 0) {
        categoryRows.push({ category: category.value, label: category.label, counts, total });
      }
    });

    const total = Object.values(totalCounts).reduce((sum, count) => sum + count, 0);
    return [{ label: 'All categories', counts: totalCounts, total }, ...categoryRows];
  }, [allArticles, langCode]);

  if (rows[0].total === 0) {
    return (
      <Text size="sm" c="dimmed">
        No translation data available for this language.
      </Text>
    );
  }

  return (
    <>
      {/* On small screens the column headers that explain the colors are scrolled out of view */}
      <Group gap="md" mb="xs" hiddenFrom="sm">
        {statuses.map((status) => (
          <Group key={status.value} gap={6} wrap="nowrap">
            <Swatch color={status.color} />
            <Text size="xs">{status.label}</Text>
          </Group>
        ))}
      </Group>
      <Table.ScrollContainer minWidth={rem(860)}>
        <Table verticalSpacing="sm">
          <Table.Thead>
            <Table.Tr>
              <Table.Th>Category</Table.Th>
              <Table.Th w="35%">Share of English pages</Table.Th>
              {statuses.map((status) => (
                <Table.Th key={status.value} ta="right">
                  <Group gap={6} justify="flex-end" wrap="nowrap">
                    <Swatch color={status.color} />
                    {status.label}
                  </Group>
                </Table.Th>
              ))}
              <Table.Th ta="right">Total</Table.Th>
            </Table.Tr>
          </Table.Thead>
          <Table.Tbody>
            {rows.map((row) => (
              <Table.Tr key={row.category ?? 'all'} fw={row.category ? undefined : 700}>
                <Table.Td style={{ whiteSpace: 'nowrap' }}>
                  {row.category ? (
                    <Anchor component={Link} to={buildListUrl(row.category, langCode)} size="sm">
                      {row.label}
                    </Anchor>
                  ) : (
                    row.label
                  )}
                </Table.Td>
                <Table.Td>
                  <StatusBar row={row} langCode={langCode} />
                </Table.Td>
                {statuses.map((status) => {
                  const count = row.counts[status.value];
                  return (
                    <Table.Td
                      key={status.value}
                      ta="right"
                      style={{ whiteSpace: 'nowrap', fontVariantNumeric: 'tabular-nums' }}
                    >
                      {row.category && count > 0 ? (
                        <Anchor
                          component={Link}
                          to={buildListUrl(row.category, langCode, status.value)}
                          size="sm"
                        >
                          {count.toLocaleString()}
                        </Anchor>
                      ) : (
                        <Text span size="sm" fw="inherit" c={count === 0 ? 'dimmed' : undefined}>
                          {count.toLocaleString()}
                        </Text>
                      )}
                      <Text span size="xs" c="dimmed" fw={400} ml={6}>
                        ({formatPercent(count, row.total)})
                      </Text>
                    </Table.Td>
                  );
                })}
                <Table.Td ta="right" style={{ fontVariantNumeric: 'tabular-nums' }}>
                  {row.total.toLocaleString()}
                </Table.Td>
              </Table.Tr>
            ))}
          </Table.Tbody>
        </Table>
      </Table.ScrollContainer>
    </>
  );
};
