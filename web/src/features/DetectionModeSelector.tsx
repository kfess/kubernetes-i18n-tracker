import { useState } from 'react';
import {
  IconAdjustmentsHorizontal,
  IconCheck,
  IconGitBranch,
  IconLayersIntersect,
  IconListTree,
  type Icon,
} from '@tabler/icons-react';
import { ActionIcon, Group, Indicator, Menu, Stack, Text, Tooltip } from '@mantine/core';
import { useDetectionMode } from '@/features/hooks/useDetectionMode';
import {
  defaultDetectionMode,
  detectionModeLabels,
  detectionModes,
  type DetectionMode,
} from '@/features/types';

const modeDetails: Record<DetectionMode, { icon: Icon; description: string }> = {
  combined: {
    icon: IconLayersIntersect,
    description:
      'Git history decides what is outdated. Structural gaps mark up-to-date pages as possibly outdated.',
  },
  git: {
    icon: IconGitBranch,
    description: 'Outdated when the English page changed after the translation was last updated.',
  },
  structure: {
    icon: IconListTree,
    description: 'Compares headings, code blocks and version references with the English page.',
  },
};

export const DetectionModeSelector = () => {
  const [selectedDetectionMode, setDetectionMode] = useDetectionMode();
  const [opened, setOpened] = useState(false);

  const isDefault = selectedDetectionMode === defaultDetectionMode;

  return (
    <Menu position="bottom-end" withinPortal width={350} opened={opened} onChange={setOpened}>
      <Menu.Target>
        {/* The tooltip would cover the dropdown, so it is hidden while the menu is open. */}
        <Tooltip
          label={`Detection mode: ${detectionModeLabels[selectedDetectionMode]}`}
          position="bottom"
          withArrow
          offset={10}
          disabled={opened}
        >
          {/* The dot tells that a non-default mode is changing the statuses on screen. */}
          <Indicator disabled={isDefault} color="orange" size={9} offset={4} withBorder>
            <ActionIcon variant="default" radius="md" size="lg" aria-label="Detection mode">
              <IconAdjustmentsHorizontal size={20} />
            </ActionIcon>
          </Indicator>
        </Tooltip>
      </Menu.Target>
      <Menu.Dropdown>
        <Menu.Label>Detection mode</Menu.Label>
        {detectionModes.map((mode) => {
          const isSelected = selectedDetectionMode === mode;
          const { icon: ModeIcon, description } = modeDetails[mode];

          return (
            <Menu.Item
              key={mode}
              aria-current={isSelected ? 'true' : undefined}
              onClick={() => setDetectionMode(mode)}
              leftSection={
                <ModeIcon
                  size={18}
                  color={isSelected ? 'var(--mantine-color-blue-6)' : undefined}
                />
              }
              rightSection={isSelected ? <IconCheck size={16} /> : null}
              style={{
                backgroundColor: isSelected ? 'var(--mantine-color-blue-0)' : undefined,
              }}
            >
              <Stack gap={2}>
                <Group gap="xs">
                  <Text size="sm" fw={isSelected ? 600 : 400} c={isSelected ? 'blue' : undefined}>
                    {detectionModeLabels[mode]}
                  </Text>
                  {mode === defaultDetectionMode && (
                    <Text size="xs" c="dimmed">
                      (Default, Recommended)
                    </Text>
                  )}
                </Group>
                <Text size="xs" c="dimmed" lh={1.35}>
                  {description}
                </Text>
              </Stack>
            </Menu.Item>
          );
        })}

        <Menu.Divider />

        <Text size="xs" c="dimmed" ta="center" px="sm" py={6}>
          Changes how every page decides a translation's status
        </Text>
      </Menu.Dropdown>
    </Menu>
  );
};
