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
OUTPUT_FILE="${OUTPUT_DIR}/git_history.jsonl"
CACHE_DIR="${ROOT_DIR}/data/cache"
LAST_COMMIT_FILE="${CACHE_DIR}/last_commit.txt"

mkdir -p "${OUTPUT_DIR}"
mkdir -p "${CACHE_DIR}"

fetch_history_jsonl() {
  local start_commit=$1
  local end_commit=$2
  local output_file=$3

  local commit_range=""
  if [ -n "$start_commit" ] && [ -n "$end_commit" ]; then
    commit_range="${start_commit}..${end_commit}"
  fi

# find_last_commit() {
#   local merge_commit=$1
#   local file=$2
  
#   # ^1 と ^2 の両方から検索
#   local from_first=$(git log --pretty=format:'%H' "${merge_commit}^1" ^"${merge_commit}^2" -- "$file" | head -n 1)
#   local from_second=$(git log --pretty=format:'%H' "${merge_commit}^2" ^"${merge_commit}^1" -- "$file" | head -n 1)
  
#   # 両方が空の場合
#   if [ -z "$from_first" ] && [ -z "$from_second" ]; then
#     echo ""
#     return
#   fi
  
#   # どちらか一方のみの場合
#   if [ -z "$from_first" ]; then
#     local last="$from_second"
#   elif [ -z "$from_second" ]; then
#     local last="$from_first"
#   else
#     # 両方見つかった場合、より新しい方（コミット日時）を選択
#     local date_first=$(git log --pretty=format:'%at' -n 1 "$from_first")
#     local date_second=$(git log --pretty=format:'%at' -n 1 "$from_second")
    
#     if [ "$date_second" -gt "$date_first" ]; then
#       local last="$from_second"
#     else
#       local last="$from_first"
#     fi
#   fi
  
#   # マージコミットの場合は再帰
#   local parent_count=$(git rev-list --parents -n 1 "$last" | wc -w)
  
#   if [ "$parent_count" -gt 2 ]; then
#     find_last_commit "$last" "$file"
#   else
#     echo "$last"
#   fi
# }

find_last_commit() {
  local merge_commit=$1
  local file=$2
  local original_main_parent=$3  # 元のマージコミットのmain側の親
  
  # 初回呼び出しの場合、original_main_parentを設定
  if [ -z "$original_main_parent" ]; then
    original_main_parent="${merge_commit}^1"
  fi
  
  # ^1 と ^2 の両方から検索
  local from_first=$(git log --pretty=format:'%H' "${merge_commit}^1" ^"${merge_commit}^2" -- "$file" | head -n 1)
  local from_second=$(git log --pretty=format:'%H' "${merge_commit}^2" ^"${merge_commit}^1" -- "$file" | head -n 1)
  
  # 両方が空の場合
  if [ -z "$from_first" ] && [ -z "$from_second" ]; then
    echo ""
    return
  fi
  
  # mainブランチの第1親 (元のマージコミットのものを使用)
  local main_parent="$original_main_parent"
  
  # feature branchのコミットを選択
  local last=""
  
  if [ -n "$from_first" ] && [ -n "$from_second" ]; then
    # 両方見つかった場合、mainに含まれない方を選択
    if git merge-base --is-ancestor "$from_first" "$main_parent" 2>/dev/null; then
      # from_first はmainに含まれる → from_second を採用
      last="$from_second"
    elif git merge-base --is-ancestor "$from_second" "$main_parent" 2>/dev/null; then
      # from_second はmainに含まれる → from_first を採用
      last="$from_first"
    else
      # 両方ともmainに含まれない → より新しい方を選択
      local date_first=$(git log --pretty=format:'%at' -n 1 "$from_first")
      local date_second=$(git log --pretty=format:'%at' -n 1 "$from_second")
      
      if [ "$date_second" -gt "$date_first" ]; then
        last="$from_second"
      else
        last="$from_first"
      fi
    fi
  elif [ -n "$from_first" ]; then
    # from_first のみ存在
    if ! git merge-base --is-ancestor "$from_first" "$main_parent" 2>/dev/null; then
      last="$from_first"
    fi
  elif [ -n "$from_second" ]; then
    # from_second のみ存在
    if ! git merge-base --is-ancestor "$from_second" "$main_parent" 2>/dev/null; then
      last="$from_second"
    fi
  fi
  
  # feature branchのコミットが見つからなかった場合
  if [ -z "$last" ]; then
    echo ""
    return
  fi
  
  # マージコミットの場合は再帰
  local parent_count=$(git rev-list --parents -n 1 "$last" | wc -w)
  
  if [ "$parent_count" -gt 2 ]; then
    find_last_commit "$last" "$file" "$original_main_parent"
  else
    echo "$last"
  fi
}


git log --first-parent main --pretty=format:'%H' -- content/ | \
while read -r commit_hash; do
  parent_count=$(git rev-list --parents -n 1 "$commit_hash" | wc -w)

  if [ "$parent_count" -gt 2 ]; then
    merge_subject=$(git log --pretty=format:'%s' -n 1 "$commit_hash")
    
    diff_output=$(git diff --numstat "${commit_hash}^1" "${commit_hash}" -- content/)
    
    while IFS=$'\t' read -r added removed file; do
      # last_commit=$(find_last_commit "$commit_hash" "$file")
      last_commit=$(find_last_commit "$commit_hash" "$file" "")


      if [ -n "$last_commit" ]; then
        git show --pretty=format:'%H%x1F%an%x1F%ad%x1F'"${merge_subject}" -s --date=iso "$last_commit"
        echo ""
        echo -e "$added\t$removed\t$file"
      else
        git show --pretty=format:'%H%x1F%an%x1F%ad%x1F%s' -s --date=iso "$commit_hash"
        echo ""
        echo -e "$added\t$removed\t$file"
      fi
      echo ""
    done <<< "$diff_output"
  else
    git show --pretty=format:'%H%x1F%an%x1F%ad%x1F%s' --numstat --date=iso "$commit_hash" -- content/
    echo ""
  fi
done | tee /tmp/raw_output.txt | \
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
  ' >> "$output_file"
}

if [ ! -d "$REPO_PATH" ]; then
  log_info "Clone the Kubernetes repository..."
  mkdir -p "$(dirname "$REPO_PATH")"
  git clone -b main --single-branch https://github.com/kubernetes/website.git "$REPO_PATH"
fi

cd "$REPO_PATH"

log_info "Updating the repository to the latest version..."
git pull origin main

CURRENT_HEAD=$(git rev-parse HEAD)

if [ -f "$LAST_COMMIT_FILE" ]; then
  LAST_COMMIT=$(cat "$LAST_COMMIT_FILE")
  log_info "Fetching new commits since last processed commit: ${LAST_COMMIT} → ${CURRENT_HEAD}"

  if git merge-base --is-ancestor "$LAST_COMMIT" "$CURRENT_HEAD"; then
    fetch_history_jsonl "$LAST_COMMIT" "$CURRENT_HEAD" "$OUTPUT_FILE"
  else
    log_warn "Previous commit not found in history. Fetching full history."
    fetch_history_jsonl "" "" "$OUTPUT_FILE"
  fi
else
  log_info "First run or no previous commit information. Fetching full history."
  fetch_history_jsonl "" "" "$OUTPUT_FILE"
fi

echo "$CURRENT_HEAD" > "$LAST_COMMIT_FILE"
log_success "JSONL file created/updated: $OUTPUT_FILE"
log_success "Last processed commit hash: $CURRENT_HEAD"