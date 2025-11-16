import { Box, Card, Divider, Group, Stack, Text } from '@mantine/core';
import { EnglishSourceInfo } from '@/features/EnglishSourceInfo';
import { languageCodes, type LanguageCode } from '@/features/language/languageCodes';
import { NotTranslatedItem } from '@/features/mobile/NotTranslatedItem';
import { OutdatedTranslationItem } from '@/features/mobile/OutdatedTranslationItem';
import { UpToDateTranslationItem } from '@/features/mobile/UpToDateTranslationItem';
import { type ArticleCategory, type ArticleTranslation } from '@/features/translations';

interface Props {
  article: ArticleTranslation;
  selectedLanguages: LanguageCode[];
  selectedArticleCategory: ArticleCategory;
}

export const MobileArticleCard = ({
  article,
  selectedLanguages,
  selectedArticleCategory,
}: Props) => {
  const sortedLangCodes = languageCodes.filter((code) => selectedLanguages.includes(code.value));

  const upToDateOrPossiblyOutdatedLangs = sortedLangCodes.filter(
    (code) =>
      code.value !== 'en' &&
      (article.translations[code.value]?.status === 'up_to_date' ||
        article.translations[code.value]?.status === 'possibly_outdated')
  );

  const outdatedLangs = sortedLangCodes.filter(
    (code) => code.value !== 'en' && article.translations[code.value]?.status === 'outdated'
  );

  const notTranslatedLangs = sortedLangCodes.filter(
    (code) =>
      code.value !== 'en' &&
      (!article.translations[code.value] ||
        article.translations[code.value]?.status === 'not_translated')
  );

  return (
    <Card withBorder radius="md" p="sm" shadow="sm">
      <EnglishSourceInfo article={article} />
      <Divider my="md" />
      {sortedLangCodes.length === 1 ? (
        <Text c="dimmed" size="sm">
          Select target languages to view translation status
        </Text>
      ) : (
        <Stack gap="md">
          {/* Up to date languages */}
          {upToDateOrPossiblyOutdatedLangs.length > 0 && (
            <Stack gap="xs">
              {upToDateOrPossiblyOutdatedLangs.map((code) => {
                const translation = article.translations[code.value];
                return (
                  <UpToDateTranslationItem
                    key={code.value}
                    code={code.value}
                    translation={translation}
                    article={article}
                  />
                );
              })}
            </Stack>
          )}

          {/* Outdated languages */}
          {outdatedLangs.length > 0 && (
            <Stack gap="xs">
              {outdatedLangs.map((code) => {
                const translation = article.translations[code.value];
                return (
                  <OutdatedTranslationItem
                    key={code.value}
                    code={code.value}
                    translation={translation}
                    article={article}
                    selectedArticleCategory={selectedArticleCategory}
                  />
                );
              })}
            </Stack>
          )}

          {/* Not translated languages */}
          {notTranslatedLangs.length > 0 && (
            <Box
              p="xs"
              bg="#f8fafc"
              style={{
                borderRadius: 4,
                borderLeft: '4px solid #cbd5e1',
              }}
            >
              <Group gap="xs" mb="xs">
                <Text size="xs" fw={600} c="dimmed">
                  — Not translated
                </Text>
              </Group>
              <Group gap="xs">
                {notTranslatedLangs.map((code) => {
                  const translation = article.translations[code.value];
                  return (
                    <NotTranslatedItem
                      key={code.value}
                      code={code.value}
                      translation={translation}
                      article={article}
                    />
                  );
                })}
              </Group>
            </Box>
          )}
        </Stack>
      )}
    </Card>
  );
};
