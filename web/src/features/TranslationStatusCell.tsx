import { useInView } from 'react-intersection-observer';
import { rem, Table } from '@mantine/core';
import { type LanguageCode } from '@/features/language/languageCodes';
import { ArticleCategory, type ArticleTranslation } from '@/features/translations';
import { NotTranslatedStatusCell } from './NotTranslatedStatusCell';
import { OutdatedStatusCell } from './OutdatedStatusCell';
import { PossiblyOutdatedStatusCell } from './PossiblyOutdatedStatusCell';
import { UpToDateStatusCell } from './UpToDateStatusCell';

export const TranslationStatusCell = ({
  article,
  langCode,
  category,
}: {
  article: ArticleTranslation;
  langCode: LanguageCode;
  category: ArticleCategory;
}) => {
  const { ref, inView } = useInView({
    triggerOnce: true,
    rootMargin: '100px', // load before it comes into view
  });

  const status = article.translations[langCode]?.status;

  if (!inView) {
    return (
      <Table.Td
        style={{
          textAlign: 'center',
          whiteSpace: 'nowrap',
          minWidth: rem(200),
        }}
        ref={ref}
      >
        <div style={{ minHeight: '100px' }} />
      </Table.Td>
    );
  }

  if (status === 'not_translated') {
    return <NotTranslatedStatusCell article={article} langCode={langCode} />;
  }

  if (status === 'up_to_date') {
    return <UpToDateStatusCell article={article} langCode={langCode} />;
  }

  if (status === 'possibly_outdated') {
    return <PossiblyOutdatedStatusCell article={article} langCode={langCode} />;
  }

  if (status === 'outdated') {
    return <OutdatedStatusCell article={article} langCode={langCode} category={category} />;
  }

  return null;
};
