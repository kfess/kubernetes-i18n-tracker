import { IconInfoCircle } from '@tabler/icons-react';
import { Alert, Anchor } from '@mantine/core';
import { getLanguageName, LanguageCode } from '@/features/language/languageCodes';

interface Props {
  selectedLanguages: LanguageCode[];
}

const CASE_STUDIES_REMOVAL_PR_URL = 'https://github.com/kubernetes/website/pull/57043';

// Languages that still have `content/<lang>/case-studies/` on main (as of 2026-10-12).
// en, ja, fa, pt-br, uk have already removed it.
const LANGUAGES_WITH_CASE_STUDIES: ReadonlySet<LanguageCode> = new Set<LanguageCode>([
  'bn',
  'de',
  'es',
  'fr',
  'hi',
  'id',
  'it',
  'ko',
  'pl',
  'ru',
  'vi',
  'zh-cn',
]);

export const TemporaryNotification = ({ selectedLanguages }: Props) => {
  const matchedLanguages = (selectedLanguages || []).filter((l) =>
    LANGUAGES_WITH_CASE_STUDIES.has(l)
  );

  if (matchedLanguages.length === 0) {
    return null;
  }

  const languageNames = matchedLanguages.map((l) => getLanguageName(l)).join(', ');

  return (
    <Alert
      withCloseButton
      variant="light"
      color="orange"
      title="Case Studies Removed"
      icon={<IconInfoCircle size={16} />}
      mb={20}
    >
      The case studies section has been removed from the English documentation (
      <Anchor href={CASE_STUDIES_REMOVAL_PR_URL} size="sm" target="_blank" rel="noopener">
        kubernetes/website#57043
      </Anchor>
      ), so it is no longer tracked here. The following selected languages still have localized
      case studies: {languageNames}. Please consider removing them as well.
    </Alert>
  );
};
