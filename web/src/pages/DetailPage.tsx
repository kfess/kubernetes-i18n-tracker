import React, { useEffect, useRef } from 'react';
import { IconInfoCircle } from '@tabler/icons-react';
import { Diff2HtmlUI } from 'diff2html/lib/ui/js/diff2html-ui';
import { useLocation } from 'react-router-dom';
import { Alert, Container } from '@mantine/core';
import { useMediaQuery } from '@mantine/hooks';
import { useFetchDiff } from '@/features/hooks/useFetchDiff';

import 'diff2html/bundles/css/diff2html.min.css';

export const DetailPage = () => {
  const isMobile = useMediaQuery('(max-width: 768px)');

  const location = useLocation();
  const params = new URLSearchParams(location.search);
  const category = params.get('category') || '';
  const translationPath = decodeURIComponent(params.get('translationPath') || '');
  const language = params.get('language') || '';

  const diff = useFetchDiff(category);
  const diffString = diff[translationPath].diff;

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

  return (
    <Container fluid>
      <Alert
        withCloseButton
        variant="light"
        color="orange"
        title="Experimental: Translation Diff Viewer"
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
      <div ref={containerRef} />
    </Container>
  );
};
