import { IconGitPullRequest } from '@tabler/icons-react';
import { ActionIcon, Tooltip } from '@mantine/core';
import { useClipboard } from '@mantine/hooks';
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

  // Description section
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

  // Issue section
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

export const GitHubPRTemplateGenerator = ({
  englishPath,
  englishUrl,
  translationUrl,
  langCode,
  isNewTranslation,
}: Props) => {
  const clipboard = useClipboard({ timeout: 500 });

  return (
    <Tooltip label="Copy PR template markdown to clipboard" position="top" withArrow>
      <ActionIcon
        component="button"
        size="xs"
        radius="xs"
        c={clipboard.copied ? 'teal' : 'blue'}
        variant="subtle"
        onClick={() => {
          const markdown = generatePRTemplateMarkdown(
            englishPath,
            englishUrl,
            translationUrl,
            langCode,
            isNewTranslation
          );
          clipboard.copy(markdown);
        }}
      >
        <IconGitPullRequest size={14} />
      </ActionIcon>
    </Tooltip>
  );
};
