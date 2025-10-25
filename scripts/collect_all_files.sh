#!/bin/bash

set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REPO_PATH="${ROOT_DIR}/k8s-repo/website"
REPO_CONTENT_PATH="${REPO_PATH}/content"
OUTPUT_DIR="${ROOT_DIR}/data/master"

mkdir -p "${OUTPUT_DIR}"

LANGUAGES=("bn" "de" "en" "es" "fr" "hi" "it" "ja" "ko" "pl" "pt-br" "ru" "uk" "vi" "zh-cn")

log_info "Starting to collect files for each language..."

# Create a temporary file to store all file paths
temp_file="${OUTPUT_DIR}/all_files.csv.tmp" > "$temp_file"  # Clear the file

for lang in "${LANGUAGES[@]}"; do
  lang_dir="${REPO_CONTENT_PATH}/${lang}"
  
  if [ ! -d "$lang_dir" ]; then
    log_warn "Directory for language '$lang' not found at $lang_dir. Skipping."
    continue
  fi
  
  log_info "Processing language: $lang"

  cd "${REPO_PATH}"
  
  # Collect all file paths for the current language and append to the temporary file
  find "content/${lang}" -type f | sort >> "$temp_file"
  
  cd - > /dev/null
  
  lang_file_count=$(find "${REPO_CONTENT_PATH}/${lang}" -type f | wc -l)
  log_success "Found $lang_file_count files for $lang"
done

# Move the temporary file to the final file
mv "$temp_file" "${OUTPUT_DIR}/all_files.csv"

total_file_count=$(wc -l < "${OUTPUT_DIR}/all_files.csv")

log_success "File collection complete. Total $total_file_count files saved to ${OUTPUT_DIR}/all_files.csv"