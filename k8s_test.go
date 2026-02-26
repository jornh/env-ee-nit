package main

import (
	"testing"
	"k8s.io/client-go/kubernetes/fake"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
)

func TestParseImage(t *testing.T) {
	tests := []struct {
		input    string
		wantApp  string
		wantVer  string
	}{
		{"nginx", "nginx", "latest"},
		{"my-org/app:1.2.3", "app", "1.2.3"},
		{"ghcr.io/org/sub/app:v1", "app", "v1"},
		{"app@sha256:digest", "app", "sha256:digest"},
		{"registry:5000/app:beta", "app", "beta"},
	}

	for _, tt := range tests {
		app, ver := parseImage(tt.input)
		if app != tt.wantApp || ver != tt.wantVer {
			t.Errorf("parseImage(%q) = (%q, %q); want (%q, %q)",
				tt.input, app, ver, tt.wantApp, tt.wantVer)
		}
	}
}

func TestFetchVersions(t *testing.T) {
	// Create a fake clientset with a pre-defined deployment
	client := fake.NewSimpleClientset(&appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: "app-a", Namespace: "dev"},
		Spec: appsv1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{Image: "org/app-a:0.1.2"}},
				},
			},
		},
	})

	versions, err := FetchVersionsFromNamespace(client, "dev")
	if err != nil || len(versions) != 1 {
		t.Fatalf("Expected 1 version, got %d", len(versions))
	}

	if versions[0].Version != "0.1.2" {
		t.Errorf("Expected 0.1.2, got %s", versions[0].Version)
	}
}

func TestFetchVersionsFromGoldenFiles(t *testing.T) {
	// 1. Point to your testdata directory
	// Ensure testdata/deployment.yaml exists with a valid Deployment manifest
	client, err := GetK8sFilesClient("testdata")
	if err != nil {
		t.Fatalf("Failed to load golden files: %v", err)
	}

	// 2. The namespace here must match the 'namespace' field in your golden YAML
	// or be empty if the YAML doesn't specify one.
	//versions, err := FetchVersionsFromNamespace(client, "default")
	versions, err := FetchVersionsFromNamespace(client, "dev")
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	if len(versions) == 0 {
		t.Fatal("Expected to find versions in golden files, but got none")
	}

	t.Logf("Found app: %s version: %s", versions[0].App, versions[0].Version)
}
