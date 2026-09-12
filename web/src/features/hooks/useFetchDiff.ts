import blogDiffs from '@/data/output/diff/blog_diff.json';
import careerDiffs from '@/data/output/diff/careers_diff.json';
import communityDiffs from '@/data/output/diff/community_diff.json';
import docsConceptsDiffs from '@/data/output/diff/docs_concepts_diff.json';
import docsContributionDiffs from '@/data/output/diff/docs_contribute_diff.json';
import docsHomeDiffs from '@/data/output/diff/docs_home_diff.json';
import docsReferenceDiffs from '@/data/output/diff/docs_reference_diff.json';
import docsSetupDiffs from '@/data/output/diff/docs_setup_diff.json';
import docsTasksDiffs from '@/data/output/diff/docs_tasks_diff.json';
import docsTutorialsDiffs from '@/data/output/diff/docs_tutorials_diff.json';
import examplesDiffs from '@/data/output/diff/examples_diff.json';
import includesDiffs from '@/data/output/diff/includes_diff.json';
import partnerDiffs from '@/data/output/diff/partners_diff.json';
import releaseDiffs from '@/data/output/diff/releases_diff.json';
import trainingDiffs from '@/data/output/diff/training_diff.json';

import { Diff } from '@/features/translations';

const diffs: Record<string, Diff> = {
  blog: blogDiffs as Diff,
  community: communityDiffs as Diff,
  examples: examplesDiffs as Diff,
  docsConcept: docsConceptsDiffs as Diff,
  docsContribute: docsContributionDiffs as Diff,
  docsHome: docsHomeDiffs as Diff,
  docsTask: docsTasksDiffs as Diff,
  docsReference: docsReferenceDiffs as Diff,
  docsSetup: docsSetupDiffs as Diff,
  docsTutorial: docsTutorialsDiffs as Diff,
  includes: includesDiffs as Diff,
  partner: partnerDiffs as Diff,
  release: releaseDiffs as Diff,
  training: trainingDiffs as Diff,
  career: careerDiffs as Diff,
};

export const useFetchDiff = (articleCategory: string): Diff => {
  return diffs[articleCategory] || {};
};
