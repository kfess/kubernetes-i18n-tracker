import { memo } from 'react';
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

/area localization
/language ${langCode}
/assign`;

    return (
      `https://github.com/kubernetes/website/issues/new?` +
      `template=feature-request.md&` +
      `title=${encodeURIComponent(title)}&` +
      `body=${encodeURIComponent(body)}`
    );
  } else {
    const title = `[${langCode}] Update ${translationPath}`;

    let whatToAdd = `Update the ${languageName} translation of \`${translationPath}\` to match the latest English version.`;

    const hasLinks = englishUrl || translationUrl;
    if (hasLinks) {
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

/area localization
/language ${langCode}
/assign`;

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

export const GitHubIssueButton = memo(
  ({ englishPath, englishUrl, translationUrl, langCode, variant }: GitHubIssueButtonProps) => {
    const isNewTranslation = variant === 'new';
    const tooltipLabel = isNewTranslation
      ? 'Request new translation Issue on GitHub'
      : 'Report outdated translation Issue on GitHub';

    return (
      <Tooltip label={tooltipLabel} position="top" withArrow>
        <ActionIcon
          component="a"
          href={generateIssueUrl(
            englishPath,
            englishUrl,
            translationUrl,
            langCode,
            isNewTranslation
          )}
          target="_blank"
          rel="noopener noreferrer"
          size="xs"
          radius="xs"
          c="blue"
          variant="subtle"
        >
          <IconCircleDot size={14} />
        </ActionIcon>
      </Tooltip>
    );
  }
);

GitHubIssueButton.displayName = 'GitHubIssueButton';
