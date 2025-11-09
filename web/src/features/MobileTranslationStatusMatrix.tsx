import { IconExternalLink, IconGitBranch } from '@tabler/icons-react';
import { useNavigate } from 'react-router-dom';
import { ActionIcon, Anchor, Box, Card, Divider, Group, Stack, Text, Tooltip } from '@mantine/core';
import { EnglishSourceInfo } from '@/features/EnglishSourceInfo';
import { GitHubIssueButton } from '@/features/GitHubIssueButton';
import {
  getLanguageName,
  languageCodes,
  type LanguageCode,
} from '@/features/language/languageCodes';
import { StatusBadge } from '@/features/StatusBadge';
import { type ArticleCategory, type ArticleTranslation } from '@/features/translations';
import { formatDateISO } from '@/utils/date';

interface Props {
  articles: ArticleTranslation[];
  selectedLanguages: LanguageCode[];
  selectedArticleCategory: ArticleCategory;
}

export const MobileTranslationStatusMatrix = ({
  articles,
  selectedLanguages,
  selectedArticleCategory,
}: Props) => {
  const navigate = useNavigate();

  const sortedLangCodes = languageCodes.filter((code) => selectedLanguages.includes(code.value));

  if (articles.length === 0) {
    return (
      <Stack gap="md">
        <Card withBorder radius="md" p="lg" shadow="xs">
          <Text c="dimmed" ta="center">
            No articles found
          </Text>
        </Card>
      </Stack>
    );
  }

  return (
    <Stack gap="md">
      {articles.map((article) => {
        const upToDateOrPossiblyOutdatedLangs = sortedLangCodes.filter(
          (code) =>
            code.value !== 'en' &&
            (article.translations[code.value]?.status === 'up_to_date' ||
              article.translations[code.value]?.status === 'possibly_outdated')
        );
        const outdatedLangs = sortedLangCodes.filter(
          (code) => code.value !== 'en' && article.translations[code.value]?.status === 'outdated'
        );
        const notTranslatedLangs = sortedLangCodes.filter(
          (code) =>
            code.value !== 'en' &&
            (!article.translations[code.value] ||
              article.translations[code.value]?.status === 'not_translated')
        );

        return (
          <Card key={article.englishPath} withBorder radius="md" p="sm" shadow="sm">
            <EnglishSourceInfo article={article} />
            <Divider my="md" />
            {sortedLangCodes.length === 1 ? (
              <Text c="dimmed" size="sm">
                Select target languages to view translation status
              </Text>
            ) : (
              <Stack gap="md">
                {/* Up to date languages */}
                {upToDateOrPossiblyOutdatedLangs.length > 0 && (
                  <Stack gap="xs">
                    {upToDateOrPossiblyOutdatedLangs.map((code) => {
                      const translation = article.translations[code.value];
                      const translationPath = article.englishPath.replace(
                        '/en/',
                        `/${code.value}/`
                      );

                      return (
                        <Box
                          key={code.value}
                          p="xs"
                          bg="rgba(34, 139, 34, 0.05)"
                          style={{
                            borderLeft: '4px solid #10b981',
                            borderRadius: 4,
                          }}
                        >
                          <Group justify="space-between" align="center">
                            <Group gap="xs" align="center">
                              <StatusBadge status={translation.status} />
                              <Anchor
                                href={`https://github.com/kubernetes/website/blob/main/${translationPath}`}
                                target="_blank"
                                rel="noopener noreferrer"
                                underline="hover"
                                c="inherit"
                              >
                                <Text fw={600} c="dark" size="sm">
                                  {getLanguageName(code.value)}
                                </Text>
                              </Anchor>
                              <Text size="xs" c="dimmed">
                                Updated at{' '}
                                {translation?.targetLatestDate
                                  ? formatDateISO(translation.targetLatestDate)
                                  : ''}{' '}
                                (UTC)
                              </Text>
                              <Group gap="0">
                                {translation?.translationUrl && (
                                  <ActionIcon
                                    component="a"
                                    href={translation.translationUrl || ''}
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    size="xs"
                                    radius="xs"
                                    c="gray"
                                    variant="subtle"
                                    title="View on Kubernetes site"
                                  >
                                    <IconExternalLink size={14} />
                                  </ActionIcon>
                                )}
                                {translation.status === 'possibly_outdated' && (
                                  <GitHubIssueButton
                                    englishPath={article.englishPath}
                                    langCode={code.value}
                                    variant="update"
                                  />
                                )}
                              </Group>
                            </Group>
                          </Group>
                          {translation.prs.length > 0 && (
                            <Text size="xs" c="dimmed" mt="xs">
                              PR:{' '}
                              {translation?.prs.map((pr) => (
                                <Tooltip label={`PR #${pr.number} - ${pr.title}`} key={pr.number}>
                                  <Text size="xs" c="dimmed" component="span">
                                    <Anchor
                                      href={`${pr.url}`}
                                      target="_blank"
                                      rel="noopener noreferrer"
                                    >
                                      #{pr.number}{' '}
                                    </Anchor>
                                  </Text>
                                </Tooltip>
                              ))}
                            </Text>
                          )}
                          {translation.issues.length > 0 && (
                            <Text size="xs" c="dimmed">
                              Issue:{' '}
                              {translation.issues.map((issue) => (
                                <Tooltip
                                  label={`Issue #${issue.number} - ${issue.title}`}
                                  key={issue.number}
                                >
                                  <Text size="xs" c="dimmed" component="span">
                                    <Anchor
                                      href={`${issue.url}`}
                                      target="_blank"
                                      rel="noopener noreferrer"
                                    >
                                      #{issue.number}{' '}
                                    </Anchor>
                                  </Text>
                                </Tooltip>
                              ))}
                            </Text>
                          )}
                        </Box>
                      );
                    })}
                  </Stack>
                )}

                {/* Outdated languages */}
                {outdatedLangs.length > 0 && (
                  <Stack gap="xs">
                    {outdatedLangs.map((code) => {
                      const translation = article.translations[code.value];
                      const translationPath = article.englishPath.replace(
                        '/en/',
                        `/${code.value}/`
                      );

                      const handleDiffClick = () => {
                        const params = new URLSearchParams({
                          category: selectedArticleCategory,
                          translationPath,
                          language: code.value,
                        });
                        navigate(`/detail?${params.toString()}`);
                      };

                      return (
                        <Box
                          key={code.value}
                          p="xs"
                          bg="rgba(255, 165, 0, 0.1)"
                          style={{
                            borderLeft: '4px solid #f59e0b',
                            borderRadius: 4,
                          }}
                        >
                          <Stack gap={4}>
                            <Group justify="space-between" align="center">
                              <Group gap="xs">
                                <StatusBadge status="outdated" />
                                <Anchor
                                  href={`https://github.com/kubernetes/website/blob/main/${translationPath}`}
                                  target="_blank"
                                  rel="noopener noreferrer"
                                  underline="hover"
                                  c="inherit"
                                >
                                  <Text fw={600} c="dark" size="sm" pb={2}>
                                    {getLanguageName(code.value)}
                                  </Text>
                                </Anchor>
                                <Text size="xs" c="dimmed">
                                  Updated at{' '}
                                  {translation?.targetLatestDate
                                    ? formatDateISO(translation.targetLatestDate)
                                    : ''}{' '}
                                  (UTC)
                                </Text>
                                <Group gap="0">
                                  {translation?.translationUrl && (
                                    <ActionIcon
                                      component="a"
                                      href={translation.translationUrl || ''}
                                      target="_blank"
                                      rel="noopener noreferrer"
                                      size="xs"
                                      radius="xs"
                                      c="gray"
                                      variant="subtle"
                                      title="View on Kubernetes site"
                                    >
                                      <IconExternalLink size={14} />
                                    </ActionIcon>
                                  )}
                                  {translation.status === 'outdated' && (
                                    <>
                                      <ActionIcon
                                        onClick={handleDiffClick}
                                        size="xs"
                                        radius="xs"
                                        c="blue"
                                        variant="subtle"
                                        title="View translation diff"
                                        style={{ cursor: 'pointer' }}
                                      >
                                        <IconGitBranch size={14} />
                                      </ActionIcon>
                                      <GitHubIssueButton
                                        englishPath={article.englishPath}
                                        langCode={code.value}
                                        variant="update"
                                      />
                                    </>
                                  )}
                                </Group>
                              </Group>
                            </Group>
                            <Text size="xs" c="dimmed">
                              {[
                                translation?.totalChangeLines &&
                                  `${translation.totalChangeLines} lines changed`,
                                translation?.commitsBehind &&
                                  translation?.daysBehind &&
                                  `${translation.commitsBehind} commits`,
                                translation?.daysBehind && `${translation.daysBehind} days behind`,
                              ]
                                .filter(Boolean)
                                .join(' • ')}
                            </Text>
                            {translation.prs.length > 0 && (
                              <Text size="xs" c="dimmed">
                                PR:{' '}
                                {translation?.prs.map((pr) => (
                                  <Tooltip key={pr.number} label={`PR #${pr.number} - ${pr.title}`}>
                                    <Text key={pr.number} size="xs" c="dimmed" component="span">
                                      <Anchor
                                        href={`${pr.url}`}
                                        target="_blank"
                                        rel="noopener noreferrer"
                                      >
                                        #{pr.number}{' '}
                                      </Anchor>
                                    </Text>
                                  </Tooltip>
                                ))}
                              </Text>
                            )}
                            {translation.issues.length > 0 && (
                              <Text size="xs" c="dimmed">
                                Issue:{' '}
                                {translation.issues.map((issue) => (
                                  <Tooltip
                                    label={`Issue #${issue.number} - ${issue.title}`}
                                    key={issue.number}
                                  >
                                    <Text size="xs" c="dimmed" component="span">
                                      <Anchor
                                        href={`${issue.url}`}
                                        target="_blank"
                                        rel="noopener noreferrer"
                                      >
                                        #{issue.number}{' '}
                                      </Anchor>
                                    </Text>
                                  </Tooltip>
                                ))}
                              </Text>
                            )}
                          </Stack>
                        </Box>
                      );
                    })}
                  </Stack>
                )}

                {/* Not translated languages */}
                {notTranslatedLangs.length > 0 && (
                  <Box
                    p="xs"
                    bg="#f8fafc"
                    style={{
                      borderRadius: 4,
                      borderLeft: '4px solid #cbd5e1',
                    }}
                  >
                    <Group gap="xs" mb="xs">
                      <Text size="xs" fw={600} c="dimmed">
                        — Not translated
                      </Text>
                    </Group>
                    <Group gap="xs">
                      {notTranslatedLangs.map((code) => {
                        const translation = article.translations[code.value];
                        return (
                          <Box
                            key={code.value}
                            px="xs"
                            py={4}
                            bg="white"
                            style={{
                              borderRadius: 16,
                              fontSize: 11,
                              color: '#64748b',
                              border: '1px solid #e2e8f0',
                              fontWeight: 500,
                            }}
                          >
                            <Group gap={4} align="center">
                              {getLanguageName(code.value)}
                              <GitHubIssueButton
                                englishPath={article.englishPath}
                                langCode={code.value}
                                variant="new"
                              />
                            </Group>
                            {translation.issues.length > 0 && (
                              <Text size="xs" c="dimmed" mt="xs" component="span">
                                {' '}
                                {translation.issues.map((issue) => (
                                  <Tooltip
                                    label={`Issue #${issue.number} - ${issue.title}`}
                                    key={issue.number}
                                  >
                                    <Text size="xs" c="dimmed" component="span">
                                      <Anchor
                                        href={`${issue.url}`}
                                        target="_blank"
                                        rel="noopener noreferrer"
                                      >
                                        #{issue.number}{' '}
                                      </Anchor>
                                    </Text>
                                  </Tooltip>
                                ))}
                              </Text>
                            )}
                            {translation?.prs.length > 0 && (
                              <Text size="xs" c="dimmed" component="span">
                                {translation?.prs.map((pr) => (
                                  <Tooltip key={pr.number} label={`PR #${pr.number} - ${pr.title}`}>
                                    <Text size="xs" c="dimmed" component="span">
                                      <Anchor
                                        href={`${pr.url}`}
                                        target="_blank"
                                        rel="noopener noreferrer"
                                      >
                                        #{pr.number}{' '}
                                      </Anchor>
                                    </Text>
                                  </Tooltip>
                                ))}
                              </Text>
                            )}
                          </Box>
                        );
                      })}
                    </Group>
                  </Box>
                )}
              </Stack>
            )}
          </Card>
        );
      })}
    </Stack>
  );
};
