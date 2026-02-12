package sidecar

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Return the client for kubernetes
func kubernetesClient() (*kubernetes.Clientset, error) {
	// Config in-cluster (si ton code tourne dans un Pod)
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("error on get cluster config: %w", err)
	}

	return kubernetes.NewForConfig(config)
}

// Watch a Secret and apply the treatment when it's modified
// From the context
//   - namespace String
//   - secretName String
//
// Add in the context :
//   - secret *v1.Secret
func WatchSecretModified(c *Callback) error {
	// Check context
	namespace, ok := c.Context["namespace"].(string)
	if !ok {
		return fmt.Errorf("missing or incorrect type for the context: namespace")
	}
	secretName, ok := c.Context["secretName"].(string)
	if !ok {
		return fmt.Errorf("missing or incorrect type for the context: secretName")
	}

	clientset, err := kubernetesClient()
	if err != nil {
		return fmt.Errorf("error on kubernetes client: %w", err)
	}

	// Watch on the Secret inside the namespace
	watchInterface, err := clientset.CoreV1().Secrets(namespace).Watch(
		context.TODO(),
		metav1.ListOptions{
			FieldSelector: fmt.Sprintf("metadata.name=%s", secretName),
		},
	)
	if err != nil {
		return fmt.Errorf("error on watching %s.%s : %w", namespace, secretName, err)
	}

	go func() {
		fmt.Println("[Sidecar : WatchSecretModified] Routine launched")
		// Boucle sur les événements
		for event := range watchInterface.ResultChan() {
			switch event.Type {
			case watch.Modified:
				fmt.Println("Secret modified")
				secret, ok := event.Object.(*v1.Secret)
				if !ok {
					fmt.Println("Object unexpected")
					continue
				}
				c.Context["secret"] = secret
				err = c.Clone().CallNext()
				if err != nil {
					fmt.Println("error on watching secret callback : ", err)
				}
			}
		}
		fmt.Println("[Sidecar : WatchSecretModified] Routine stopped")
	}()
	return nil
}

// Read a Secret in Kubernetes
// From the context
//   - namespace String
//   - secretName String
//
// Add in the context :
//   - secret *v1.Secret
func ReadSecret(c *Callback) error {
	// Check context
	namespace, ok := c.Context["namespace"].(string)
	if !ok {
		return fmt.Errorf("missing or incorrect type for the context: namespace")
	}
	secretName, ok := c.Context["secretName"].(string)
	if !ok {
		return fmt.Errorf("missing or incorrect type for the context: secretName")
	}

	clientset, err := kubernetesClient()
	if err != nil {
		return fmt.Errorf("error on kubernetes client: %w", err)
	}

	// Get a Secret by his name inside the namespace
	secret, err := clientset.CoreV1().Secrets(namespace).Get(
		context.TODO(),
		secretName,
		metav1.GetOptions{},
	)
	c.Context["secret"] = secret

	if err != nil {
		return fmt.Errorf("get Secret call fail: %w", err)
	}

	return c.CallNext()
}

// Patch a secret
// From the context
//   - namespace String
//   - secretName String
//   - secretKey String
//   - secretData []byte
func PatchSecret(c *Callback) error {
	// Check context
	namespace, ok := c.Context["namespace"].(string)
	if !ok {
		return fmt.Errorf("missing or incorrect type for the context: namespace")
	}
	secretName, ok := c.Context["secretName"].(string)
	if !ok {
		return fmt.Errorf("missing or incorrect type for the context: secretName")
	}
	secretKey, ok := c.Context["secretKey"].(string)
	if !ok {
		return fmt.Errorf("missing or incorrect type for the context: secretKey")
	}
	secretData, ok := c.Context["secretData"].([]byte)
	if !ok {
		return fmt.Errorf("missing or incorrect type for the context: secretData")
	}

	clientset, err := kubernetesClient()
	if err != nil {
		return fmt.Errorf("error on kubernetes client: %w", err)
	}

	// secret's data must be encoded in base64
	encoded := base64.StdEncoding.EncodeToString(secretData)
	fmt.Println("[Sidecar : PatchSecret] patching", namespace, "/", secretName, "key:", secretKey, "size:", len(secretData))

	// Patch JSON to replace/add the key
	patch := map[string]interface{}{
		"data": map[string]string{
			secretKey: encoded,
		},
	}

	patchBytes, err := json.Marshal(patch)
	if err != nil {
		return fmt.Errorf("error marshalling patch: %w", err)
	}

	// Send the patch
	_, err = clientset.CoreV1().Secrets(namespace).Patch(
		context.TODO(),
		secretName,
		types.MergePatchType,
		patchBytes,
		metav1.PatchOptions{},
	)
	if err != nil {
		return fmt.Errorf("error patching secret: %w", err)
	}
	fmt.Println("[Sidecar : PatchSecret] patched", secretName, secretKey)

	return c.CallNext()
}
