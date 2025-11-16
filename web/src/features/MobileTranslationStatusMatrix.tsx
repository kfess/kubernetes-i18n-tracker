import { Card, Stack, Text } from '@mantine/core';
import { type LanguageCode } from '@/features/language/languageCodes';
import { MobileArticleCard } from '@/features/mobile/MobileArticleCard';
import { type ArticleCategory, type ArticleTranslation } from '@/features/translations';

interface Props {
  articles: ArticleTranslation[];
  selectedLanguages: LanguageCode[];
  selectedArticleCategory: ArticleCategory;
}

export const MobileTranslationStatusMatrix = ({
  articles,
  selectedLanguages,
  selectedArticleCategory,
}: Props) => {
  if (articles.length === 0) {
    return (
      <Stack gap="md">
        <Card withBorder radius="md" p="lg" shadow="xs">
          <Text c="dimmed" ta="center">
            No articles found
          </Text>
        </Card>
      </Stack>
    );
  }

  return (
    <Stack gap="md">
      {articles.map((article) => (
        <MobileArticleCard
          key={article.englishPath}
          article={article}
          selectedLanguages={selectedLanguages}
          selectedArticleCategory={selectedArticleCategory}
        />
      ))}
    </Stack>
  );
};
