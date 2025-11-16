import { IconExternalLink, IconGitBranch } from '@tabler/icons-react';
import { useNavigate } from 'react-router-dom';
import { ActionIcon, Anchor, Group, rem, Table, Text } from '@mantine/core';
import { GitHubIssueButton } from '@/features/GitHubIssueButton';
import { GitHubPRTemplateGenerator } from '@/features/GitHubPRTemplateGenerator';
import { IssueList } from '@/features/IssueList';
import { type LanguageCode } from '@/features/language/languageCodes';
import { PRList } from '@/features/PRList';
import { StatusBadge } from '@/features/StatusBadge';
import { ArticleCategory, type ArticleTranslation } from '@/features/translations';
import { formatDateISO } from '@/utils/date';

interface OutdatedStatusCellProps {
  article: ArticleTranslation;
  langCode: LanguageCode;
  category: ArticleCategory;
}

export const OutdatedStatusCell = ({ article, langCode, category }: OutdatedStatusCellProps) => {
  const navigate = useNavigate();
  const translation = article.translations[langCode];
  const translationPath = article.englishPath.replace('/en/', `/${langCode}/`);

  const handleDiffClick = () => {
    const params = new URLSearchParams({
      category,
      translationPath,
      language: langCode,
    });
    navigate(`/detail?${params.toString()}`);
  };

  return (
    <Table.Td
      style={{
        textAlign: 'center',
        whiteSpace: 'nowrap',
        backgroundColor: 'rgba(255, 165, 0, 0.1)',
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
        <StatusBadge status="outdated" />
      </Anchor>
      {translation.totalChangeLines > 0 && (
        <Text size="sm">{translation.totalChangeLines.toLocaleString()} lines changed</Text>
      )}
      {translation.daysBehind && (
        <Text size="xs" c="dimmed">
          {translation.commitsBehind.toLocaleString()} commit
          {translation.commitsBehind > 1 ? 's' : ''} / {translation.daysBehind.toLocaleString()} day
          {translation.daysBehind !== 1 ? 's' : ''} behind
        </Text>
      )}
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
          <ActionIcon
            onClick={handleDiffClick}
            size="xs"
            radius="xs"
            c="blue"
            variant="subtle"
            title="View translation diff"
            style={{ cursor: 'pointer' }}
          >
            <IconGitBranch size={14} />
          </ActionIcon>
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
