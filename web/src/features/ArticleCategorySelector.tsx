import { useState } from 'react';
import { Box, FloatingIndicator, Tabs } from '@mantine/core';
import { articleCategories, type ArticleCategory } from '@/features/translations';

interface Props {
  articleCategory: ArticleCategory;
  onArticleCategoryChange: (category: ArticleCategory) => void;
}

export const ArticleCategorySelector = ({ articleCategory, onArticleCategoryChange }: Props) => {
  const [rootRef, setRootRef] = useState<HTMLDivElement | null>(null);
  const [controlsRefs, setControlsRefs] = useState<Record<string, HTMLButtonElement | null>>({});

  const setControlRef = (val: string) => (node: HTMLButtonElement) => {
    controlsRefs[val] = node;
    setControlsRefs(controlsRefs);
  };

  const tabListStyle = {
    position: 'relative' as const,
    padding: '4px',
    display: 'flex',
    gap: '4px',
  };

  const tabStyle = {
    border: 'none',
    backgroundColor: 'transparent',
    borderRadius: 'var(--mantine-radius-sm)',
    padding: '8px 16px',
    fontWeight: 800,
    fontSize: '14px',
    color: 'var(--mantine-color-gray-7)',
    cursor: 'pointer',
    transition: 'all 0.2s ease',
    position: 'relative' as const,
    zIndex: 1,
  };

  const indicatorStyle = {
    backgroundColor: 'var(--mantine-color-gray-1)',
    borderRadius: 'var(--mantine-radius-sm)',
    border: '1px solid var(--mantine-color-gray-2)',
    transition: 'all 0.2s ease',
  };

  const handleTabChange = (value: string | null) => {
    onArticleCategoryChange(value as ArticleCategory);
  };

  return (
    <Box mb="md">
      <Tabs value={articleCategory} onChange={handleTabChange} variant="none">
        <Tabs.List ref={setRootRef} style={tabListStyle}>
          {articleCategories.map((category) => (
            <Tabs.Tab
              key={category.value}
              value={category.value}
              ref={setControlRef(category.value)}
              style={tabStyle}
            >
              {category.label}
            </Tabs.Tab>
          ))}
          <FloatingIndicator
            target={articleCategory ? controlsRefs[articleCategory] : null}
            parent={rootRef}
            style={indicatorStyle}
          />
        </Tabs.List>
      </Tabs>
    </Box>
  );
};
