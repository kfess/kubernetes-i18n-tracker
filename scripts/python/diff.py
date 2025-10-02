import subprocess


def diff(l_commit: str, r_commit: str, file_path: str) -> str:
    """Generate a diff between two commits.

    Args:
    ----
        l_commit (str): The left commit hash or reference.
        r_commit (str): The right commit hash or reference.
        file_path (str): The path to the file to diff.

    Returns:
    -------
        str: The diff output as a string.

    """
    try:
        result = subprocess.run(
            ["git", "diff", l_commit, r_commit, "--", file_path],
            capture_output=True,
            text=True,
            check=True,
            cwd="../../k8s-repo/website",
        )
        return result.stdout
    except subprocess.CalledProcessError as e:
        print(f"Error generating diff: {e.stderr}")
        return ""
