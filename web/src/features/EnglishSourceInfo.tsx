import { IconClock, IconExternalLink, IconEye, IconUsers } from '@tabler/icons-react';
import { ActionIcon, Anchor, Badge, Group, rem, Stack, Text, Tooltip } from '@mantine/core';
import { useMediaQuery } from '@mantine/hooks';
import { type ArticleTranslation } from '@/features/translations';
import { formatDateISO, formatSecondsToMMSS } from '@/utils/date';

interface MetricBadgeProps {
  label: string;
  color: string;
  icon: React.ComponentType<{ size: number }>;
  value: string;
}

const MetricBadge = ({ label, color, icon: Icon, value }: MetricBadgeProps) => (
  <Tooltip label={label} withArrow>
    <Badge
      variant="light"
      color={color}
      leftSection={<Icon size={12} />}
      size="sm"
      style={{ textTransform: 'none' }}
    >
      {value}
    </Badge>
  </Tooltip>
);

interface Props {
  article: ArticleTranslation;
}

export const EnglishSourceInfo = ({ article }: Props) => {
  const isMobile = useMediaQuery('(max-width: 768px)');

  return (
    <Stack gap={4}>
      <Anchor
        href={`https://github.com/kubernetes/website/blob/main/${article.englishPath}`}
        target="_blank"
        rel="noopener noreferrer"
        underline="hover"
        c="inherit"
      >
        <Text
          size="sm"
          fw={500}
          title={article.englishPath}
          style={{
            maxWidth: isMobile ? '' : rem(250),
            whiteSpace: 'normal',
            overflowWrap: 'break-word',
            wordBreak: 'keep-all',
          }}
        >
          {article.englishPath.replace('content/en/', '')}
        </Text>
      </Anchor>
      <Group gap="3" align="center" wrap="nowrap" c="dimmed">
        <Text size="xs">
          Updated: {formatDateISO(article.translations.en?.englishLatestDate)} (UTC)
        </Text>
        {article.englishUrl && (
          <ActionIcon
            component="a"
            href={article.englishUrl || ''}
            target="_blank"
            rel="noopener noreferrer"
            size="xs"
            radius="xs"
            variant="subtle"
            color="gray"
            title="View on Kubernetes site"
          >
            <IconExternalLink size={14} />
          </ActionIcon>
        )}
      </Group>
      <Group gap="xs" align="center">
        <MetricBadge
          label="Views"
          color="blue"
          icon={IconEye}
          value={article.translations.en.views.toLocaleString()}
        />
        <MetricBadge
          label="New Users"
          color="green"
          icon={IconUsers}
          value={article.translations.en.newUsers.toLocaleString()}
        />
        <MetricBadge
          label="Average Session Duration"
          color="indigo"
          icon={IconClock}
          value={formatSecondsToMMSS(
            Math.round(article.translations.en.averageSessionDuration) || 0
          )}
        />
      </Group>
    </Stack>
  );
};
