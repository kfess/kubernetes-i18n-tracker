import { IconInfoCircle } from '@tabler/icons-react';
import { Alert } from '@mantine/core';

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
      The case studies section has been removed from the English documentation, so it is no longer
      tracked here. Please consider removing the localized versions as well.
    </Alert>
  );
};
