import json
import re
from collections import defaultdict
from datetime import datetime, timezone
from pathlib import Path, PurePosixPath
from typing import Any

from diff import diff
from issue import GitHubIssue, get_issues_by_file
from log import logger
from page_view import PageView, summarize_view
from pull_requests import GitHubPullRequest, get_prs_by_file
from translation_status import TranslationStatusResult
from url_builder import build_url
from utils import convert_keys_to_camel_case, serialize_datetime


def should_process(result: TranslationStatusResult) -> bool:
    """Determine if a translation result should be processed.

    Args:
    ----
        result (TranslationStatusResult): The translation status result to check.

    Returns:
    -------
        bool: True if the result should be processed, False otherwise.

    """

    def is_supported_extension(path: PurePosixPath) -> bool:
        return path.suffix in {".md", ".html"}

    def is_known_category(category: str) -> bool:
        return category != "unknown"

    path = PurePosixPath(result["english_path"].lower())

    return is_supported_extension(path) and is_known_category(result["category"])


def extract_blog_date_from_en_path(en_path: str) -> str:
    """Extract date from blog post path.

    Args:
    ----
        en_path (str): The English file path.

    Returns:
    -------
        str: The extracted date in YYYY-MM-DD format.

    Example:
    -------
        content/en/blog/_posts/2024-10-02-xxxx.md -> 2024-10-02
        content/ja/blog/_posts/2025-03-26.md -> 2025-03-26

    """
    match = re.search(r"content/en/blog/_posts/(\d{4}-\d{2}-\d{2})", en_path)
    return match.group(1) if match else "0000-00-00"


def extract_docs_subcategory(en_path: str) -> str | None:
    """Extract subcategory from docs path.

    Args:
    ----
        en_path (str): The English file path.

    Returns:
    -------
        str | None: The extracted subcategory, or
                    None if not found or if it's just a filename.

    Example:
    -------
        content/en/docs/getting-started/installation.md -> getting-started
        content/en/docs/api/reference.md -> api
        content/en/docs/tutorial/basic.md -> tutorial
        content/en/docs/overview.md -> None  # Just a filename, not a subcategory

    """
    match = re.search(r"content/[^/]+/docs/([^/]+)", en_path)
    if not match:
        return None

    potential_subcategory = match.group(1)

    # Check if this is just a filename (has .md extension)
    if potential_subcategory.endswith(".md"):
        return None

    return potential_subcategory


def build_category_name(original_category: str, english_path: str) -> str:
    """Build the final category name, expanding docs with subcategories.

    Args:
    ----
        original_category (str): The original category name.
        english_path (str): The English file path.

    Returns:
    -------
        str: The final category name (e.g., 'docs_getting-started' for docs).

    """
    if original_category == "docs":
        subcategory = extract_docs_subcategory(english_path)
        if subcategory:
            return f"docs_{subcategory}"
        else:
            return "docs_misc"
    return original_category


def create_matrix_data(
    results: dict[str, TranslationStatusResult],
    issues_by_file: dict[str, list[GitHubIssue]],
    prs_by_file: dict[str, list[GitHubPullRequest]],
    existing_urls: set[str],
) -> dict[str, dict[str, Any]]:
    """Create matrix data grouped by category."""
    matrix_data = defaultdict(
        lambda: {
            "last_updated": datetime.now(tz=timezone.utc).isoformat(),  # noqa: UP017
            "articles": [],
        }
    )
    page_views = summarize_view("../../data/master/page_view.csv", existing_urls)

    articles_by_english_path = defaultdict(dict)

    for result in results.values():
        english_path = result["english_path"]
        language = result["language"]
        target_path = result["english_path"].replace(
            "content/en/", f"content/{language}/"
        )
        english_url = build_url(english_path, "en", existing_urls)
        translation_url = build_url(english_path, language, existing_urls)
        page_view = page_views.get(
            translation_url,
            PageView(views=0, new_users=0, average_session_duration=0.0),
        )

        translation_data = {
            "status": result["status"],
            "severity": result["severity"],
            "days_behind": result["days_behind"],
            "commits_behind": result["commits_behind"],
            "total_change_lines": result["total_change_lines"],
            "target_latest_date": result["target_latest_date"],
            "english_latest_date": result["english_latest_date"],
            "translation_url": translation_url,
            "views": page_view.views,
            "new_users": page_view.new_users,
            "average_session_duration": page_view.average_session_duration,
            "issues": issues_by_file.get(target_path, []),
            "prs": prs_by_file.get(target_path, []),
            "english_latest_commit_hash": result.get("english_latest_commit_hash"),
            "ref_english_commit_hash": result.get("ref_english_commit_hash"),
            "ref_english_commit_date": result.get("ref_english_commit_date"),
        }

        articles_by_english_path[english_path][language] = translation_data

    for english_path, translations in articles_by_english_path.items():
        original_category = next(
            result["category"]
            for result in results.values()
            if result["english_path"] == english_path
        )

        category_name = build_category_name(original_category, english_path)
        english_url = build_url(english_path, "en", existing_urls)

        article_data = {
            "english_path": english_path,
            "english_url": english_url,
            "translations": translations,
        }
        matrix_data[category_name]["articles"].append(article_data)

    # Sort blog articles by date
    for category, data in matrix_data.items():
        if category == "blog":
            data["articles"].sort(
                key=lambda x: extract_blog_date_from_en_path(x["english_path"]),
                reverse=True,
            )

    return dict(matrix_data)


def create_detail_data(
    result: TranslationStatusResult,
    existing_urls: set[str],
    issues_by_file: dict[str, list[GitHubIssue]],
    prs_by_file: dict[str, Any],
) -> dict[str, Any]:
    """Create detail data for a single translation result."""
    english_path = result["english_path"]
    language = result["language"]

    # Generate URLs
    english_url = build_url(english_path, "en", existing_urls)
    translation_url = build_url(english_path, language, existing_urls)

    return {
        "target_path": result["target_path"],
        "english_path": result["english_path"],
        "english_url": english_url,
        "translation_url": translation_url,
        "target_latest_date": result["target_latest_date"],
        "english_latest_date": result["english_latest_date"],
        "days_behind": result["days_behind"],
        "commits_behind": result["commits_behind"],
        "total_change_lines": result["total_change_lines"],
        "insertions_behind_lines": result["insertions_behind_lines"],
        "deletions_behind_lines": result["deletions_behind_lines"],
        "status": result["status"],
        "severity": result["severity"],
        "missing_commits": result["missing_commits"],
        "issues": issues_by_file.get(result["target_path"], []),
        "prs": prs_by_file.get(result["target_path"], []),
        "english_latest_commit_hash": result.get("english_latest_commit_hash"),
        "ref_english_commit_hash": result.get("ref_english_commit_hash"),
        "ref_english_commit_date": result.get("ref_english_commit_date"),
    }


def save_matrix_files(
    matrix_data: dict[str, dict[str, Any]], output_dir: str = "data"
) -> None:
    """Save matrix data to category-specific JSON files."""
    matrix_dir = Path(output_dir) / "matrix"
    matrix_dir.mkdir(parents=True, exist_ok=True)

    for category, data in matrix_data.items():
        file_path = matrix_dir / f"{category}.json"
        with file_path.open("w", encoding="utf-8") as f:
            json.dump(
                convert_keys_to_camel_case(data),
                f,
                indent=2,
                default=serialize_datetime,
            )


def save_detail_files(
    results: dict[str, TranslationStatusResult],
    issues_by_file: dict[str, list[GitHubIssue]],
    prs_by_file: dict[str, list[GitHubPullRequest]],
    existing_urls: set[str],
    output_dir: str = "data",
) -> None:
    """Save detail data grouped by language and category.

    Args:
    ----
        results (dict[str, TranslationStatusResult]): The translation status results.
        prs_by_file (dict[str, list[GitHubPullRequest]]): A mapping of file paths to
                                                          their associated PRs.
        existing_urls (set[str]): A set of existing file paths to check against.
        output_dir (str): The directory where detail files will be saved.

    Returns:
    -------
        None: The function saves JSON files to the specified output directory.

    """
    details_dir = Path(output_dir) / "details"
    details_dir.mkdir(parents=True, exist_ok=True)

    details_by_language_category = defaultdict(lambda: defaultdict(dict))

    for result in results.values():
        language = result["language"]
        english_path = result["english_path"]
        original_category = result["category"]

        category_name = build_category_name(original_category, english_path)

        detail_data = create_detail_data(
            result,
            existing_urls,
            issues_by_file,
            prs_by_file,
        )
        details_by_language_category[language][category_name][english_path] = (
            detail_data
        )

    for language, categories in details_by_language_category.items():
        lang_dir = details_dir / language
        lang_dir.mkdir(exist_ok=True)

        for category, details in categories.items():
            file_path = lang_dir / f"{category}.json"
            with file_path.open("w", encoding="utf-8") as f:
                json.dump(
                    convert_keys_to_camel_case(details),
                    f,
                    indent=2,
                    default=serialize_datetime,
                )


def save_diff(matrix_data: dict[str, dict[str, Any]], output_dir: str = "data") -> None:
    """Save diff data grouped by category."""
    diff_dir = Path(output_dir) / "diff"
    diff_dir.mkdir(parents=True, exist_ok=True)

    # カテゴリーごとにデータを整理
    diffs_by_category = defaultdict(dict)

    for category, data in matrix_data.items():
        for article in data["articles"]:
            english_path = article["english_path"]

            for language, translation in article["translations"].items():
                status = translation.get("status")
                if status != "outdated":
                    continue

                ref_hash = translation.get("ref_english_commit_hash")
                latest_hash = translation.get("english_latest_commit_hash")

                if ref_hash and latest_hash and ref_hash != latest_hash:
                    diff_output = diff(ref_hash, latest_hash, english_path)

                    # Get actual translation file path from translation data
                    target_path = translation.get("target_path")
                    if not target_path:
                        # Fallback: simple replacement
                        target_path = english_path.replace(
                            "content/en/", f"content/{language}/"
                        )

                    # Use translation file path as key, organized by category
                    diffs_by_category[category][target_path] = {
                        "english_path": english_path,
                        "language": language,
                        "ref_english_commit_hash": ref_hash,
                        "english_latest_commit_hash": latest_hash,
                        "diff": diff_output,
                    }

    # カテゴリーごとにファイルを保存
    for category, diffs in diffs_by_category.items():
        if not diffs:  # 空の場合はスキップ
            continue

        file_path = diff_dir / f"{category}_diff.json"
        with file_path.open("w", encoding="utf-8") as f:
            converted_diffs = {
                file_path: convert_keys_to_camel_case(diff_data)
                for file_path, diff_data in diffs.items()
            }
            json.dump(
                converted_diffs,
                f,
                indent=2,
                default=serialize_datetime,
            )


def process_translation_results(
    results: dict[str, TranslationStatusResult],
    existing_urls: set[str],
    output_dir: str = "data",
) -> None:
    """Process translation results and save them to JSON files.

    Args:
    ----
        results (dict[str, TranslationStatusResult]): The translation status results.
        existing_urls (set[str]): A set of existing urls to check against.
        output_dir (str): The directory where output files will be saved.

    Returns:
    -------
        None: The function saves JSON files to the specified output directory.

    """
    results = dict(sorted(results.items(), key=lambda item: item[0].lower()))

    filtered_results = {
        target_path: result
        for target_path, result in results.items()
        if should_process(result)
    }

    # issues
    issues_by_file = get_issues_by_file()

    # prs
    prs_by_file = get_prs_by_file()

    # Create matrix data from results
    matrix_data = create_matrix_data(
        filtered_results,
        issues_by_file,
        prs_by_file,
        existing_urls,
    )

    logger.info("Saving diff files...")
    save_diff(matrix_data, output_dir)

    logger.info("Saving matrix files...")
    save_matrix_files(matrix_data, output_dir)

    # Save detailed translation results
    logger.info("Saving detail files...")
    save_detail_files(
        filtered_results,
        issues_by_file,
        prs_by_file,
        existing_urls,
        output_dir,
    )
