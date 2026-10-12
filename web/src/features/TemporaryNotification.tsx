import { IconInfoCircle } from '@tabler/icons-react';
import { Alert, Anchor } from '@mantine/core';
import { LanguageCode } from '@/features/language/languageCodes';

interface Props {
  selectedLanguages: LanguageCode[];
}

const CASE_STUDIES_REMOVAL_PR_URL = 'https://github.com/kubernetes/website/pull/57043';

// Languages whose `content/<lang>/case-studies/` had not yet been removed on main
// when this notice was added (2026-10-12). Used only to decide whether to show the notice.
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
      ), so it is no longer tracked here. Please consider removing the localized versions as well.
    </Alert>
  );
};
