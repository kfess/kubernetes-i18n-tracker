import { useEffect, useRef } from 'react';
import { IconExternalLink, IconGitBranch, IconInfoCircle } from '@tabler/icons-react';
import { Diff2HtmlUI } from 'diff2html/lib/ui/js/diff2html-ui';
import { useLocation } from 'react-router-dom';
import {
  Alert,
  Anchor,
  Badge,
  Breadcrumbs,
  Container,
  Group,
  Paper,
  Stack,
  Text,
  Title,
} from '@mantine/core';
import { useMediaQuery } from '@mantine/hooks';
import { basename } from '@/const';
import { useFetchDiff } from '@/features/hooks/useFetchDiff';
import { useFetchTranslationArticles } from '@/features/hooks/useFetchTranslationArticles';
import { StatusBadge } from '@/features/StatusBadge';
import { ArticleCategory } from '@/features/translations';
import { formatDateISO } from '@/utils/date';

import 'diff2html/bundles/css/diff2html.min.css';

export const DetailPage = () => {
  const isMobile = useMediaQuery('(max-width: 768px)');

  const location = useLocation();
  const params = new URLSearchParams(location.search);
  const category = params.get('category') || '';
  const translationPath = decodeURIComponent(params.get('translationPath') || '');
  const language = params.get('language') || '';

  const breadcrumbItems = [
    { title: 'Home', href: basename },
    { title: 'Detail', href: basename + location.pathname + location.search },
  ].map((item, index) => (
    <Anchor href={item.href} key={index}>
      {item.title}
    </Anchor>
  ));

  const diff = useFetchDiff(category);
  const articles = useFetchTranslationArticles(category as ArticleCategory);

  const diffData = diff[translationPath];
  const diffString = diffData?.diff || '';

  const targetArticle = articles.find(
    (article) => article.englishPath.replace('/en/', `/${language}/`) === translationPath
  );
  const translationInfo =
    targetArticle?.translations?.[language as keyof typeof targetArticle.translations];

  const containerRef = useRef(null);

  useEffect(() => {
    if (containerRef.current && diffString) {
      const diff2htmlUi = new Diff2HtmlUI(containerRef.current, diffString, {
        drawFileList: true,
        matching: 'lines',
        outputFormat: isMobile ? 'line-by-line' : 'side-by-side',
      });
      diff2htmlUi.draw();
    }
  }, [diffString, isMobile]);

  const fileName = translationPath.split('/').pop() || 'Unknown file';

  return (
    <Container fluid>
      <Breadcrumbs mb="md">{breadcrumbItems}</Breadcrumbs>
      <Alert
        withCloseButton
        variant="light"
        color="orange"
        title="Translation Diff Viewer"
        icon={<IconInfoCircle size={16} />}
        mb={20}
      >
        <div>
          This page shows the differences between the English version at the time when the
          translation was last updated and the latest English version. These changes represent what
          needs to be translated to bring the "{language}" version up to date with the latest
          English content.
        </div>
        <div>Please note that this is still an experimental feature and may contain errors.</div>
      </Alert>
      <Paper withBorder p="md" mb="md">
        <Stack gap="md">
          <div>
            <Group justify="space-between" align="baseline" mb="xs">
              <Title order={2} size="h3">
                {fileName}
              </Title>
              {translationInfo?.targetLatestDate && (
                <Text size="xs" c="dimmed">
                  Last updated: {formatDateISO(translationInfo.targetLatestDate)}
                </Text>
              )}
            </Group>
            <Text c="gray.6" size="sm" mb="xs">
              {translationPath}
            </Text>
            {(translationInfo?.refEnglishCommitHash ||
              translationInfo?.englishLatestCommitHash) && (
              <Group gap="md" mb="sm" wrap="wrap">
                {translationInfo?.refEnglishCommitHash && translationInfo?.refEnglishCommitDate && (
                  <Anchor
                    href={`https://github.com/kubernetes/website/blob/${translationInfo.refEnglishCommitHash}/${targetArticle?.englishPath || ''}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    size="xs"
                    fw={500}
                    c="gray.7"
                    td="none"
                  >
                    <Group gap={4} wrap="nowrap">
                      <IconGitBranch size={12} />
                      Ref English: {translationInfo.refEnglishCommitHash.substring(0, 7)} (
                      {formatDateISO(translationInfo.refEnglishCommitDate)})
                    </Group>
                  </Anchor>
                )}
                {translationInfo?.englishLatestCommitHash && translationInfo?.englishLatestDate && (
                  <Anchor
                    href={`https://github.com/kubernetes/website/blob/${translationInfo.englishLatestCommitHash}/${targetArticle?.englishPath || ''}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    size="xs"
                    fw={500}
                    c="gray.7"
                    td="none"
                  >
                    <Group gap={4} wrap="nowrap">
                      <IconGitBranch size={12} />
                      Latest English: {translationInfo.englishLatestCommitHash.substring(0, 7)} (
                      {formatDateISO(translationInfo.englishLatestDate)})
                    </Group>
                  </Anchor>
                )}
              </Group>
            )}
            {translationInfo && (
              <Group gap="xs" wrap="wrap">
                <StatusBadge status={translationInfo.status} />
                {translationInfo.status === 'outdated' && (
                  <>
                    <Badge variant="light" color="blue" style={{ textTransform: 'none' }}>
                      {translationInfo.commitsBehind} commits behind
                    </Badge>
                    <Badge variant="light" color="blue" style={{ textTransform: 'none' }}>
                      {translationInfo.daysBehind} days behind
                    </Badge>
                    <Badge variant="light" color="blue" style={{ textTransform: 'none' }}>
                      {translationInfo.totalChangeLines.toLocaleString()} lines changed
                    </Badge>
                  </>
                )}
              </Group>
            )}
          </div>

          <Group gap="sm" wrap="wrap" justify={isMobile ? 'flex-start' : 'flex-end'}>
            {targetArticle?.englishUrl && (
              <Anchor
                href={targetArticle.englishUrl}
                target="_blank"
                rel="noopener noreferrer"
                size="xs"
                c="gray.7"
                td="none"
              >
                <Group gap={4} wrap="nowrap">
                  <IconExternalLink size={12} />
                  English (Live)
                </Group>
              </Anchor>
            )}

            {translationInfo?.translationUrl && (
              <Anchor
                href={translationInfo.translationUrl}
                target="_blank"
                rel="noopener noreferrer"
                size="xs"
                c="gray.7"
                td="none"
              >
                <Group gap={4} wrap="nowrap">
                  <IconExternalLink size={12} />
                  {language} (Live)
                </Group>
              </Anchor>
            )}

            <Anchor
              href={`https://github.com/kubernetes/website/blob/main/${targetArticle?.englishPath || ''}`}
              target="_blank"
              rel="noopener noreferrer"
              size="xs"
              c="gray.7"
              td="none"
            >
              <Group gap={4} wrap="nowrap">
                <IconGitBranch size={12} />
                English (GitHub)
              </Group>
            </Anchor>

            <Anchor
              href={`https://github.com/kubernetes/website/blob/main/${translationPath}`}
              target="_blank"
              rel="noopener noreferrer"
              size="xs"
              c="gray.7"
              td="none"
            >
              <Group gap={4} wrap="nowrap">
                <IconGitBranch size={12} />
                {language} (GitHub)
              </Group>
            </Anchor>
          </Group>
        </Stack>
      </Paper>

      {diffString ? (
        <div ref={containerRef} />
      ) : (
        <Alert variant="light" color="blue" title="No diff available">
          No differences found for this file. The translation might be up to date or the diff data
          might not be available.
        </Alert>
      )}
    </Container>
  );
};
