export const languageCodes = [
  { value: 'en', label: 'English' },
  { value: 'bn', label: 'Bengali' },
  { value: 'zh-cn', label: 'Chinese' },
  { value: 'fr', label: 'French' },
  { value: 'de', label: 'German' },
  { value: 'hi', label: 'Hindi' },
  { value: 'id', label: 'Indonesian' },
  { value: 'it', label: 'Italian' },
  { value: 'ja', label: 'Japanese' },
  { value: 'ko', label: 'Korean' },
  { value: 'fa', label: 'Persian' },
  { value: 'pl', label: 'Polish' },
  { value: 'pt-br', label: 'Portuguese' },
  { value: 'ru', label: 'Russian' },
  { value: 'es', label: 'Spanish' },
  { value: 'uk', label: 'Ukrainian' },
  { value: 'vi', label: 'Vietnamese' },
] as const;

export type LanguageCode = (typeof languageCodes)[number]['value'];
export type LanguageCodeWithAll = LanguageCode | 'all';

// List of language map according to RFC 5646.
// See <http://tools.ietf.org/html/rfc5646>
export const browserLanguageMap: Record<string, LanguageCode> = {
  // English
  en: 'en',
  'en-US': 'en',
  'en-GB': 'en',
  'en-AU': 'en',
  'en-CA': 'en',
  'en-IE': 'en',
  'en-NZ': 'en',
  'en-ZA': 'en',

  // Bengali
  bn: 'bn',
  'bn-BD': 'bn',
  'bn-IN': 'bn',

  // Chinese
  zh: 'zh-cn',
  'zh-CN': 'zh-cn',
  'zh-SG': 'zh-cn',
  'zh-TW': 'zh-cn',
  'zh-HK': 'zh-cn',
  'zh-MO': 'zh-cn',

  // French
  fr: 'fr',
  'fr-FR': 'fr',
  'fr-CA': 'fr',
  'fr-BE': 'fr',
  'fr-CH': 'fr',
  'fr-LU': 'fr',
  'fr-MC': 'fr',

  // German
  de: 'de',
  'de-DE': 'de',
  'de-AT': 'de',
  'de-CH': 'de',
  'de-LI': 'de',
  'de-LU': 'de',

  // Hindi
  hi: 'hi',
  'hi-IN': 'hi',

  // Indonesian
  id: 'id',
  'id-ID': 'id',

  // Italian
  it: 'it',
  'it-IT': 'it',
  'it-CH': 'it',

  // Japanese
  ja: 'ja',
  'ja-JP': 'ja',

  // Korean
  ko: 'ko',
  'ko-KR': 'ko',

  // Persian
  fa: 'fa',
  'fa-IR': 'fa',

  // Polish
  pl: 'pl',
  'pl-PL': 'pl',

  // Portuguese
  pt: 'pt-br',
  'pt-BR': 'pt-br',
  'pt-PT': 'pt-br',

  // Russian
  ru: 'ru',
  'ru-RU': 'ru',

  // Spanish
  es: 'es',
  'es-ES': 'es',
  'es-AR': 'es',
  'es-BO': 'es',
  'es-CL': 'es',
  'es-CO': 'es',
  'es-CR': 'es',
  'es-DO': 'es',
  'es-EC': 'es',
  'es-GT': 'es',
  'es-HN': 'es',
  'es-MX': 'es',
  'es-NI': 'es',
  'es-PA': 'es',
  'es-PE': 'es',
  'es-PR': 'es',
  'es-PY': 'es',
  'es-SV': 'es',
  'es-UY': 'es',
  'es-VE': 'es',

  // Ukrainian
  uk: 'uk',
  'uk-UA': 'uk',

  // Vietnamese
  vi: 'vi',
  'vi-VN': 'vi',
};

export const getSortedLangCodes = (
  selectedLanguages?: LanguageCode[]
): { value: LanguageCode; label: string }[] => {
  const head = ['en'];

  if (selectedLanguages && Array.isArray(selectedLanguages)) {
    selectedLanguages.forEach((lang) => {
      if (lang !== 'en' && !head.includes(lang)) {
        head.push(lang);
      }
    });
  }
  const rest = languageCodes.filter((l) => !head.includes(l.value));
  return [...head.map((code) => languageCodes.find((l) => l.value === code)!), ...rest];
};

export const getLanguageName = (langCode: LanguageCode): string => {
  const languageMap = {
    en: 'English',
    bn: 'Bengali',
    'zh-cn': 'Chinese',
    fr: 'French',
    de: 'German',
    hi: 'Hindi',
    id: 'Indonesian',
    it: 'Italian',
    ja: 'Japanese',
    ko: 'Korean',
    fa: 'Persian',
    pl: 'Polish',
    'pt-br': 'Portuguese',
    ru: 'Russian',
    es: 'Spanish',
    uk: 'Ukrainian',
    vi: 'Vietnamese',
  };
  return languageMap[langCode] || langCode;
};
