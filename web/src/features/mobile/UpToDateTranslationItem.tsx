import { IconExternalLink } from '@tabler/icons-react';
import { ActionIcon, Anchor, Box, Group, Text, Tooltip } from '@mantine/core';
import { GitHubIssueButton } from '@/features/GitHubIssueButton';
import { GitHubPRTemplateGenerator } from '@/features/GitHubPRTemplateGenerator';
import { getLanguageName, type LanguageCode } from '@/features/language/languageCodes';
import { StatusBadge } from '@/features/StatusBadge';
import { type ArticleTranslation } from '@/features/translations';
import { formatDateISO } from '@/utils/date';

type TranslationInfo = ArticleTranslation['translations'][LanguageCode];

interface Props {
  code: LanguageCode;
  translation: TranslationInfo;
  article: ArticleTranslation;
}

export const UpToDateTranslationItem = ({ code, translation, article }: Props) => {
  const translationPath = article.englishPath.replace('/en/', `/${code}/`);

  return (
    <Box
      p="xs"
      bg="rgba(34, 139, 34, 0.05)"
      style={{
        borderLeft: '4px solid #10b981',
        borderRadius: 4,
      }}
    >
      <Group justify="space-between" align="center">
        <Group gap="xs" align="center">
          <StatusBadge status={translation.status} />
          <Anchor
            href={`https://github.com/kubernetes/website/blob/main/${translationPath}`}
            target="_blank"
            rel="noopener noreferrer"
            underline="hover"
            c="inherit"
          >
            <Text fw={600} c="dark" size="sm">
              {getLanguageName(code)}
            </Text>
          </Anchor>
          <Text size="xs" c="dimmed">
            Updated at{' '}
            {translation?.targetLatestDate ? formatDateISO(translation.targetLatestDate) : ''} (UTC)
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
            {translation.status === 'possibly_outdated' && (
              <>
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
      {translation.prs.length > 0 && (
        <Text size="xs" c="dimmed" mt="xs">
          PR:{' '}
          {translation?.prs.map((pr) => (
            <Tooltip label={`PR #${pr.number} - ${pr.title}`} key={pr.number}>
              <Text size="xs" c="dimmed" component="span">
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
    </Box>
  );
};
