package structure

import "regexp"

var (
	versionExpr             = regexp.MustCompile(`(?:^|[^A-Za-z0-9])v([0-9]+)\.([0-9]{1,3})`)
	featureStateVersionExpr = regexp.MustCompile(`for_k8s_version="(v\d+\.\d+)"`)
	featureStateGateExpr    = regexp.MustCompile(`feature_gate_name="([A-Za-z_][A-Za-z0-9_]*)"`)
	anchorExpr              = regexp.MustCompile(`\{#([^}]+)\}`)
	frontMatterExpr         = regexp.MustCompile(`(?s)^---\s*\n.*?\n---\s*\n`)
	codeExpr                = regexp.MustCompile("(?s)```.*?```")
	commentExpr             = regexp.MustCompile(`(?s)<!--.*?-->`)
	inlineExpr              = regexp.MustCompile("`[^`\n]+`")
	bodyWordExpr            = regexp.MustCompile(`[\p{L}\p{N}]{2,}`)
	paragraphSplitExpr      = regexp.MustCompile(`\n{2,}`)
	apiVersionLineExpr      = regexp.MustCompile(`(?m)^[ \t]*-?[ \t]*apiVersion:\s*\"?([A-Za-z0-9./_-]+)\"?\s*$`)
	kindLineExpr            = regexp.MustCompile(`(?m)^[ \t]*-?[ \t]*kind:\s*\"?([A-Z][A-Za-z0-9]+)\"?\s*$`)
)
