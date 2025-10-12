import json
import re
import xml.etree.ElementTree as ET
from pathlib import Path

import requests
from const import LANGUAGE_CODES
from exporter import process_translation_results
from history import GitFileHistoryTracker
from log import logger
from models import GitCommitDict
from translation_status import TranslationStatusTracker

ROOT_DIR = Path(__file__).resolve().parent.parent.parent
INPUT_FILE = ROOT_DIR / "data" / "master" / "git_history.jsonl"
OUTPUT_DIR = ROOT_DIR / "data" / "output"


def sanitize(line: str) -> str:
    """Sanitize a string by replacing backslashes with double backslashes.

    Args:
    ----
        line (str): The input string to sanitize.

    Returns:
    -------
        str: The sanitized string with backslashes replaced.

    """
    return re.compile(r'\\(?![\\bfnrt"/])').sub(r"\\\\", line)


def transform_records(records: list[dict]) -> list[dict]:
    """Transform a list of Git commit records into a structured format."""
    result = []
    grouped = {}

    for record in records:
        key = (record["hash"], record["author"], record["date"], record["message"])
        if key not in grouped:
            grouped[key] = []
        grouped[key].append(record["file"])

    for (hash_val, author, date, message), files in grouped.items():
        total_insertions = sum(
            f["insertions"] for f in files if f["insertions"] is not None
        )
        total_deletions = sum(
            f["deletions"] for f in files if f["deletions"] is not None
        )

        result.append(
            {
                "hash": hash_val,
                "author": author,
                "date": date,
                "message": message,
                "files": files,
                "summary": {
                    "total_files": len(files),
                    "total_insertions": total_insertions,
                    "total_deletions": total_deletions,
                    "total_changes": total_insertions + total_deletions,
                },
            }
        )

    return result


def load_json_records(filepath: Path | str) -> list[GitCommitDict]:
    """Load JSON records from a file, handling potential formatting issues.

    Args:
    ----
        filepath (Path | str): The path to the JSONL file.

    Returns:
    -------
        list[dict[str, Any]]: A list of JSON objects loaded from the file.

    Raises:
    ------
        FileNotFoundError: If the file doesn't exist.

    """
    records: list[GitCommitDict] = []
    path = Path(filepath)

    if not path.exists():
        msg = f"File not found: {path}"
        raise FileNotFoundError(msg)

    logger.info("Start Loading JSON records from: %s", path)

    with path.open(encoding="utf-8") as f:
        for line_no, line in enumerate(f, 1):
            original = line.strip()
            if not original:
                continue
            try:
                obj = json.loads(original)
            except json.JSONDecodeError:
                fixed = sanitize(original)
                try:
                    logger.warning(
                        "Fixed line %d - Original: %s => Fixed: %s",
                        line_no,
                        original,
                        fixed,
                    )
                    obj = json.loads(fixed)
                except json.JSONDecodeError:
                    logger.exception("Failed to decode line %d", line_no)
                    continue

            records.append(obj)

    records = transform_records(records)

    logger.info("Successfully loaded %d records from %s", len(records), INPUT_FILE)
    return records


def load_existing_paths() -> set[str]:
    """Load existing file paths from text files.

    Returns
    -------
        set[str]: A set of existing file paths from the JSONL file.

    """
    target_file = ROOT_DIR / "data" / "master" / "all_files.csv"
    with Path.open(target_file, "r", encoding="utf-8") as f:
        return {line.strip() for line in f if line.strip()}


def load_existing_urls() -> set[str]:
    """Load existing URLs from sitemap.xml files.

    Returns
    -------
        set[str]: A set of existing URLs from sitemap files.

    Comments:
    --------
        This function fetches sitemap.xml files for each language and extracts URLs.
        It handles both all localized sitemaps.
        If a sitemap cannot be loaded, it logs an error but continues processing others.

    """
    urls: set[str] = set()

    for lang in LANGUAGE_CODES:
        sitemap_url = f"https://kubernetes.io/{lang}/sitemap.xml"
        try:
            response = requests.get(sitemap_url, timeout=10)
            response.raise_for_status()

            root = ET.fromstring(response.text)  # noqa: S314
            for url_elem in root.findall(
                ".//{http://www.sitemaps.org/schemas/sitemap/0.9}loc"
            ):
                urls.add(url_elem.text.strip())

        except Exception:
            logger.exception("Error loading sitemap for %s", lang)

    return urls


def main() -> None:
    """Load JSONL file and save translation results to output directory."""
    try:
        records = load_json_records(INPUT_FILE)
    except FileNotFoundError:
        logger.exception("File not found: %s")
        return []
    except Exception:
        logger.exception("An unexpected error occurred: %s")
        return []

    existing_paths = load_existing_paths()
    existing_urls = load_existing_urls()
    file_history_tracker = GitFileHistoryTracker(
        commits=records, current_files=existing_paths
    )
    translation_tracker = TranslationStatusTracker(
        file_history_tracker=file_history_tracker,
        existing_paths=existing_paths,
    )
    status_result = translation_tracker.analyze()
    process_translation_results(status_result, existing_urls, output_dir=OUTPUT_DIR)

    return None


if __name__ == "__main__":
    main()
