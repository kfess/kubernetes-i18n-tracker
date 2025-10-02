import os
from collections import defaultdict
from dataclasses import asdict, dataclass

from dotenv import load_dotenv
from github import Auth, Github
from log import logger

TOO_MANY_FILES_CHANGED = 1000
TOO_MANY_COMMITS = 1000


@dataclass
class GitHubPullRequest:
    """Represents a GitHub pull request."""

    number: int
    title: str
    url: str
    files: list[str]


def _get_prs(repo_name: str = "kubernetes/website") -> list[GitHubPullRequest]:
    """Get open pull requests for a GitHub repository."""
    load_dotenv()

    auth = Auth.Token(os.getenv("KUBERNETES_WEBSITE_READ_GITHUB_TOKEN"))
    g = Github(auth=auth)
    repository = g.get_repo(repo_name)

    logger.info("Start fetching pull requests from %s", repo_name)

    raw_pull_requests = repository.get_pulls(
        state="open", sort="created", direction="desc"
    )

    pull_requests: list[GitHubPullRequest] = []

    for pr in raw_pull_requests:
        number = pr.number
        title = pr.title
        url = pr.html_url
        label_names = [label.name for label in pr.labels]
        files = pr.get_files()
        commits = pr.commits

        if "area/localization" not in label_names or "cncf-cla: yes" not in label_names:
            continue

        file_changes = files.totalCount

        # We suppose wrong PR if it has too many files changed
        if file_changes >= TOO_MANY_FILES_CHANGED or commits >= TOO_MANY_COMMITS:
            logger.warning(
                "Skipping PR #%d - Too many files changed: %d, commits: %d",
                number,
                file_changes,
                commits,
            )
            continue

        pull_requests.append(
            GitHubPullRequest(
                number=number,
                title=title,
                url=url,
                files=[f.filename for f in files],
            )
        )

    g.close()
    logger.info("Finished fetching pull requests from %s", repo_name)

    return pull_requests


def get_prs_by_file(
    repo_name: str = "kubernetes/website",
) -> dict[str, list[GitHubPullRequest]]:
    """Group PRs by the files they modify."""
    file_to_prs: dict[str, list[GitHubPullRequest]] = defaultdict(list)

    prs = _get_prs(repo_name)

    for pr in prs:
        for file_path in pr.files:
            file_to_prs[file_path].append(asdict(pr))

    return dict(file_to_prs)


if __name__ == "__main__":
    print(get_prs_by_file())
