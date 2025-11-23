import { JSX } from 'react';
import {
  IconArrowsSort,
  IconCalendar,
  IconClock,
  IconEye,
  IconHome,
  IconSortAscending,
  IconSortDescending,
  IconUsers,
} from '@tabler/icons-react';
import { ActionIcon, Menu } from '@mantine/core';
import { type SortDirection, type SortMode } from './types';

interface Props {
  sortMode: SortMode;
  setSortMode: (mode: SortMode) => void;
  sortDirection: SortDirection;
  setSortDirection: (direction: SortDirection) => void;
}

const sortOptions: { mode: SortMode; label: string; icon: JSX.Element }[] = [
  { mode: 'default', label: 'Default', icon: <IconHome size={14} /> },
  { mode: 'views', label: 'Views', icon: <IconEye size={14} /> },
  { mode: 'newUsers', label: 'New Users', icon: <IconUsers size={14} /> },
  {
    mode: 'averageSessionDuration',
    label: 'Session Duration',
    icon: <IconClock size={14} />,
  },
  { mode: 'updatedAt', label: 'Updated At', icon: <IconCalendar size={14} /> },
];

export const SortMenu = ({ sortMode, setSortMode, sortDirection, setSortDirection }: Props) => {
  const handleClick = (mode: SortMode) => {
    const isSame = sortMode === mode;
    setSortMode(mode);
    setSortDirection(isSame ? (sortDirection === 'asc' ? 'desc' : 'asc') : 'desc');
  };

  const renderSortIcon = (mode: SortMode) => {
    if (sortMode !== mode) {
      return null;
    }
    return sortDirection === 'desc' ? (
      <IconSortDescending size={14} />
    ) : (
      <IconSortAscending size={14} />
    );
  };

  return (
    <Menu shadow="md" width={200}>
      <Menu.Target>
        <ActionIcon variant="light" size="lg" title={`Sort by ${sortMode} (${sortDirection})`}>
          <IconArrowsSort size={17} />
        </ActionIcon>
      </Menu.Target>
      <Menu.Dropdown>
        <Menu.Label>Sort By</Menu.Label>
        {sortOptions.map(({ mode, label, icon }) => (
          <Menu.Item
            key={label}
            leftSection={icon}
            onClick={() => handleClick(mode)}
            rightSection={renderSortIcon(mode)}
            bg={sortMode === mode ? 'var(--mantine-color-gray-1)' : undefined}
          >
            {label}
          </Menu.Item>
        ))}
      </Menu.Dropdown>
    </Menu>
  );
};
