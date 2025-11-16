import { Anchor, Text, Tooltip } from '@mantine/core';
import { type LanguageCode } from '@/features/language/languageCodes';
import { type ArticleTranslation } from '@/features/translations';

interface PRListProps {
  article: ArticleTranslation;
  langCode: LanguageCode;
}

export const PRList = ({ article, langCode }: PRListProps) => {
  const translation = article.translations[langCode];
  if (!translation || !translation.prs || translation.prs.length === 0) {
    return null;
  }

  return (
    <Text size="xs" c="dimmed">
      PR:{' '}
      {translation.prs.map((pr) => (
        <Tooltip key={`pr-${pr.number}`} label={`PR #${pr.number} - ${pr.title}`}>
          <Text size="xs" c="dimmed" component="span">
            <Anchor href={`${pr.url}`} target="_blank" rel="noopener noreferrer">
              #{pr.number}{' '}
            </Anchor>
          </Text>
        </Tooltip>
      ))}
    </Text>
  );
};
