import { memo, useCallback } from 'react';
import { IconGitPullRequest } from '@tabler/icons-react';
import { ActionIcon, Tooltip } from '@mantine/core';
import { languageCodes, type LanguageCode } from '@/features/language/languageCodes';

const generatePRTemplateMarkdown = (
  englishPath: string,
  englishUrl: string | null,
  translationUrl: string | null,
  langCode: LanguageCode,
  isNewTranslation: boolean,
  issueNumber?: number
): string => {
  const translationPath = englishPath.replace('/en/', `/${langCode}/`);
  const languageName =
    languageCodes.find((lang) => lang.value === langCode)?.label || langCode.toUpperCase();

  let description: string;
  if (isNewTranslation) {
    description = `### Description\n\nTranslated \`${englishPath}\` into ${languageName}: \`${translationPath}\`.\n\n`;
  } else {
    description = `### Description\n\nUpdated ${languageName} translation: \`${translationPath}\`.\n\n`;
  }

  const hasLinks = translationUrl || englishUrl;
  if (hasLinks) {
    description += '**Website Link**:\n\n';
    if (translationUrl) {
      description += `- ${languageName}: ${translationUrl}\n`;
    }
    if (englishUrl) {
      description += `- English: ${englishUrl}\n`;
    }
    description += '\n';
  }

  const issueSection = `\n### Issue\n\nCloses: ${issueNumber ? `#${issueNumber}` : '#'}\n\n/area localization\n/language ${langCode}\n`;

  return description + issueSection;
};

interface Props {
  englishPath: string;
  englishUrl: string | null;
  translationUrl: string | null;
  langCode: LanguageCode;
  isNewTranslation: boolean;
}

export const GitHubPRTemplateGenerator = memo(
  ({ englishPath, englishUrl, translationUrl, langCode, isNewTranslation }: Props) => {
    const handleCopy = useCallback(() => {
      const markdown = generatePRTemplateMarkdown(
        englishPath,
        englishUrl,
        translationUrl,
        langCode,
        isNewTranslation
      );
      navigator.clipboard.writeText(markdown);
    }, [englishPath, englishUrl, translationUrl, langCode, isNewTranslation]);

    return (
      <Tooltip label="Copy PR template markdown to clipboard" position="top" withArrow>
        <ActionIcon component="button" size="xs" radius="xs" variant="subtle" onClick={handleCopy}>
          <IconGitPullRequest size={14} />
        </ActionIcon>
      </Tooltip>
    );
  }
);

GitHubPRTemplateGenerator.displayName = 'GitHubPRTemplateGenerator';
