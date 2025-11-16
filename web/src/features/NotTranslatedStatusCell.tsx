import { rem, Table } from '@mantine/core';
import { GitHubIssueButton } from '@/features/GitHubIssueButton';
import { GitHubPRTemplateGenerator } from '@/features/GitHubPRTemplateGenerator';
import { IssueList } from '@/features/IssueList';
import { type LanguageCode } from '@/features/language/languageCodes';
import { PRList } from '@/features/PRList';
import { StatusBadge } from '@/features/StatusBadge';
import { type ArticleTranslation } from '@/features/translations';

interface NotTranslatedStatusCellProps {
  article: ArticleTranslation;
  langCode: LanguageCode;
}

export const NotTranslatedStatusCell = ({ article, langCode }: NotTranslatedStatusCellProps) => {
  return (
    <Table.Td
      style={{
        textAlign: 'center',
        whiteSpace: 'nowrap',
        minWidth: rem(200),
      }}
    >
      <StatusBadge status="not_translated" />
      <div>
        <GitHubIssueButton
          englishPath={article.englishPath}
          englishUrl={article.englishUrl}
          translationUrl={null}
          langCode={langCode}
          variant="new"
        />
        <GitHubPRTemplateGenerator
          englishPath={article.englishPath}
          englishUrl={article.englishUrl}
          translationUrl={null}
          langCode={langCode}
          isNewTranslation={true}
        />
      </div>
      <IssueList article={article} langCode={langCode} />
      <PRList article={article} langCode={langCode} />
    </Table.Td>
  );
};
