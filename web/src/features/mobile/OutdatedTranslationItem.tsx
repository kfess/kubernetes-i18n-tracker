import { IconExternalLink, IconGitBranch } from '@tabler/icons-react';
import { useNavigate } from 'react-router-dom';
import { ActionIcon, Anchor, Box, Group, Stack, Text, Tooltip } from '@mantine/core';
import { GitHubIssueButton } from '@/features/GitHubIssueButton';
import { GitHubPRTemplateGenerator } from '@/features/GitHubPRTemplateGenerator';
import { getLanguageName, type LanguageCode } from '@/features/language/languageCodes';
import { StatusBadge } from '@/features/StatusBadge';
import { type ArticleCategory, type ArticleTranslation } from '@/features/translations';
import { formatDateISO } from '@/utils/date';

type TranslationInfo = ArticleTranslation['translations'][LanguageCode];

interface Props {
  code: LanguageCode;
  translation: TranslationInfo;
  article: ArticleTranslation;
  selectedArticleCategory: ArticleCategory;
}

export const OutdatedTranslationItem = ({
  code,
  translation,
  article,
  selectedArticleCategory,
}: Props) => {
  const navigate = useNavigate();
  const translationPath = article.englishPath.replace('/en/', `/${code}/`);

  const handleDiffClick = () => {
    const params = new URLSearchParams({
      category: selectedArticleCategory,
      translationPath,
      language: code,
    });
    navigate(`/detail?${params.toString()}`);
  };

  return (
    <Box
      p="xs"
      bg="rgba(255, 165, 0, 0.1)"
      style={{
        borderLeft: '4px solid #f59e0b',
        borderRadius: 4,
      }}
    >
      <Stack gap={4}>
        <Group justify="space-between" align="center">
          <Group gap="xs">
            <StatusBadge status="outdated" />
            <Anchor
              href={`https://github.com/kubernetes/website/blob/main/${translationPath}`}
              target="_blank"
              rel="noopener noreferrer"
              underline="hover"
              c="inherit"
            >
              <Text fw={600} c="dark" size="sm" pb={2}>
                {getLanguageName(code)}
              </Text>
            </Anchor>
            <Text size="xs" c="dimmed">
              Updated at{' '}
              {translation?.targetLatestDate ? formatDateISO(translation.targetLatestDate) : ''}{' '}
              (UTC)
            </Text>
            <Group gap="0">
              {translation?.translationUrl && (
                <ActionIcon
                  component="a"
                  href={translation.translationUrl || ''}
                  target="_blank"
                  rel="noopener noreferrer"
                  size="xs"
                  radius="xs"
                  c="gray"
                  variant="subtle"
                  title="View on Kubernetes site"
                >
                  <IconExternalLink size={14} />
                </ActionIcon>
              )}
              {translation.status === 'outdated' && (
                <>
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
                    translationUrl={translation.translationUrl}
                    langCode={code}
                    variant="update"
                  />
                  <GitHubPRTemplateGenerator
                    englishPath={article.englishPath}
                    englishUrl={article.englishUrl}
                    translationUrl={translation.translationUrl}
                    langCode={code}
                    isNewTranslation={false}
                  />
                </>
              )}
            </Group>
          </Group>
        </Group>
        <Text size="xs" c="dimmed">
          {[
            translation?.totalChangeLines && `${translation.totalChangeLines} lines changed`,
            translation?.commitsBehind &&
              translation?.daysBehind &&
              `${translation.commitsBehind} commits`,
            translation?.daysBehind && `${translation.daysBehind} days behind`,
          ]
            .filter(Boolean)
            .join(' • ')}
        </Text>
        {translation.prs.length > 0 && (
          <Text size="xs" c="dimmed">
            PR:{' '}
            {translation?.prs.map((pr) => (
              <Tooltip key={pr.number} label={`PR #${pr.number} - ${pr.title}`}>
                <Text key={pr.number} size="xs" c="dimmed" component="span">
                  <Anchor href={`${pr.url}`} target="_blank" rel="noopener noreferrer">
                    #{pr.number}{' '}
                  </Anchor>
                </Text>
              </Tooltip>
            ))}
          </Text>
        )}
        {translation.issues.length > 0 && (
          <Text size="xs" c="dimmed">
            Issue:{' '}
            {translation.issues.map((issue) => (
              <Tooltip label={`Issue #${issue.number} - ${issue.title}`} key={issue.number}>
                <Text size="xs" c="dimmed" component="span">
                  <Anchor href={`${issue.url}`} target="_blank" rel="noopener noreferrer">
                    #{issue.number}{' '}
                  </Anchor>
                </Text>
              </Tooltip>
            ))}
          </Text>
        )}
      </Stack>
    </Box>
  );
};
