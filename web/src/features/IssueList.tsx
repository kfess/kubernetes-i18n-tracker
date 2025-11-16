import { Anchor, Text, Tooltip } from '@mantine/core';
import { type LanguageCode } from '@/features/language/languageCodes';
import { type ArticleTranslation } from '@/features/translations';

interface IssueListProps {
  article: ArticleTranslation;
  langCode: LanguageCode;
}

export const IssueList = ({ article, langCode }: IssueListProps) => {
  const translation = article.translations[langCode];
  if (!translation || !translation.issues || translation.issues.length === 0) {
    return null;
  }

  return (
    <Text size="xs" c="dimmed">
      Issue:{' '}
      {translation.issues.map((issue) => (
        <Tooltip key={`issue-${issue.number}`} label={`Issue #${issue.number} - ${issue.title}`}>
          <Text size="xs" c="dimmed" component="span">
            <Anchor href={`${issue.url}`} target="_blank" rel="noopener noreferrer">
              #{issue.number}{' '}
            </Anchor>
          </Text>
        </Tooltip>
      ))}
    </Text>
  );
};
