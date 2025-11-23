export const sortModes = [
  'default',
  'views',
  'newUsers',
  'updatedAt',
  'averageSessionDuration',
] as const;
export type SortMode = (typeof sortModes)[number];
export type SortDirection = 'asc' | 'desc';
export type IssueStatus = 'all' | 'withIssues' | 'withoutIssues';
export type PrStatus = 'all' | 'withPr' | 'withoutPr';
