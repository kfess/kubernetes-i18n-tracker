import { IconCircleDot } from '@tabler/icons-react';
import { ActionIcon, Tooltip } from '@mantine/core';
import { languageCodes, type LanguageCode } from '@/features/language/languageCodes';

export const generateIssueUrl = (
  englishPath: string,
  englishUrl: string | null,
  translationUrl: string | null,
  langCode: LanguageCode,
  isNewTranslation: boolean
): string => {
  const translationPath = englishPath.replace('/en/', `/${langCode}/`);
  const languageName =
    languageCodes.find((lang) => lang.value === langCode)?.label || langCode.toUpperCase();

  if (isNewTranslation) {
    const title = `[${langCode}] Translate ${englishPath} into ${languageName}`;

    let whatToAdd = `Translate \`${englishPath}\` into ${languageName}`;
    if (englishUrl) {
      whatToAdd += `\n\n**Website Link**\n\n- English: ${englishUrl}`;
    }

    const body = `**This is a Feature Request**

**What would you like to be added**

${whatToAdd}

**Why is this needed**

This page is not translated yet.

/language ${langCode}`;

    return (
      `https://github.com/kubernetes/website/issues/new?` +
      `template=feature-request.md&` +
      `title=${encodeURIComponent(title)}&` +
      `body=${encodeURIComponent(body)}`
    );
  } else {
    const title = `[${langCode}] Update ${translationPath}`;

    let whatToAdd = `Update the ${languageName} translation of \`${translationPath}\` to match the latest English version.`;

    if (englishUrl || translationUrl) {
      whatToAdd += '\n\n**Website Link**\n';
      if (translationUrl) {
        whatToAdd += `\n- ${languageName}: ${translationUrl}`;
      }
      if (englishUrl) {
        whatToAdd += `\n- English: ${englishUrl}`;
      }
    }

    const body = `**This is a Feature Request**

**What would you like to be added**

${whatToAdd}

**Why is this needed**

The current ${languageName} translation is outdated.

/language ${langCode}`;

    return (
      `https://github.com/kubernetes/website/issues/new?` +
      `template=feature-request.md&` +
      `title=${encodeURIComponent(title)}&` +
      `body=${encodeURIComponent(body)}`
    );
  }
};

interface GitHubIssueButtonProps {
  englishPath: string;
  englishUrl: string | null;
  translationUrl: string | null;
  langCode: LanguageCode;
  variant: 'update' | 'new';
}

export const GitHubIssueButton = ({
  englishPath,
  englishUrl,
  translationUrl,
  langCode,
  variant,
}: GitHubIssueButtonProps) => {
  const isNewTranslation = variant === 'new';
  const tooltipLabel = isNewTranslation
    ? 'Request new translation Issue'
    : 'Report outdated translation Issue';
  const title = isNewTranslation ? 'Request translation on GitHub' : 'Report issue on GitHub';

  return (
    <Tooltip label={tooltipLabel} position="top" withArrow>
      <ActionIcon
        component="a"
        href={generateIssueUrl(englishPath, englishUrl, translationUrl, langCode, isNewTranslation)}
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
