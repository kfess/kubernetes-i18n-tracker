import { useEffect, useRef, useState } from 'react';
import { Diff2HtmlUI } from 'diff2html/lib/ui/js/diff2html-ui';

import 'diff2html/bundles/css/diff2html.min.css';

interface Props {
  diffString?: string;
}

// export const DevPage = ({ diffString }: Props) => {
//   const containerRef = useRef(null);

//   useEffect(() => {
//     if (containerRef.current && diffString) {
//       containerRef.current.innerHTML = '';

//       const diff2htmlUi = new Diff2HtmlUI(containerRef.current, diffString, {
//         drawFileList: true,
//         matching: 'lines',
//         // outputFormat: 'side-by-side',
//         outputFormat: 'line-by-line',
//       });
//       diff2htmlUi.draw();
//     }
//   }, [diffString]);

//   return <div ref={containerRef} />;
// };

interface DiffJson {
  diff: string;
}

export const DevPage = ({ diffString }: Props) => {
  const [diff, setDiff] = useState<string>('');

  useEffect(() => {
    fetch('https://api.github.com/repos/kfess/kubernetes-i18n-tracker/releases/tags/test-tag')
      .then((res) => res.json())
      .then((release) => fetch(release.assets[0].browser_download_url))
      .then((res) => res.json())
      .then((data) => setDiff(data.diff));
  }, []);

  //   useEffect(() => {
  //     fetch('https://github.com/kfess/kubernetes-i18n-tracker/releases/download/test-tag/diff.json')
  //       .then((res) => res.json())
  //       .then((data: DiffJson) => setDiff(data.diff));
  //   }, []);

  //   const containerRef = useRef(null);

  //   useEffect(() => {
  //     if (containerRef.current && diffString) {
  //       containerRef.current.innerHTML = '';

  //       const diff2htmlUi = new Diff2HtmlUI(containerRef.current, diffString, {
  //         drawFileList: true,
  //         matching: 'lines',
  //         // outputFormat: 'side-by-side',
  //         outputFormat: 'line-by-line',
  //       });
  //       diff2htmlUi.draw();
  //     }
  //   }, [diffString]);

  //   return <div ref={containerRef} />;

  return <div>{diff}</div>;
};
