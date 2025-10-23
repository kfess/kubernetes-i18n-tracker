#!/bin/bash

set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() {
  echo -e "${BLUE}[INFO]${NC} $*" >&2
}

log_warn() {
  echo -e "${YELLOW}[WARN]${NC} $*" >&2
}

log_error() {
  echo -e "${RED}[ERROR]${NC} $*" >&2
}

log_success() {
  echo -e "${GREEN}[OK]${NC} $*" >&2
}

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
REPO_PATH="${ROOT_DIR}/k8s-repo/website"
OUTPUT_DIR="${ROOT_DIR}/data/master"
OUTPUT_FILE="${OUTPUT_DIR}/git_history_refactor.jsonl"

mkdir -p "${OUTPUT_DIR}"

if [ ! -d "$REPO_PATH" ]; then
  log_error "Repository not found: $REPO_PATH"
  exit 1
fi

cd "$REPO_PATH"

log_info "Updating the repository to the latest version..."
git pull origin main

log_info "Fetching git history for all content files..."

# すべての content/ 配下のファイルを取得
git ls-files 'content/**/*.md' 'content/**/*.html' | while read -r file; do
  # 各ファイルの履歴を取得（main ブランチに限定、リネーム追跡あり）
  git log main --follow --pretty=format:'%H%x1F%an%x1F%ad%x1F%s' --date=iso --numstat -M -- "$file"
  echo ""
done | \
  awk '
    BEGIN {
      RS="";
      FS="\n";
      ORS="";
    }
    {
      if (NF == 0) next;

      split($1, meta, "\x1F");
      hash = meta[1];
      author = meta[2];
      date = meta[3];
      message = meta[4];

      for (i = 2; i <= NF; i++) {
        if ($i ~ /^[0-9\-]/ && $i ~ /content/) {
          split($i, stats, "\t");
          insertions = stats[1];
          deletions = stats[2];
          filepath = stats[3];

          if (insertions == "-") {
            insertions_val = "null";
          } else {
            insertions_val = insertions;
          }

          if (deletions == "-") {
            deletions_val = "null";
          } else {
            deletions_val = deletions;
          }

          # リネーム検出（{old => new} 形式）
          if (filepath ~ /{.*=>.*}/) {
            old_path = parse_rename_path(filepath, "old");
            new_path = parse_rename_path(filepath, "new");
          } else {
            old_path = "";
            new_path = clean_path(filepath);
          }

          json_hash = "\"" hash "\"";
          json_author = "\"" escape_json(author) "\"";
          json_date = "\"" date "\"";
          json_message = "\"" escape_json(message) "\"";

          file_json = "{\"path\":\"" escape_json(new_path) "\",\"insertions\":" insertions_val ",\"deletions\":" deletions_val;

          if (old_path != "") {
            file_json = file_json ",\"old_path\":\"" escape_json(old_path) "\"";
          }

          file_json = file_json "}";

          printf "{\"hash\":%s,\"author\":%s,\"date\":%s,\"message\":%s,\"file\":%s}\n", json_hash, json_author, json_date, json_message, file_json;
        }
      }
    }

    function parse_rename_path(path, type) {
      brace_start = index(path, "{");
      brace_end = index(path, "}");
      
      if (brace_start == 0 || brace_end == 0) {
        return clean_path(path);
      }
      
      prefix = substr(path, 1, brace_start - 1);
      brace_content = substr(path, brace_start + 1, brace_end - brace_start - 1);
      suffix = substr(path, brace_end + 1);
      
      arrow_pos = index(brace_content, " => ");
      if (arrow_pos == 0) {
        return clean_path(path);
      }
      
      old_part = substr(brace_content, 1, arrow_pos - 1);
      new_part = substr(brace_content, arrow_pos + 4);
      
      gsub(/^ +| +$/, "", old_part);
      gsub(/^ +| +$/, "", new_part);
      
      if (type == "old") {
        return clean_path(prefix old_part suffix);
      } else {
        return clean_path(prefix new_part suffix);
      }
    }

    function clean_path(path) {
      if (path ~ /^".*"$/) {
        path = substr(path, 2, length(path) - 2);
      }
      return path;
    }

    function escape_json(str) {
      gsub(/\\/, "\\\\", str);
      gsub(/"/, "\\\"", str);
      gsub(/\r/, "\\r", str);
      gsub(/\n/, "\\n", str);
      gsub(/\t/, "\\t", str);
      gsub(/\b/, "\\b", str);
      gsub(/\f/, "\\f", str);
      gsub(/\007/, "\\u0007", str);
      gsub(/\013/, "\\u000B", str);

      return str;
    }
  ' > "$OUTPUT_FILE"

log_success "JSONL file created: $OUTPUT_FILE"

# 統計情報を表示
TOTAL_LINES=$(wc -l < "$OUTPUT_FILE")
UNIQUE_FILES=$(jq -r '.file.path' "$OUTPUT_FILE" 2>/dev/null | sort -u | wc -l)
UNIQUE_COMMITS=$(jq -r '.hash' "$OUTPUT_FILE" 2>/dev/null | sort -u | wc -l)

log_info "Total records: $TOTAL_LINES"
log_info "Unique files: $UNIQUE_FILES"
log_info "Unique commits: $UNIQUE_COMMITS"

