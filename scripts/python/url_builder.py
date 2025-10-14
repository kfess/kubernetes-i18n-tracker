import re
from pathlib import Path, PurePosixPath

import yaml

ROOT_DIR = Path(__file__).resolve().parent.parent.parent
KUBERNETES_DIR = ROOT_DIR / "k8s-repo" / "website"


def build_url(  # noqa: PLR0911, C901
    english_path: str,
    language: str,
    existing_urls: set[str],
    base_url: str = "https://kubernetes.io",
) -> str | None:
    """Build URL for a given English path based on language and existing paths.

    Args:
    ----
        english_path (str): Path to the English markdown file.
        language (str): Language code (e.g., 'en', 'fr').
        existing_urls (set[str]): Set of existing urls to check against.
        base_url (str): Base URL for the site.

    Returns:
    -------
        str | None: The constructed URL or None if the path is invalid or not found.

    Example:
    -------
        content/en/docs/api/reference.md -> https://kubernetes.io/docs/api/reference/
        content/en/docs/api/reference.md -> https://kubernetes.io/ja/docs/api/reference/
        content/en/blog/_posts/2024-10-02-xxxx.md -> https://kubernetes.io/blog/2024/10/02/xxxx/
        content/en/blog/concepts/_index.md -> https://kubernetes.io/blog/concepts/

    """
    path = PurePosixPath(english_path)
    if not str(path).startswith("content/en/") or path.suffix not in (".md", ".html"):
        return None

    try:
        relative_path = path.relative_to("content/en")
        parts = relative_path.parts
    except ValueError:
        return None

    if len(parts) < 2:
        return None

    if not is_public_url(english_path):
        return None

    # Remove _index files
    if parts[-1] in ("_index.md", "_index.html"):
        parts = parts[:-1]

    lang_prefix = "" if language == "en" else f"{language}/"
    category = parts[0]

    if category == "docs":
        if len(parts) >= 3 and parts[1] == "reference" and parts[2] == "glossary":
            return _build_glossary_url(english_path, lang_prefix, base_url)
        elif len(parts) >= 3 and parts[1] == "contribute" and parts[2] == "blog":
            return _build_contribute_blog_url(
                english_path, parts, lang_prefix, base_url
            )

        doc_path = "/".join(parts[1:]).removesuffix(".md")
        url = f"{base_url}/{lang_prefix}docs/{doc_path}/"
        if url in existing_urls:
            return url
        elif url.lower() in existing_urls:
            return url.lower()
        return None

    elif category == "blog":
        return _build_blog_url(
            english_path, parts, existing_urls, lang_prefix, base_url
        )
    elif category == "includes":
        return None
    elif category == "case-studies":
        if len(parts) == 1:
            url = f"{base_url}/{lang_prefix}case-studies/"
            if url in existing_urls:
                return url

        case_study_path = "/".join(parts[1:-1])
        url = f"{base_url}/{lang_prefix}case-studies/{case_study_path}/"
        if url in existing_urls:
            return url
        elif url.lower() in existing_urls:
            return url.lower()
        return None

    # Other categories
    other_path = "/".join(parts).removesuffix(".md")
    url = f"{base_url}/{lang_prefix}{other_path}/"
    if url in existing_urls:
        return url
    elif url.lower() in existing_urls:
        return url.lower()

    return None


def _build_contribute_blog_url(
    file_path: str,
    parts: tuple,
    lang_prefix: str,
    base_url: str = "https://kubernetes.io",
) -> str | None:
    """Build URL for a blog post in the contribute section.

    Args:
    ----
        file_path (str): Path to the English markdown file.
        parts (tuple): Parts of the path split by '/'.
        lang_prefix (str): Language prefix for the URL.
        base_url (str): Base URL for the site.

    Returns:
    -------
        str | None: The constructed URL or None if the path is invalid or not found.

    """
    # has no extension (e.g., _index.md)
    if file_path.endswith(("_index.md", "_index.html")):
        return f"{base_url}/{lang_prefix}docs/contribute/blog/"

    # Try to get front matter
    front_matter = _parse_front_matter(file_path)

    # Priority 1: slug
    if front_matter and "slug" in front_matter:
        slug = front_matter["slug"]
        return f"{base_url}/{lang_prefix}docs/contribute/blog/{slug}/"

    doc_path = "/".join(parts[1:]).removesuffix(".md")
    return f"{base_url}/{lang_prefix}docs/{doc_path}/"


def _build_glossary_url(file_path: str, lang_prefix: str, base_url: str) -> str | None:
    """Build glossary URL based on file path.

    Args:
    ----
        file_path (str): Path to the markdown file.
        parts (tuple): Parts of the path split by '/'.
        lang_prefix (str): Language prefix for the URL.
        base_url (str): Base URL for the site.

    Returns:
    -------
        str | None: The URL or None if no valid URL found.

    """

    def parse_full_link(file_path: str) -> str | None:
        """Parse the full link from the file path."""
        try:
            content = Path(KUBERNETES_DIR / file_path).read_text(encoding="utf-8")
            match = re.match(r"^---\s*\n(.*?)\n---\s*\n", content, re.DOTALL)
            if not match:
                return None

            fm_content = match.group(1)

            try:
                parsed = yaml.safe_load(fm_content)
                if not isinstance(parsed, dict):
                    return None
            except yaml.YAMLError:
                return None

            url = parsed.get("full_link")
            if url:
                return str(url).strip()

            return None

        except (FileNotFoundError, UnicodeDecodeError):
            return None

    full_link = parse_full_link(file_path)

    if not full_link:
        return None

    # For example, full_link starts with "https://"
    if not full_link.startswith("/"):
        return None

    url = f"{base_url}/{lang_prefix}{full_link.strip('/')}/"

    # If full_link is a fragment (starts with #), return the URL without trailing slash
    if "#" in full_link:
        return url.rstrip("/")

    return url


def _text_to_slug(text: str) -> str:
    """Convert text to URL-friendly slug by removing symbols & replacing spaces with -.

    Args:
        text: Input text to convert

    Returns:
        Converted slug

    """
    cleaned = re.sub(r"[^a-zA-Z0-9\s\.\/\-]", "", text)
    hyphenated = re.sub(r"\s+", "-", cleaned).lower()
    normalized = re.sub(r"-+", "-", hyphenated)
    return re.sub(r"^-+|-+$", "", normalized)


def _build_blog_url(  # noqa: PLR0911, PLR0912, C901
    file_path: str,
    parts: tuple,
    existing_urls: set,
    lang_prefix: str,
    base_url: str,
) -> str | None:
    """Build blog URL with Hugo priority: slug+date > url > filename.

    Args:
    ----
        file_path (str): Path to the markdown file.
        partsPri (tuple): Parts of the path split by '/'.
        existing_urls (set): Set of existing URLs to check against.
        lang_prefix (str): Language prefix for the URL.
        base_url (str): Base URL for the site.

    Returns:
    -------
        str | None: The URL or None if no valid URL found.

    """
    # Try to get front matter
    front_matter = _parse_front_matter(file_path)

    # Priority 1: slug + date
    if front_matter and "slug" in front_matter and "date" in front_matter:
        slug = front_matter["slug"]
        date_match = re.match(r"(\d{4})-(\d{2})-(\d{2})", front_matter["date"])
        if date_match:
            year, month, day = date_match.groups()
            if day == "00":
                url = f"{base_url}/{lang_prefix}blog/{year}/{month}/{slug}/"
            else:
                url = f"{base_url}/{lang_prefix}blog/{year}/{month}/{day}/{slug}/"
            if url in existing_urls:
                return url
            elif url.lower() in existing_urls:
                return url.lower()

    # Priority 2: explicit url
    if front_matter and "url" in front_matter:
        url = f"{base_url}/{lang_prefix}{front_matter['url'].strip('/')}/"
        if url in existing_urls:
            return url
        elif url.lower() in existing_urls:
            return url.lower()

    # Priority 3: filename
    if len(parts) >= 3 and parts[1] == "_posts":
        filename = Path(parts[2]).stem
        date_match = re.match(r"(\d{4})-(\d{2})-(\d{2})-(.+)", filename)
        if date_match:
            year, month, day, title = date_match.groups()
            if day == "00":
                url = f"{base_url}/{lang_prefix}blog/{year}/{month}/{title}/"
            else:
                url = f"{base_url}/{lang_prefix}blog/{year}/{month}/{day}/{title}/"
        else:
            url = f"{base_url}/{lang_prefix}blog/{filename}/"
        if url in existing_urls:
            return url
        elif url.lower() in existing_urls:
            return url.lower()

    # Priority 4: blog title with date
    if front_matter and "title" in front_matter and "date" in front_matter:
        title = front_matter["title"]
        title = _text_to_slug(title)

        date_match = re.match(r"(\d{4})-(\d{2})-(\d{2})", front_matter["date"])
        if date_match:
            year, month, day = date_match.groups()
            if day == "00":
                url = f"{base_url}/{lang_prefix}blog/{year}/{month}/{title}/"
            else:
                url = f"{base_url}/{lang_prefix}blog/{year}/{month}/{day}/{title}/"
            if url in existing_urls:
                return url
            elif url.lower() in existing_urls:
                return url.lower()
    # Blog category
    blog_path = "/".join(parts[1:]).removesuffix(".md")
    url = (
        f"{base_url}/{lang_prefix}blog/{blog_path}/"
        if blog_path
        else f"{base_url}/{lang_prefix}blog/"
    )
    if url in existing_urls:
        return url
    elif url.lower() in existing_urls:
        return url.lower()

    return None


def _parse_front_matter(file_path: str) -> dict | None:
    """Parse front matter fields: url, slug, date."""
    try:
        content = Path(KUBERNETES_DIR / file_path).read_text(encoding="utf-8")
        # some file has leading blank lines
        # e.g., content/en/blog/_posts/2019-08-30-announcing-etcd-3.4.md
        content = content.lstrip("\n\r\t ")
        match = re.match(r"^---\s*\n(.*?)\n---\s*\n", content, re.DOTALL)
        if not match:
            return None

        fm_content = match.group(1)
        result = {}

        for field in ["url", "slug", "date", "title"]:
            field_match = re.search(
                rf'^{field}:\s*["\']?([^"\n]*?)["\']?\s*$', fm_content, re.MULTILINE
            )
            if field_match:
                result[field] = field_match.group(1).strip()

        return result if result else None
    except (FileNotFoundError, UnicodeDecodeError):
        return None


def is_public_url(file_path: str) -> bool:
    """Check if a file will have a public URL based on Hugo front matter settings.

    Returns False if:
    - _build.render: never

    Args:
        file_path: Path to the markdown file relative to KUBERNETES_DIR

    Returns:
        bool: True if the page will be publicly accessible, False otherwise

    """
    try:
        content = Path(KUBERNETES_DIR / file_path).read_text(encoding="utf-8")
        match = re.match(r"^---\s*\n(.*?)\n---\s*\n", content, re.DOTALL)
        if not match:
            return True

        fm_content = match.group(1)

        try:
            parsed = yaml.safe_load(fm_content)
            if not isinstance(parsed, dict):
                return True
        except yaml.YAMLError:
            return True

        build_settings = parsed.get("_build", {})
        if isinstance(build_settings, dict):
            render_value = build_settings.get("render")
            if render_value in ("never", False):
                return False

        return True

    except (FileNotFoundError, UnicodeDecodeError):
        return False
