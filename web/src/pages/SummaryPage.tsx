import { Card, Container, Group, Select, Stack, Text, Title } from '@mantine/core';
import { useLocalStorage } from '@mantine/hooks';
import { getSortedLangCodes, type LanguageCode } from '@/features/language/languageCodes';
import { TranslationSummary } from '@/features/TranslationSummary';
import { useQueryParams } from '@/hooks/useQueryParams';
import { getDeploymentInfo } from '@/utils/deploy';

export function SummaryPage() {
  const { deployedAt, gitCommit } = getDeploymentInfo();

  // Selected languages from localStorage
  const [selectedLanguages] = useLocalStorage<LanguageCode[]>({
    key: 'selected-languages',
  });
  const [langParam, setLangParam] = useQueryParams<string | null>('lang', null);

  // Selected languages come first, so the default is the user's own language
  const languageOptions = getSortedLangCodes(selectedLanguages).filter(
    (lang) => lang.value !== 'en'
  );
  const langCode =
    languageOptions.find((lang) => lang.value === langParam)?.value ?? languageOptions[0].value;

  return (
    <Container fluid px={{ base: '0', sm: 'md' }} mt={{ base: 'xs', sm: 'md' }}>
      <Stack gap="sm">
        <Group justify="space-between" align="end" wrap="wrap">
          <div>
            <Title order={2} size="h3">
              Translation Summary
            </Title>
            <Text size="sm" c="dimmed">
              Translation status of every English page, by category. Click a bar or a number to see
              the pages.
            </Text>
          </div>
          <Select
            label="Language"
            c="dimmed"
            size="sm"
            value={langCode}
            onChange={(value) => value && setLangParam(value)}
            data={languageOptions}
            allowDeselect={false}
            w={180}
          />
        </Group>
        <Card withBorder radius="md" p="md">
          <TranslationSummary langCode={langCode} />
        </Card>
        {deployedAt && (
          <Text size="xs" c="dimmed" ta="right">
            Last Updated at {deployedAt} ({gitCommit})
          </Text>
        )}
      </Stack>
    </Container>
  );
}
