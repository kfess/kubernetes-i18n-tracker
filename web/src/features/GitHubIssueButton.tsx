import { IconCircleDot } from '@tabler/icons-react';
import { ActionIcon, Tooltip } from '@mantine/core';
import { languageCodes, type LanguageCode } from '@/features/language/languageCodes';

export const generateIssueUrl = (
  englishPath: string,
  langCode: LanguageCode,
  isNewTranslation: boolean
): string => {
  const englishUrl = englishPath
    .replace(/^content\/en\//, '/')
    .replace(/\/_?index\.md$/, '/')
    .replace(/\.md$/, '/');

  const translationPath = englishPath.replace('/en/', `/${langCode}/`);
  const translationUrl = englishUrl.replace(/^\//, `/${langCode}/`);

  const languageName =
    languageCodes.find((lang) => lang.value === langCode)?.label || langCode.toUpperCase();
  const languageLabel = `language/${langCode}`;

  if (isNewTranslation) {
    const title = `[${langCode}] Translate ${englishPath} into ${languageName}`;
    const body = `**This is a Feature Request**

**What would you like to be added**

Translate [${englishPath}](https://kubernetes.io${englishUrl}) into ${languageName}

**Why is this needed**

This page is not translated yet.`;

    return (
      `https://github.com/kubernetes/website/issues/new?` +
      `labels=${encodeURIComponent(languageLabel)}&` +
      `title=${encodeURIComponent(title)}&` +
      `body=${encodeURIComponent(body)}`
    );
  } else {
    const title = `[${langCode}] Update ${translationPath}`;
    const body = `**This is a Feature Request**

**What would you like to be added**

Update the ${languageName} translation of the following page to match the latest English version:

- ${languageName}: https://kubernetes.io${translationUrl}
- English: https://kubernetes.io${englishUrl}

**Why is this needed**

The current ${languageName} translation is outdated.`;

    return (
      `https://github.com/kubernetes/website/issues/new?` +
      `labels=${encodeURIComponent(languageLabel)}&` +
      `title=${encodeURIComponent(title)}&` +
      `body=${encodeURIComponent(body)}`
    );
  }
};

interface GitHubIssueButtonProps {
  englishPath: string;
  langCode: LanguageCode;
  variant: 'update' | 'new';
}

export const GitHubIssueButton = ({ englishPath, langCode, variant }: GitHubIssueButtonProps) => {
  const isNewTranslation = variant === 'new';
  const tooltipLabel = isNewTranslation
    ? 'Request new translation Issue'
    : 'Report outdated translation Issue';
  const title = isNewTranslation ? 'Request translation on GitHub' : 'Report issue on GitHub';

  return (
    <Tooltip label={tooltipLabel} position="top" withArrow>
      <ActionIcon
        component="a"
        href={generateIssueUrl(englishPath, langCode, isNewTranslation)}
        target="_blank"
        rel="noopener noreferrer"
        size="xs"
        radius="xs"
        c="blue"
        variant="subtle"
        title={title}
      >
        <IconCircleDot size={14} />
      </ActionIcon>
    </Tooltip>
  );
};
