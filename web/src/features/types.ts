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
export const detectionModes = ['combined', 'git', 'structure'] as const;
export type DetectionMode = (typeof detectionModes)[number];
export const defaultDetectionMode: DetectionMode = 'combined';
export const detectionModeLabels: Record<DetectionMode, string> = {
  combined: 'Git + Structure',
  git: 'Git only',
  structure: 'Structure only',
};
