// reflector package
package reflector

import (
	"context"
	"errors"
	"strings"

	corev1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

// ResourceType represents the type of resource.
type ResourceType int

const (
	ConfigMap ResourceType = iota
	Secret
)

// Resource defines the common operations for both ConfigMaps and Secrets.
type Resource interface {
	Get(ctx context.Context, name, namespace string, opts metav1.GetOptions) (metav1.Object, error)
	Create(ctx context.Context, resource metav1.Object, opts metav1.CreateOptions) (metav1.Object, error)
	Update(ctx context.Context, resource metav1.Object, opts metav1.UpdateOptions) (metav1.Object, error)
	Delete(ctx context.Context, name, namespace string, opts metav1.DeleteOptions) error
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
}

// ReflectorConfig holds configuration options for creating a Reflector.
type ReflectorConfig struct {
	ClientSet    kubernetes.Interface
	Resource     Resource
	ResourceType ResourceType
}

// Reflector is responsible for reflecting resources across namespaces.
type Reflector struct {
	ClientSet kubernetes.Interface
	Resource  Resource
	Type      ResourceType
}

// NewReflector creates a new instance of Reflector using the provided configuration.
func NewReflector(config ReflectorConfig) *Reflector {
	return &Reflector{
		ClientSet: config.ClientSet,
		Resource:  config.Resource,
		Type:      config.ResourceType,
	}
}

// Run starts reflecting resources across namespaces.
func (r *Reflector) Run(ctx context.Context) error {
	watcher, err := r.Resource.Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}
	defer watcher.Stop()

	for event := range watcher.ResultChan() {
		switch event.Type {
		case watch.Added, watch.Modified:
			resource, err := meta.Accessor(event.Object)
			if err != nil {
				return err
			}
			if err := r.reflectResource(ctx, resource); err != nil {
				return err
			}
		case watch.Deleted:
			resource, err := meta.Accessor(event.Object)
			if err != nil {
				return err
			}
			if err := r.deleteReflectedResources(ctx, resource); err != nil {
				return err
			}
		}
	}
	return nil
}

// reflectResource reflects the given resource across namespaces.
func (r *Reflector) reflectResource(ctx context.Context, resource metav1.Object) error {
	annotations := resource.GetAnnotations()
	val, ok := annotations["reflect.lab42.io/to-namespaces"]
	if !ok || len(val) == 0 {
		return nil // No reflection needed
	}

	namespaces := strings.Split(val, ",")
	for _, namespace := range namespaces {
		if namespace == resource.GetNamespace() {
			return errors.New("cannot overwrite self")
		}

		reflectedResource, err := r.prepareReflectedResource(resource, namespace)
		if err != nil {
			return err
		}

		_, err = r.Resource.Get(ctx, reflectedResource.GetName(), namespace, metav1.GetOptions{})
		if err != nil {
			if k8sErrors.IsNotFound(err) {
				_, err := r.Resource.Create(ctx, reflectedResource, metav1.CreateOptions{})
				if err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			_, err := r.Resource.Update(ctx, reflectedResource, metav1.UpdateOptions{})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// prepareReflectedResource prepares a deep copy of the resource for reflection.
func (r *Reflector) prepareReflectedResource(resource metav1.Object, namespace string) (metav1.Object, error) {
	var reflectedResource metav1.Object

	switch r.Type {
	case ConfigMap:
		cm, ok := resource.(*corev1.ConfigMap)
		if !ok {
			return nil, errors.New("unexpected resource type")
		}
		reflectedResource = cm.DeepCopy()
	case Secret:
		secret, ok := resource.(*corev1.Secret)
		if !ok {
			return nil, errors.New("unexpected resource type")
		}
		reflectedResource = secret.DeepCopy()
	default:
		return nil, errors.New("unknown resource type")
	}

	reflectedResource.SetResourceVersion("") // Clear the resource version
	reflectedResource.SetUID("")             // Clear the UID to prevent conflicts
	reflectedResource.SetNamespace(namespace)
	deleteAnnotations(reflectedResource)
	return reflectedResource, nil
}

// deleteReflectedResources deletes the reflected resources from the specified namespaces.
func (r *Reflector) deleteReflectedResources(ctx context.Context, resource metav1.Object) error {
	annotations := resource.GetAnnotations()
	val, ok := annotations["reflect.lab42.io/to-namespaces"]
	if !ok || len(val) == 0 {
		return nil // No reflection needed
	}

	namespaces := strings.Split(val, ",")
	for _, namespace := range namespaces {
		if err := r.Resource.Delete(ctx, resource.GetName(), namespace, metav1.DeleteOptions{}); err != nil {
			return err
		}
	}
	return nil
}

// deleteAnnotations deletes the reflection-related annotations from the resource.
func deleteAnnotations(resource metav1.Object) {
	annotations := resource.GetAnnotations()
	delete(annotations, "reflect.lab42.io/to-namespaces")
	resource.SetAnnotations(annotations)
}
