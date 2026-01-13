import { IconInfoCircle } from '@tabler/icons-react';
import { Alert, Anchor, List } from '@mantine/core';
import { getLanguageName, LanguageCode } from '@/features/language/languageCodes';

interface Props {
  selectedLanguages: LanguageCode[];
}

const LOCALIZATION_POLICY_URLS: Record<
  Extract<LanguageCode, 'zh-cn' | 'fr' | 'de' | 'ja' | 'ko' | 'pl' | 'ru'>,
  string
> = {
  'zh-cn': 'https://kubernetes.io/zh-cn/docs/contribute/localization/',
  fr: 'https://kubernetes.io/fr/docs/contribute/localization/',
  de: 'https://kubernetes.io/de/docs/contribute/localization/',
  ja: 'https://kubernetes.io/ja/docs/contribute/localization/',
  ko: 'https://kubernetes.io/ko/docs/contribute/localization/',
  pl: 'https://kubernetes.io/pl/docs/contribute/localization/',
  ru: 'https://kubernetes.io/ru/docs/contribute/localization/',
};

export const LocalizationPolicyWarning = ({ selectedLanguages }: Props) => {
  const matchedLanguages = (selectedLanguages || []).filter(
    (l) => l in LOCALIZATION_POLICY_URLS
  ) as Array<keyof typeof LOCALIZATION_POLICY_URLS>;

  if (matchedLanguages.length > 0) {
    return (
      <Alert
        withCloseButton
        variant="light"
        color="orange"
        title="Localization Policy"
        icon={<IconInfoCircle size={16} />}
        mb={20}
      >
        <div>
          Some selected languages have specific localization policies. Please review them thoroughly
          before submitting PR:
          <List size="xs" py="4">
            {matchedLanguages.map((lang) => (
              <List.Item key={lang}>
                <Anchor
                  href={LOCALIZATION_POLICY_URLS[lang]}
                  size="sm"
                  target="_blank"
                  rel="noopener"
                >
                  {getLanguageName(lang)} Localization Policy
                </Anchor>
              </List.Item>
            ))}
          </List>
        </div>
      </Alert>
    );
  }
};
