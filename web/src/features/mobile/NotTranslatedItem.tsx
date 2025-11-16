import { Anchor, Box, Group, Text, Tooltip } from '@mantine/core';
import { GitHubIssueButton } from '@/features/GitHubIssueButton';
import { GitHubPRTemplateGenerator } from '@/features/GitHubPRTemplateGenerator';
import { getLanguageName, type LanguageCode } from '@/features/language/languageCodes';
import { type ArticleTranslation } from '@/features/translations';

type TranslationInfo = ArticleTranslation['translations'][LanguageCode];

interface Props {
  code: LanguageCode;
  translation: TranslationInfo | undefined;
  article: ArticleTranslation;
}

export const NotTranslatedItem = ({ code, translation, article }: Props) => {
  return (
    <Box
      px="xs"
      py={4}
      bg="white"
      style={{
        borderRadius: 16,
        fontSize: 11,
        color: '#64748b',
        border: '1px solid #e2e8f0',
        fontWeight: 500,
      }}
    >
      <Group gap={4} align="center">
        {getLanguageName(code)}
        <GitHubIssueButton
          englishPath={article.englishPath}
          englishUrl={article.englishUrl}
          translationUrl={null}
          langCode={code}
          variant="new"
        />
        <GitHubPRTemplateGenerator
          englishPath={article.englishPath}
          englishUrl={article.englishUrl}
          translationUrl={null}
          langCode={code}
          isNewTranslation={true}
        />
      </Group>
      {translation?.issues && translation.issues.length > 0 && (
        <Text size="xs" c="dimmed" mt="xs" component="span">
          {' '}
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
      {translation?.prs && translation?.prs.length > 0 && (
        <Text size="xs" c="dimmed" component="span">
          {translation?.prs.map((pr) => (
            <Tooltip key={pr.number} label={`PR #${pr.number} - ${pr.title}`}>
              <Text size="xs" c="dimmed" component="span">
                <Anchor href={`${pr.url}`} target="_blank" rel="noopener noreferrer">
                  #{pr.number}{' '}
                </Anchor>
              </Text>
            </Tooltip>
          ))}
        </Text>
      )}
    </Box>
  );
};
