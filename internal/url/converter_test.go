package url

import (
	"context"
	"testing"

	"github.com/kfess/kubernetes-i18n-tracker/internal/language"
	"github.com/kfess/kubernetes-i18n-tracker/internal/path"
)

func TestConverter_DocsConcepts(t *testing.T) {
	config := &Config{
		BaseUrl: "https://kubernetes.io",
		ExistingUrls: map[string]bool{
			"https://kubernetes.io/docs/concepts/":                                                          true,
			"https://kubernetes.io/docs/concepts/architecture/":                                             true,
			"https://kubernetes.io/docs/concepts/architecture/cloud-controller/":                            true,
			"https://kubernetes.io/docs/concepts/cluster-administration/admission-webhooks-good-practices/": true,
		},
		SupportedLangs: language.SupportedLanguages,
		SupportedExts:  []string{".md", ".html"},
		ValidSections:  []path.Category{path.Docs},
	}
	converter := NewConverter(*config)
	ctx := context.Background()

	tests := []struct {
		name    string
		path    string
		fm      *FrontMatter
		want    string
		wantErr bool
	}{
		// docs/concepts tests
		{
			name:    "Basic docs with _index.md",
			path:    "content/en/docs/concepts/_index.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/docs/concepts/",
			wantErr: false,
		},
		{
			name:    "Basic docs with _index.md",
			path:    "content/en/docs/concepts/architecture/_index.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/docs/concepts/architecture/",
			wantErr: false,
		},
		{
			name:    "Basic docs with named markdown file",
			path:    "content/en/docs/concepts/architecture/cloud-controller.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/docs/concepts/architecture/cloud-controller/",
			wantErr: false,
		},
		{
			name:    "Basic docs with named markdown file",
			path:    "content/en/docs/concepts/cluster-administration/admission-webhooks-good-practices.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/docs/concepts/cluster-administration/admission-webhooks-good-practices/",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.Convert(ctx, tt.path, tt.fm)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConverter_DocsTasks(t *testing.T) {
	config := &Config{
		BaseUrl: "https://kubernetes.io",
		ExistingUrls: map[string]bool{
			"https://kubernetes.io/docs/tasks/":                                                  true,
			"https://kubernetes.io/docs/tasks/manage-kubernetes-objects/":                        true,
			"https://kubernetes.io/docs/tasks/administer-cluster/change-pv-reclaim-policy/":      true,
			"https://kubernetes.io/docs/tasks/tools/included/optional-kubectl-configs-bash-mac/": true,
		},
		SupportedLangs: language.SupportedLanguages,
		SupportedExts:  []string{".md", ".html"},
		ValidSections:  []path.Category{path.Docs},
	}
	converter := NewConverter(*config)
	ctx := context.Background()

	tests := []struct {
		name    string
		path    string
		fm      *FrontMatter
		want    string
		wantErr bool
	}{
		{
			name:    "Basic docs with _index.md",
			path:    "content/en/docs/tasks/_index.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/docs/tasks/",
			wantErr: false,
		},
		{
			name:    "Basic docs with _index.md",
			path:    "content/en/docs/tasks/manage-kubernetes-objects/_index.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/docs/tasks/manage-kubernetes-objects/",
			wantErr: false,
		},
		{
			name:    "Basic docs with named markdown file",
			path:    "content/en/docs/tasks/administer-cluster/change-pv-reclaim-policy.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/docs/tasks/administer-cluster/change-pv-reclaim-policy/",
			wantErr: false,
		},
		{
			name:    "Basic docs with named markdown file",
			path:    "content/en/docs/tasks/tools/included/optional-kubectl-configs-bash-mac.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/docs/tasks/tools/included/optional-kubectl-configs-bash-mac/",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.Convert(ctx, tt.path, tt.fm)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConverter_DocsSetup(t *testing.T) {
	config := &Config{
		BaseUrl: "https://kubernetes.io",
		ExistingUrls: map[string]bool{
			"https://kubernetes.io/docs/setup/":                              true,
			"https://kubernetes.io/docs/setup/best-practices/":               true,
			"https://kubernetes.io/docs/setup/best-practices/certificates/":  true,
			"https://kubernetes.io/docs/setup/best-practices/cluster-large/": true,
		},
		SupportedLangs: language.SupportedLanguages,
		SupportedExts:  []string{".md", ".html"},
		ValidSections:  []path.Category{path.Docs},
	}
	converter := NewConverter(*config)
	ctx := context.Background()

	tests := []struct {
		name    string
		path    string
		fm      *FrontMatter
		want    string
		wantErr bool
	}{
		{
			name:    "Basic docs with _index.md",
			path:    "content/en/docs/setup/_index.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/docs/setup/",
			wantErr: false,
		},
		{
			name:    "Basic docs with _index.md",
			path:    "content/en/docs/setup/best-practices/_index.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/docs/setup/best-practices/",
			wantErr: false,
		},
		{
			name:    "Basic docs with named markdown file",
			path:    "content/en/docs/setup/best-practices/certificates.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/docs/setup/best-practices/certificates/",
			wantErr: false,
		},
		{
			name:    "Basic docs with named markdown file",
			path:    "content/en/docs/setup/best-practices/cluster-large.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/docs/setup/best-practices/cluster-large/",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.Convert(ctx, tt.path, tt.fm)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConverter_DocsReference(t *testing.T) {
	config := &Config{
		BaseUrl: "https://kubernetes.io",
		ExistingUrls: map[string]bool{
			"https://kubernetes.io/docs/reference/":                                              true,
			"https://kubernetes.io/docs/reference/access-authn-authz/":                           true,
			"https://kubernetes.io/docs/reference/access-authn-authz/abac/":                      true,
			"https://kubernetes.io/docs/reference/access-authn-authz/kubelet-tls-bootstrapping/": true,
		},
		SupportedLangs: language.SupportedLanguages,
		SupportedExts:  []string{".md", ".html"},
		ValidSections:  []path.Category{path.Docs},
	}
	converter := NewConverter(*config)
	ctx := context.Background()

	tests := []struct {
		name    string
		path    string
		fm      *FrontMatter
		want    string
		wantErr bool
	}{
		{
			name:    "Basic docs with _index.md",
			path:    "content/en/docs/reference/_index.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/docs/reference/",
			wantErr: false,
		},
		{
			name:    "Basic docs with _index.md",
			path:    "content/en/docs/reference/access-authn-authz/_index.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/docs/reference/access-authn-authz/",
			wantErr: false,
		},
		{
			name:    "Basic docs with named markdown file",
			path:    "content/en/docs/reference/access-authn-authz/abac.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/docs/reference/access-authn-authz/abac/",
			wantErr: false,
		},
		{
			name:    "Basic docs with named markdown file",
			path:    "content/en/docs/reference/access-authn-authz/kubelet-tls-bootstrapping.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/docs/reference/access-authn-authz/kubelet-tls-bootstrapping/",
			wantErr: false,
		},
		{
			name:    "index.md file without URL",
			path:    "content/en/docs/reference/command-line-tools-reference/feature-gates-removed/index.md",
			fm:      &FrontMatter{},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Basic reference docs without URL",
			path:    "content/en/docs/docs/reference/command-line-tools-reference/feature-gates/AdmissionWebhookMatchConditions.md",
			fm:      &FrontMatter{},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.Convert(ctx, tt.path, tt.fm)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConverter_DocsTutorials(t *testing.T) {
	config := &Config{
		BaseUrl: "https://kubernetes.io",
		ExistingUrls: map[string]bool{
			"https://kubernetes.io/docs/tutorials/":                                                           true,
			"https://kubernetes.io/docs/tutorials/configuration/configure-redis-using-configmap/":             true,
			"https://kubernetes.io/docs/tutorials/hello-minikube/":                                            true,
			"https://kubernetes.io/docs/tutorials/kubernetes-basics/create-cluster/cluster-interactive-gone/": true,
		},
		SupportedLangs: language.SupportedLanguages,
		SupportedExts:  []string{".md", ".html"},
		ValidSections:  []path.Category{path.Docs},
	}
	converter := NewConverter(*config)
	ctx := context.Background()

	tests := []struct {
		name    string
		path    string
		fm      *FrontMatter
		want    string
		wantErr bool
	}{
		{
			name:    "Tutorials index",
			path:    "content/en/docs/tutorials/_index.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/docs/tutorials/",
			wantErr: false,
		},
		{
			name:    "Nested tutorial with md",
			path:    "content/en/docs/tutorials/configuration/configure-redis-using-configmap.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/docs/tutorials/configuration/configure-redis-using-configmap/",
			wantErr: false,
		},
		{
			name:    "Top-level tutorial file",
			path:    "content/en/docs/tutorials/hello-minikube.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/docs/tutorials/hello-minikube/",
			wantErr: false,
		},
		{
			name:    "Deep nested HTML file",
			path:    "content/en/docs/tutorials/kubernetes-basics/create-cluster/cluster-interactive-gone.html",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/docs/tutorials/kubernetes-basics/create-cluster/cluster-interactive-gone/",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.Convert(ctx, tt.path, tt.fm)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConverter_DocsContribute(t *testing.T) {
	config := &Config{
		BaseUrl: "https://kubernetes.io",
		ExistingUrls: map[string]bool{
			"https://kubernetes.io/docs/contribute/":                           true,
			"https://kubernetes.io/docs/contribute/localization/":              true,
			"https://kubernetes.io/docs/contribute/generate-ref-docs/kubectl/": true,
		},
		SupportedLangs: language.SupportedLanguages,
		SupportedExts:  []string{".md", ".html"},
		ValidSections:  []path.Category{path.Docs},
	}
	converter := NewConverter(*config)
	ctx := context.Background()

	tests := []struct {
		name    string
		path    string
		fm      *FrontMatter
		want    string
		wantErr bool
	}{
		{
			name:    "Contribute index",
			path:    "content/en/docs/contribute/_index.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/docs/contribute/",
			wantErr: false,
		},
		{
			name:    "Top-level contribute page",
			path:    "content/en/docs/contribute/localization.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/docs/contribute/localization/",
			wantErr: false,
		},
		{
			name:    "Nested contribute page",
			path:    "content/en/docs/contribute/generate-ref-docs/kubectl.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/docs/contribute/generate-ref-docs/kubectl/",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.Convert(ctx, tt.path, tt.fm)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConverter_DocsHome(t *testing.T) {
	config := &Config{
		BaseUrl: "https://kubernetes.io",
		ExistingUrls: map[string]bool{
			"https://kubernetes.io/docs/home/":                        true,
			"https://kubernetes.io/ja/docs/home/":                     true,
			"https://kubernetes.io/docs/home/supported-doc-versions/": true,
		},
		SupportedLangs: language.SupportedLanguages,
		SupportedExts:  []string{".md", ".html"},
		ValidSections:  []path.Category{path.Docs},
	}
	converter := NewConverter(*config)
	ctx := context.Background()

	tests := []struct {
		name    string
		path    string
		fm      *FrontMatter
		want    string
		wantErr bool
	}{
		{
			name:    "Home index (root)",
			path:    "content/en/docs/home/_index.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/docs/home/",
			wantErr: false,
		},
		{
			name:    "Japanese home index",
			path:    "content/ja/docs/home/_index.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/ja/docs/home/",
			wantErr: false,
		},
		{
			name:    "Supported doc versions",
			path:    "content/en/docs/home/supported-doc-versions.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/docs/home/supported-doc-versions/",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.Convert(ctx, tt.path, tt.fm)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConverter_Blog(t *testing.T) {
	config := &Config{
		BaseUrl: "https://kubernetes.io",
		ExistingUrls: map[string]bool{
			// Blog index
			"https://kubernetes.io/blog/": true,

			// public blog posts
			"https://kubernetes.io/blog/2025/10/06/introducing-headlamp-plugin-for-karpenter/":                                                                        true,
			"https://kubernetes.io/blog/2025/09/08/kubernetes-v1-34-volume-attributes-class/":                                                                         true,
			"https://kubernetes.io/blog/2025/07/03/navigating-failures-in-pods-with-devices/":                                                                         true,
			"https://kubernetes.io/blog/2025/03/12/sig-apps-spotlight-2025/":                                                                                          true,
			"https://kubernetes.io/blog/2025/03/04/sig-etcd-spotlight/":                                                                                               true,
			"https://kubernetes.io/blog/2025/02/28/nftables-kube-proxy/":                                                                                              true,
			"https://kubernetes.io/blog/2025/02/14/cloud-controller-manager-chicken-egg-problem/":                                                                     true,
			"https://kubernetes.io/blog/2024/08/07/sig-api-machinery-spotlight-2024/":                                                                                 true,
			"https://kubernetes.io/blog/2023/03/17/upcoming-changes-in-kubernetes-v1-27/":                                                                             true,
			"https://kubernetes.io/blog/2020/06/30/sig-windows-spotlight-2020/":                                                                                       true,
			"https://kubernetes.io/blog/2021/12/15/kubernetes-1-23-prevent-persistentvolume-leaks-when-deleting-out-of-order/":                                        true,
			"https://kubernetes.io/blog/2021/12/01/contribution-containers-and-cricket-the-kubernetes-1.22-release-interview/":                                        true,
			"https://kubernetes.io/blog/2020/08/03/physics-politics-and-pull-requests-the-kubernetes-1.18-release-interview/":                                         true,
			"https://kubernetes.io/blog/2019/12/06/when-youre-in-the-release-team-youre-family-the-kubernetes-1.16-release-interview/":                                true,
			"https://kubernetes.io/blog/2019/04/24/hardware-accelerated-ssl/tls-termination-in-ingress-controllers-using-kubernetes-device-plugins-and-runtimeclass/": true,
			"https://kubernetes.io/blog/2019/02/28/automate-operations-on-your-cluster-with-operatorhub.io/":                                                          true,
			"https://kubernetes.io/blog/2018/10/10/kubernetes-v1.12-introducing-runtimeclass/":                                                                        true,
			"https://kubernetes.io/blog/2018/10/04/introducing-the-non-code-contributors-guide/":                                                                      true,
			"https://kubernetes.io/blog/2018/09/27/kubernetes-1.12-kubelet-tls-bootstrap-and-azure-virtual-machine-scale-sets-vmss-move-to-general-availability/":     true,
			"https://kubernetes.io/blog/2018/09/18/hands-on-with-linkerd-2.0/":                                                                                        true,
			"https://kubernetes.io/blog/2018/07/16/how-the-sausage-is-made-the-kubernetes-1.11-release-interview-from-the-kubernetes-podcast/":                        true,
			"https://kubernetes.io/blog/2018/04/13/local-persistent-volumes-beta/":                                                                                    true,
			"https://kubernetes.io/blog/2017/02/postgresql-clusters-kubernetes-statefulsets/":                                                                         true,
			"https://kubernetes.io/blog/2016/05/coreosfest2016-kubernetes-community/":                                                                                 true,
			"https://kubernetes.io/blog/2015/06/cluster-level-logging-with-kubernetes/":                                                                               true,
			"https://kubernetes.io/blog/2015/03/welcome-to-kubernetes-blog/":                                                                                          true,

			// Private blog posts (dummy url)
			"https://kubernetes.io/blog/2025/11/01/private-post/": true,
		},
		SupportedLangs: language.SupportedLanguages,
		SupportedExts:  []string{".md", ".html"},
		ValidSections:  []path.Category{path.Blog},
	}
	converter := NewConverter(*config)
	ctx := context.Background()

	tests := []struct {
		name    string
		path    string
		fm      *FrontMatter
		want    string
		wantErr bool
	}{
		{
			name:    "Blog index",
			path:    "content/en/blog/_index.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/blog/",
			wantErr: false,
		},
		{
			name: "Public blog post",
			path: "content/en/blog/_posts/2025-10-06-headlamp-karpenter-plugin/index.md",
			fm: &FrontMatter{
				Title: "Introducing Headlamp Plugin for Karpenter - Scaling and Visibility",
				Date:  "2025-10-06",
				Slug:  "introducing-headlamp-plugin-for-karpenter",
			},
			want:    "https://kubernetes.io/blog/2025/10/06/introducing-headlamp-plugin-for-karpenter/",
			wantErr: false,
		},
		{
			name: "Public blog post",
			path: "content/en/blog/_posts/2025-09-08-volume-attributes-class-ga/index.md",
			fm: &FrontMatter{
				Title: "Kubernetes v1.34: VolumeAttributesClass for Volume Modification GA",
				Date:  "2025-09-08 10:30:00 -0800",
				Slug:  "kubernetes-v1-34-volume-attributes-class",
			},
			want:    "https://kubernetes.io/blog/2025/09/08/kubernetes-v1-34-volume-attributes-class/",
			wantErr: false,
		},
		{
			name: "Public blog post",
			path: "content/en/blog/_posts/2025-07-03-devices-failure-handling/index.md",
			fm: &FrontMatter{
				Title: "Navigating Failures in Pods With Devices",
				Date:  "2025-07-03",
				Slug:  "navigating-failures-in-pods-with-devices",
			},
			want:    "https://kubernetes.io/blog/2025/07/03/navigating-failures-in-pods-with-devices/",
			wantErr: false,
		},
		{
			name: "Public blog post",
			path: "content/en/blog/_posts/2025-03-12-sig-apps-spotlight.md",
			fm: &FrontMatter{
				Title: "Spotlight on SIG Apps",
				Slug:  "sig-apps-spotlight-2025",
				Date:  "2025-03-12",
			},
			want:    "https://kubernetes.io/blog/2025/03/12/sig-apps-spotlight-2025/",
			wantErr: false,
		},
		{
			name: "Public blog post",
			path: "content/en/blog/_posts/2025-03-04-sig-etcd-spotlight/index.md",
			fm: &FrontMatter{
				Title: "Spotlight on SIG etcd",
				Slug:  "sig-etcd-spotlight",
				Date:  "2025-03-04",
			},
			want:    "https://kubernetes.io/blog/2025/03/04/sig-etcd-spotlight/",
			wantErr: false,
		},
		{
			name: "Public blog post",
			path: "content/en/blog/_posts/2025-02-28-nftables-kube-proxy/index.md",
			fm: &FrontMatter{
				Title: "NFTables mode for kube-proxy",
				Slug:  "nftables-kube-proxy",
				Date:  "2025-02-28",
			},
			want:    "https://kubernetes.io/blog/2025/02/28/nftables-kube-proxy/",
			wantErr: false,
		},
		{
			name: "Public blog post",
			path: "content/en/blog/_posts/2025-02-14-cloud-controller-manager-chicken-egg-problem/index.md",
			fm: &FrontMatter{
				Title: "The Cloud Controller Manager Chicken and Egg Problem",
				Slug:  "cloud-controller-manager-chicken-egg-problem",
				Date:  "2025-02-14",
			},
			want:    "https://kubernetes.io/blog/2025/02/14/cloud-controller-manager-chicken-egg-problem/",
			wantErr: false,
		},
		{
			name: "Public blog post",
			path: "content/en/blog/_posts/2023-03-17-kubernetes-1.27-deprecations-and-removals.md",
			fm: &FrontMatter{
				Title: "Kubernetes Removals and Major Changes In v1.27",
				Date:  "2023-03-17 14:00:00 +0000",
				Slug:  "upcoming-changes-in-kubernetes-v1-27",
			},
			want:    "https://kubernetes.io/blog/2023/03/17/upcoming-changes-in-kubernetes-v1-27/",
			wantErr: false,
		},
		{
			name: "Public blog post",
			path: "content/en/blog/_posts/2021-12-15-prevent-persistentvolume-leaks-when-deleting-out-of-order.md",
			fm: &FrontMatter{
				Title: "Kubernetes 1.23: Prevent PersistentVolume leaks when deleting out of order",
				Date:  "2021-12-15 10:00:00 -0800",
				Slug:  "kubernetes-1-23-prevent-persistentvolume-leaks-when-deleting-out-of-order",
			},
			want:    "https://kubernetes.io/blog/2021/12/15/kubernetes-1-23-prevent-persistentvolume-leaks-when-deleting-out-of-order/",
			wantErr: false,
		},
		{
			name: "Public blog post",
			path: "content/en/blog/_posts/2021-12-01-kubernetes-1.22-release-interview.md",
			fm: &FrontMatter{
				Title: "Contribution, containers and cricket: the Kubernetes 1.22 release interview",
				Date:  "2021-12-01",
			},
			want:    "https://kubernetes.io/blog/2021/12/01/contribution-containers-and-cricket-the-kubernetes-1.22-release-interview/",
			wantErr: false,
		},
		{
			name: "Public blog post",
			path: "content/en/blog/_posts/2020-08-03-kubernetes-1-18-release-interview.md",
			fm: &FrontMatter{
				Title: "Physics, politics and Pull Requests: the Kubernetes 1.18 release interview",
				Date:  "2020-08-03",
			},
			want:    "https://kubernetes.io/blog/2020/08/03/physics-politics-and-pull-requests-the-kubernetes-1.18-release-interview/",
			wantErr: false,
		},
		{
			name: "Public blog post",
			path: "content/en/blog/_posts/2020-06-30-SIG-Windows-Spotlight/index.md",
			fm: &FrontMatter{
				Title: "SIG-Windows Spotlight",
				Slug:  "sig-windows-spotlight-2020",
				Date:  "2020-06-30",
			},
			want:    "https://kubernetes.io/blog/2020/06/30/sig-windows-spotlight-2020/",
			wantErr: false,
		},
		{
			name: "Public blog post",
			path: "content/en/blog/_posts/2019-12-06-kubernetes-1-16-release-interview.md",
			fm: &FrontMatter{
				Title: "When you're in the release team, you're family: the Kubernetes 1.16 release interview",
				Date:  "2019-12-06",
			},
			want:    "https://kubernetes.io/blog/2019/12/06/when-youre-in-the-release-team-youre-family-the-kubernetes-1.16-release-interview/",
			wantErr: false,
		},
		{
			name: "Blog Title with special characters (slash)",
			path: "content/en/blog/_posts/2019-04-24-Hardware-Accelerated-SSLTLS-Termination-in-Ingress-Controllers-using-Kubernetes-Device-Plugins-and-RuntimeClass.md",
			fm: &FrontMatter{
				Title: "Hardware Accelerated SSL/TLS Termination in Ingress Controllers using Kubernetes Device Plugins and RuntimeClass",
				Date:  "2019-04-24",
			},
			want:    "https://kubernetes.io/blog/2019/04/24/hardware-accelerated-ssl/tls-termination-in-ingress-controllers-using-kubernetes-device-plugins-and-runtimeclass/",
			wantErr: false,
		},
		{
			name: "Public blog post",
			path: "content/en/blog/_posts/2019-02-28-automate-operations-on-your-cluster-with-operatorhub.md",
			fm: &FrontMatter{
				Title: "Automate Operations on your Cluster with OperatorHub.io",
				Date:  "2019-02-28",
			},
			want:    "https://kubernetes.io/blog/2019/02/28/automate-operations-on-your-cluster-with-operatorhub.io/",
			wantErr: false,
		},
		{
			name: "Public blog post",
			path: "content/en/blog/_posts/2018-10-10-runtimeclass.md",
			fm: &FrontMatter{
				Title: "Kubernetes v1.12: Introducing RuntimeClass",
				Date:  "2018-10-10",
			},
			want:    "https://kubernetes.io/blog/2018/10/10/kubernetes-v1.12-introducing-runtimeclass/",
			wantErr: false,
		},
		{
			name: "Public blog post",
			path: "content/en/blog/_posts/2018-10-04-non-code-contributors-guide.md",
			fm: &FrontMatter{
				Title: "Introducing the Non-Code Contributor’s Guide",
				Date:  "2018-10-04",
			},
			want:    "https://kubernetes.io/blog/2018/10/04/introducing-the-non-code-contributors-guide/",
			wantErr: false,
		},
		{
			name: "Public blog post",
			path: "content/en/blog/_posts/2018-09-27-kubernetes-1-12-release-announcement.md",
			fm: &FrontMatter{
				Title: "Kubernetes 1.12: Kubelet TLS Bootstrap and Azure Virtual Machine Scale Sets (VMSS) Move to General Availability",
				Date:  "2018-09-27",
			},
			want:    "https://kubernetes.io/blog/2018/09/27/kubernetes-1.12-kubelet-tls-bootstrap-and-azure-virtual-machine-scale-sets-vmss-move-to-general-availability/",
			wantErr: false,
		},
		{
			name: "Public blog post",
			path: "content/en/blog/_posts/2018-09-18-2018-linkerd-2-0.md",
			fm: &FrontMatter{
				Title: "Hands On With Linkerd 2.0",
				Date:  "2018-09-18",
			},
			want:    "https://kubernetes.io/blog/2018/09/18/hands-on-with-linkerd-2.0/",
			wantErr: false,
		},
		{
			name: "Public blog post",
			path: "content/en/blog/_posts/2018-07-16-kubernetes-1-11-release-interview.md",
			fm: &FrontMatter{
				Title: "How the sausage is made: the Kubernetes 1.11 release interview, from the Kubernetes Podcast",
				Date:  "2018-07-16",
			},
			want:    "https://kubernetes.io/blog/2018/07/16/how-the-sausage-is-made-the-kubernetes-1.11-release-interview-from-the-kubernetes-podcast/",
			wantErr: false,
		},
		{
			name: "Public blog post",
			path: "content/en/blog/_posts/2018-04-11-migrating-the-kubernetes-blog.md",
			fm: &FrontMatter{
				Title: "Local Persistent Volumes for Kubernetes Goes Beta",
				Slug:  "local-persistent-volumes-beta",
				Date:  "2018-04-13",
			},
			want:    "https://kubernetes.io/blog/2018/04/13/local-persistent-volumes-beta/",
			wantErr: false,
		},
		{
			name: "Public blog post",
			path: "content/en/blog/_posts/2017-02-00-Postgresql-Clusters-Kubernetes-Statefulsets.md",
			fm: &FrontMatter{
				Title: "Deploying PostgreSQL Clusters using StatefulSets",
				Slug:  "postgresql-clusters-kubernetes-statefulsets",
				Date:  "2017-02-24",
				Url:   "/blog/2017/02/Postgresql-Clusters-Kubernetes-Statefulsets",
			},
			want:    "https://kubernetes.io/blog/2017/02/postgresql-clusters-kubernetes-statefulsets/",
			wantErr: false,
		},
		{
			name: "Public blog post",
			path: "content/en/blog/_posts/2016-05-00-Coreosfest2016-Kubernetes-Community.md",
			fm: &FrontMatter{
				Title: "CoreOS Fest 2016: CoreOS and Kubernetes Community meet in Berlin (& San Francisco)",
				Slug:  "coreosfest2016-kubernetes-community",
				Date:  "2016-05-03",
				Url:   "/blog/2016/05/Coreosfest2016-Kubernetes-Community",
			},
			want:    "https://kubernetes.io/blog/2016/05/coreosfest2016-kubernetes-community/",
			wantErr: false,
		},
		{
			name: "Public blog post",
			path: "content/en/blog/_posts/2015-06-00-Cluster-Level-Logging-With-Kubernetes.md",
			fm: &FrontMatter{
				Title: "Cluster Level Logging with Kubernetes",
				Slug:  "cluster-level-logging-with-kubernetes",
				Date:  "2015-06-11",
				Url:   "/blog/2015/06/Cluster-Level-Logging-With-Kubernetes",
			},
			want:    "https://kubernetes.io/blog/2015/06/cluster-level-logging-with-kubernetes/",
			wantErr: false,
		},
		{
			name: "Public blog post",
			path: "content/en/blog/_posts/2015-03-00-Welcome-To-Kubernetes-Blog.md",
			fm: &FrontMatter{
				Title: "Welcome to the Kubernetes Blog!",
				Slug:  "welcome-to-kubernetes-blog",
				Date:  "2015-03-20",
				Url:   "/blog/2015/03/Welcome-To-Kubernetes-Blog",
			},
			want:    "https://kubernetes.io/blog/2015/03/welcome-to-kubernetes-blog/",
			wantErr: false,
		},
		{
			name: "Private blog post does not generate URL",
			path: "content/en/blog/_posts/2025-11-01-private-post.md",
			fm: &FrontMatter{
				Title: "Private Post",
				Slug:  "private-post",
				Date:  "2025-11-01",
				Url:   "/blog/2025/11/private-post",
				Build: &BuildSettings{
					Render: "never",
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Private blog post does not generate URL",
			path: "content/en/blog/_posts/2025-11-01-private-post.md",
			fm: &FrontMatter{
				Title: "Private Post",
				Slug:  "private-post",
				Date:  "2025-11-01",
				Url:   "/blog/2025/11/private-post",
				Build: &BuildSettings{
					Render: false,
				},
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.Convert(ctx, tt.path, tt.fm)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConverter_Community(t *testing.T) {
	config := &Config{
		BaseUrl: "https://kubernetes.io",
		ExistingUrls: map[string]bool{
			"https://kubernetes.io/community/":                             true,
			"https://kubernetes.io/community/code-of-conduct/":             true,
			"https://kubernetes.io/community/static/readme/":               true,
			"https://kubernetes.io/community/static/cncf-code-of-conduct/": true,
			"https://kubernetes.io/ja/community/":                          true,
			"https://kubernetes.io/ja/community/code-of-conduct/":          true,
			"https://kubernetes.io/zh-cn/community/code-of-conduct/":       true,
		},
		SupportedLangs: language.SupportedLanguages,
		SupportedExts:  []string{".md", ".html"},
		ValidSections:  []path.Category{path.Community},
	}
	converter := NewConverter(*config)
	ctx := context.Background()

	tests := []struct {
		name    string
		path    string
		fm      *FrontMatter
		want    string
		wantErr bool
	}{
		{
			name:    "Community index HTML",
			path:    "content/en/community/_index.html",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/community/",
			wantErr: false,
		},
		{
			name:    "Community code of conduct",
			path:    "content/en/community/code-of-conduct.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/community/code-of-conduct/",
			wantErr: false,
		},
		{
			name:    "Community static README",
			path:    "content/en/community/static/README.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/community/static/readme/",
			wantErr: false,
		},
		{
			name:    "Community static cncf code of conduct",
			path:    "content/en/community/static/cncf-code-of-conduct.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/community/static/cncf-code-of-conduct/",
			wantErr: false,
		},
		{
			name:    "Japanese community index",
			path:    "content/ja/community/_index.html",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/ja/community/",
			wantErr: false,
		},
		{
			name:    "Japanese community code of conduct",
			path:    "content/ja/community/code-of-conduct.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/ja/community/code-of-conduct/",
			wantErr: false,
		},
		{
			name:    "Chinese community code of conduct",
			path:    "content/zh-cn/community/code-of-conduct.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/zh-cn/community/code-of-conduct/",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.Convert(ctx, tt.path, tt.fm)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConverter_Example(t *testing.T) {
	config := &Config{
		BaseUrl: "https://kubernetes.io",
		ExistingUrls: map[string]bool{
			"https://kubernetes.io/examples/readme/":    true,
			"https://kubernetes.io/ja/examples/readme/": true,
		},
		SupportedLangs: language.SupportedLanguages,
		SupportedExts:  []string{".md", ".html"},
		ValidSections:  []path.Category{path.Example},
	}
	converter := NewConverter(*config)
	ctx := context.Background()

	tests := []struct {
		name    string
		path    string
		fm      *FrontMatter
		want    string
		wantErr bool
	}{
		{
			name:    "Examples README",
			path:    "content/en/examples/README.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/examples/readme/",
			wantErr: false,
		},
		{
			name:    "Japanese examples README",
			path:    "content/ja/examples/README.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/ja/examples/readme/",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.Convert(ctx, tt.path, tt.fm)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConverter_Includes(t *testing.T) {
	config := &Config{
		BaseUrl:        "https://kubernetes.io",
		ExistingUrls:   map[string]bool{},
		SupportedLangs: language.SupportedLanguages,
		SupportedExts:  []string{".md", ".html"},
		ValidSections:  []path.Category{path.Includes},
	}
	converter := NewConverter(*config)
	ctx := context.Background()

	tests := []struct {
		name    string
		path    string
		fm      *FrontMatter
		want    string
		wantErr bool
	}{
		{
			name:    "Includes index.md - should return empty URL",
			path:    "content/en/includes/index.md",
			fm:      &FrontMatter{},
			want:    "",
			wantErr: false,
		},
		{
			name:    "Includes default-storage-class-prereqs.md - should return empty URL",
			path:    "content/en/includes/default-storage-class-prereqs.md",
			fm:      &FrontMatter{},
			want:    "",
			wantErr: false,
		},
		{
			name:    "Includes federation-deprecation-warning-note.md - should return empty URL",
			path:    "content/en/includes/federation-deprecation-warning-note.md",
			fm:      &FrontMatter{},
			want:    "",
			wantErr: false,
		},
		{
			name:    "Includes task-tutorial-prereqs.md - should return empty URL",
			path:    "content/en/includes/task-tutorial-prereqs.md",
			fm:      &FrontMatter{},
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.Convert(ctx, tt.path, tt.fm)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConverter_Release(t *testing.T) {
	config := &Config{
		BaseUrl: "https://kubernetes.io",
		ExistingUrls: map[string]bool{
			"https://kubernetes.io/releases/":                        true,
			"https://kubernetes.io/releases/download/":               true,
			"https://kubernetes.io/releases/notes/":                  true,
			"https://kubernetes.io/releases/patch-releases/":         true,
			"https://kubernetes.io/releases/version-skew-policy/":    true,
			"https://kubernetes.io/ja/releases/":                     true,
			"https://kubernetes.io/ja/releases/version-skew-policy/": true,
			"https://kubernetes.io/zh-cn/releases/":                  true,
		},
		SupportedLangs: language.SupportedLanguages,
		SupportedExts:  []string{".md", ".html"},
		ValidSections:  []path.Category{path.Release},
	}
	converter := NewConverter(*config)
	ctx := context.Background()

	tests := []struct {
		name    string
		path    string
		fm      *FrontMatter
		want    string
		wantErr bool
	}{
		{
			name:    "Releases index",
			path:    "content/en/releases/_index.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/releases/",
			wantErr: false,
		},
		{
			name:    "Releases download",
			path:    "content/en/releases/download.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/releases/download/",
			wantErr: false,
		},
		{
			name:    "Releases notes",
			path:    "content/en/releases/notes.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/releases/notes/",
			wantErr: false,
		},
		{
			name:    "Releases patch-releases",
			path:    "content/en/releases/patch-releases.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/releases/patch-releases/",
			wantErr: false,
		},
		{
			name:    "Releases version-skew-policy",
			path:    "content/en/releases/version-skew-policy.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/releases/version-skew-policy/",
			wantErr: false,
		},
		{
			name:    "Japanese releases index",
			path:    "content/ja/releases/_index.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/ja/releases/",
			wantErr: false,
		},
		{
			name:    "Japanese releases version-skew-policy",
			path:    "content/ja/releases/version-skew-policy.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/ja/releases/version-skew-policy/",
			wantErr: false,
		},
		{
			name:    "Chinese releases index",
			path:    "content/zh-cn/releases/_index.md",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/zh-cn/releases/",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.Convert(ctx, tt.path, tt.fm)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConverter_Partner(t *testing.T) {
	config := &Config{
		BaseUrl: "https://kubernetes.io",
		ExistingUrls: map[string]bool{
			"https://kubernetes.io/partners/": true,
		},
		SupportedLangs: language.SupportedLanguages,
		SupportedExts:  []string{".md", ".html"},
		ValidSections:  []path.Category{path.Partner},
	}
	converter := NewConverter(*config)
	ctx := context.Background()

	tests := []struct {
		name    string
		path    string
		fm      *FrontMatter
		want    string
		wantErr bool
	}{
		{
			name:    "Partners index",
			path:    "content/en/partners/_index.html",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/partners/",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.Convert(ctx, tt.path, tt.fm)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConverter_Training(t *testing.T) {
	config := &Config{
		BaseUrl: "https://kubernetes.io",
		ExistingUrls: map[string]bool{
			"https://kubernetes.io/training/": true,
		},
		SupportedLangs: language.SupportedLanguages,
		SupportedExts:  []string{".md", ".html"},
		ValidSections:  []path.Category{path.Training},
	}
	converter := NewConverter(*config)
	ctx := context.Background()

	tests := []struct {
		name    string
		path    string
		fm      *FrontMatter
		want    string
		wantErr bool
	}{
		{
			name:    "Training index",
			path:    "content/en/training/_index.html",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/training/",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.Convert(ctx, tt.path, tt.fm)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConverter_Career(t *testing.T) {
	config := &Config{
		BaseUrl: "https://kubernetes.io",
		ExistingUrls: map[string]bool{
			"https://kubernetes.io/careers/": true,
		},
		SupportedLangs: language.SupportedLanguages,
		SupportedExts:  []string{".md", ".html"},
		ValidSections:  []path.Category{path.Career},
	}
	converter := NewConverter(*config)
	ctx := context.Background()

	tests := []struct {
		name    string
		path    string
		fm      *FrontMatter
		want    string
		wantErr bool
	}{
		{
			name:    "Careers index",
			path:    "content/en/careers/_index.html",
			fm:      &FrontMatter{},
			want:    "https://kubernetes.io/careers/",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.Convert(ctx, tt.path, tt.fm)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConverter_Glossary(t *testing.T) {
	config := &Config{
		BaseUrl:        "https://kubernetes.io",
		ExistingUrls:   map[string]bool{},
		SupportedLangs: language.SupportedLanguages,
		SupportedExts:  []string{".md", ".html"},
		ValidSections:  []path.Category{path.Docs},
	}
	converter := NewConverter(*config)
	ctx := context.Background()

	tests := []struct {
		name    string
		path    string
		fm      *FrontMatter
		want    string
		wantErr bool
	}{
		{
			name: "Glossary with tool tag",
			path: "content/en/docs/reference/glossary/addons.md",
			fm: &FrontMatter{
				Tags: []string{"tool"},
			},
			want:    "https://kubernetes.io/docs/reference/glossary/?tool=true",
			wantErr: false,
		},
		{
			name: "Glossary with user-type tag",
			path: "content/en/docs/reference/glossary/cluster-architect.md",
			fm: &FrontMatter{
				Tags: []string{"user-type"},
			},
			want:    "https://kubernetes.io/docs/reference/glossary/?user-type=true",
			wantErr: false,
		},
		{
			name: "Glossary with fundamental tag",
			path: "content/en/docs/reference/glossary/ephemeral-container.md",
			fm: &FrontMatter{
				Tags: []string{"fundamental"},
			},
			want:    "https://kubernetes.io/docs/reference/glossary/?fundamental=true",
			wantErr: false,
		},
		{
			name: "Glossary with no tags (fallback to all)",
			path: "content/en/docs/reference/glossary/sample.md",
			fm: &FrontMatter{
				Tags: []string{},
			},
			want:    "https://kubernetes.io/docs/reference/glossary/?all=true",
			wantErr: false,
		},
		{
			name: "Glossary with multiple tags",
			path: "content/en/docs/reference/glossary/multi-tag.md",
			fm: &FrontMatter{
				Tags: []string{"fundamental", "tool"},
			},
			want:    "https://kubernetes.io/docs/reference/glossary/?fundamental=true&tool=true",
			wantErr: false,
		},
		{
			name: "Glossary with id, no tags",
			path: "content/en/docs/reference/glossary/addons.md",
			fm: &FrontMatter{
				Id:   "addons",
				Tags: []string{},
			},
			want:    "https://kubernetes.io/docs/reference/glossary/?all=true#term-addons",
			wantErr: false,
		},
		{
			name: "Glossary with id and single tag",
			path: "content/en/docs/reference/glossary/cluster-architect.md",
			fm: &FrontMatter{
				Id:   "cluster-architect",
				Tags: []string{"user-type"},
			},
			want:    "https://kubernetes.io/docs/reference/glossary/?user-type=true#term-cluster-architect",
			wantErr: false,
		},
		{
			name: "Glossary with id and multiple tags",
			path: "content/en/docs/reference/glossary/ephemeral-container.md",
			fm: &FrontMatter{
				Id:   "ephemeral-container",
				Tags: []string{"fundamental", "core-object"},
			},
			want:    "https://kubernetes.io/docs/reference/glossary/?fundamental=true&core-object=true#term-ephemeral-container",
			wantErr: false,
		},
		{
			name: "Glossary with Japanese language",
			path: "content/ja/docs/reference/glossary/cluster-architect.md",
			fm: &FrontMatter{
				Tags: []string{"user-type"},
			},
			want:    "https://kubernetes.io/ja/docs/reference/glossary/?user-type=true",
			wantErr: false,
		},
		{
			name: "Glossary with Japanese language and id",
			path: "content/ja/docs/reference/glossary/addons.md",
			fm: &FrontMatter{
				Id:   "addons",
				Tags: []string{"tool"},
			},
			want:    "https://kubernetes.io/ja/docs/reference/glossary/?tool=true#term-addons",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.Convert(ctx, tt.path, tt.fm)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Convert() = %v, want %v", got, tt.want)
			}
		})
	}
}
