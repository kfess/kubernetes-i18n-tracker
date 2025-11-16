import { IconExternalLink } from '@tabler/icons-react';
import { ActionIcon, Anchor, Group, rem, Table, Text } from '@mantine/core';
import { GitHubIssueButton } from '@/features/GitHubIssueButton';
import { GitHubPRTemplateGenerator } from '@/features/GitHubPRTemplateGenerator';
import { IssueList } from '@/features/IssueList';
import { type LanguageCode } from '@/features/language/languageCodes';
import { PRList } from '@/features/PRList';
import { StatusBadge } from '@/features/StatusBadge';
import { type ArticleTranslation } from '@/features/translations';
import { formatDateISO } from '@/utils/date';

interface PossiblyOutdatedStatusCellProps {
  article: ArticleTranslation;
  langCode: LanguageCode;
}

export const PossiblyOutdatedStatusCell = ({
  article,
  langCode,
}: PossiblyOutdatedStatusCellProps) => {
  const translation = article.translations[langCode];
  const translationPath = article.englishPath.replace('/en/', `/${langCode}/`);

  return (
    <Table.Td
      style={{
        textAlign: 'center',
        whiteSpace: 'nowrap',
        backgroundColor: 'rgba(34, 139, 34, 0.05)',
        minWidth: rem(200),
      }}
    >
      <Anchor
        href={`https://github.com/kubernetes/website/blob/main/${translationPath}`}
        target="_blank"
        rel="noopener noreferrer"
        underline="never"
        title={`Edit ${langCode} translation on GitHub`}
      >
        <StatusBadge status="possibly_outdated" />
      </Anchor>
      {translation.targetLatestDate && (
        <Group gap="2" justify="center" align="center">
          <Text size="xs" c="dimmed">
            Updated: {formatDateISO(translation.targetLatestDate)} (UTC)
          </Text>
          {translation.translationUrl && (
            <ActionIcon
              component="a"
              href={translation.translationUrl}
              target="_blank"
              rel="noopener noreferrer"
              size="xs"
              radius="xs"
              c="gray"
              variant="subtle"
              title="Kubernetes documentation"
            >
              <IconExternalLink size={14} />
            </ActionIcon>
          )}
          <GitHubIssueButton
            englishPath={article.englishPath}
            englishUrl={article.englishUrl}
            translationUrl={translation.translationUrl || null}
            langCode={langCode}
            variant="update"
          />
          <GitHubPRTemplateGenerator
            englishPath={article.englishPath}
            englishUrl={article.englishUrl}
            translationUrl={translation.translationUrl || null}
            langCode={langCode}
            isNewTranslation={false}
          />
        </Group>
      )}
      <IssueList article={article} langCode={langCode} />
      <PRList article={article} langCode={langCode} />
    </Table.Td>
  );
};
