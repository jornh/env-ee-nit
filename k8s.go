package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// GetK8sClient initializes a standard kubernetes clientset
func GetK8sClient() (*kubernetes.Clientset, error) {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

// GetK8sFilesClient reads all .yaml files in a directory and
// populates a fake clientset to simulate a real cluster.
func GetK8sFilesClient(dir string) (kubernetes.Interface, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil || len(files) == 0 {
		return nil, fmt.Errorf("no yaml files found in %s", dir)
	}

	var objects []metav1.Object
	decode := scheme.Codecs.UniversalDeserializer().Decode

	for _, f := range files {
		content, err := os.ReadFile(f)
		if err != nil {
			continue
		}

		obj, _, err := decode(content, nil, nil)
		if err != nil {
			continue
		}

		if dep, ok := obj.(*appsv1.Deployment); ok {
			objects = append(objects, dep)
		}
	}

	// Create a fake client pre-populated with these "files"
	return fake.NewSimpleClientset(convertToRuntimeObjects(objects)...), nil
}

// Helper to convert slice for fake.NewSimpleClientset
func convertToRuntimeObjects(objs []metav1.Object) []runtime.Object {
	var out []runtime.Object
	for _, o := range objs {
		out = append(out, o.(runtime.Object))
	}
	return out
}

// FetchVersionsFromNamespace scans deployments to find app names and versions
func FetchVersionsFromNamespace(clientset kubernetes.Interface, namespace string) ([]Version, error) {
	var versions []Version

	// We query Deployments; you could also query StatefulSets or Pods depending on your stack
	deps, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, d := range deps.Items {
		// Taking the first container image as the source of truth for the app version
		if len(d.Spec.Template.Spec.Containers) > 0 {
			fullImage := d.Spec.Template.Spec.Containers[0].Image

			// Parse "repo/app-name:version"
			app, version := parseImage(fullImage)

			versions = append(versions, Version{
				App:     app,
				Env:     namespace,
				Version: version,
			})
		}
	}

	return versions, nil
}

// parseImage splits "myorg/app-a:1.2.3" into ("app-a", "1.2.3")
func parseImage(image string) (string, string) {
	// 1. Isolate the part after the last slash (the name and tag/digest)
	// Example: "ghcr.io/org/my-app:1.2.3" -> "my-app:1.2.3"
	parts := strings.Split(image, "/")
	lastPart := parts[len(parts)-1]

	// 2. Check for SHA digest (@sha256:...)
	// Example: "my-app@sha256:abcdef"
	if strings.Contains(lastPart, "@") {
		split := strings.SplitN(lastPart, "@", 2)
		return split[0], split[1]
	}

	// 3. Check for standard Tag (:)
	// Example: "my-app:v1.0.0"
	if strings.Contains(lastPart, ":") {
		split := strings.SplitN(lastPart, ":", 2)
		return split[0], split[1]
	}

	// 4. Default if no tag or digest is found
	return lastPart, "latest"
}
