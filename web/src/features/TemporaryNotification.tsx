import { IconInfoCircle } from '@tabler/icons-react';
import { Alert, Anchor } from '@mantine/core';

const CASE_STUDIES_REMOVAL_PR_URL = 'https://github.com/kubernetes/website/pull/57043';

export const TemporaryNotification = () => {
  return (
    <Alert
      withCloseButton
      variant="light"
      color="orange"
      title="Case Studies Removed"
      icon={<IconInfoCircle size={16} />}
      mb={20}
    >
      The case studies section has been removed from the English documentation (
      <Anchor href={CASE_STUDIES_REMOVAL_PR_URL} size="sm" target="_blank" rel="noopener">
        kubernetes/website#57043
      </Anchor>
      ), so it is no longer tracked here. Please consider removing the localized versions as well.
    </Alert>
  );
};
