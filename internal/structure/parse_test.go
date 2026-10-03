package structure

import (
	"reflect"
	"testing"
)

func set(keys ...string) map[string]struct{} {
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		m[k] = struct{}{}
	}
	return m
}

// Expected values were produced by upstream's parse_markdown() on the same input.
func TestParse(t *testing.T) {
	content := "---\r\ntitle: x\r\n---\r\n" +
		"## H {#Abc }\r\n" +
		"\r\n" +
		"text v1.28 v1.29,v1.30 kubev1.1 v1.2845 _v1.7\r\n" +
		"  ```yaml\r\n" +
		"  # not heading\r\n" +
		"  apiVersion: apps/v1\r\n" +
		"  kind: Deployment\r\n" +
		"  ```\r\n" +
		"- kind: \"Pod\"\r\n" +
		"{{< feature-state for_k8s_version=\"v1.31\" state=\"beta\" >}}\r\n" +
		"{{< feature-state feature_gate_name=\"SidecarContainers\" >}}\r\n"

	want := Features{
		VisibleLines: 10,
		H2:           1,
		H3:           0,
		CodeBlocks:   0, // indented fences toggle code state but are not counted
		Anchors:      set("abc"),
		BodyWords:    29,
		Versions: map[[2]int]struct{}{
			{1, 7}: {}, {1, 28}: {}, {1, 29}: {}, {1, 30}: {}, {1, 31}: {},
		},
		FeatureStateTokens: set("version:v1.31", "gate:SidecarContainers"),
		ApiKindTokens:      set("api:apps/v1", "kind:Deployment", "kind:Pod"),
	}

	got := Parse(content)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Parse() mismatch\n got: %+v\nwant: %+v", got, want)
	}
}

func TestParseStripsComments(t *testing.T) {
	content := "<!--\n## English {#english}\n```\ncode\n```\nfeature_gate_name=\"Foo\"\nSince v1.31.\n-->\n" +
		"## 标题 {#title}\n\n自 v1.28 起可用。\n"

	got := Parse(content)

	if got.H2 != 1 || got.CodeBlocks != 0 {
		t.Errorf("H2 = %d, CodeBlocks = %d, want 1 and 0", got.H2, got.CodeBlocks)
	}
	if !reflect.DeepEqual(got.Anchors, set("title")) {
		t.Errorf("Anchors = %v, want only title", got.Anchors)
	}
	if len(got.FeatureStateTokens) != 0 {
		t.Errorf("FeatureStateTokens = %v, want none", got.FeatureStateTokens)
	}
	// Upstream does not strip comments for versions, so v1.31 is still collected.
	wantVersions := map[[2]int]struct{}{{1, 28}: {}, {1, 31}: {}}
	if !reflect.DeepEqual(got.Versions, wantVersions) {
		t.Errorf("Versions = %v, want %v", got.Versions, wantVersions)
	}
}

func TestParseEmpty(t *testing.T) {
	got := Parse("")
	if got.VisibleLines != 0 || got.BodyWords != 0 || got.H2 != 0 || got.CodeBlocks != 0 {
		t.Errorf("Parse(\"\") = %+v, want zero counts", got)
	}
}
