package resource

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type ConfigMapResource struct {
	Clientset kubernetes.Interface
}

func (r *ConfigMapResource) Get(ctx context.Context, name, namespace string, opts metav1.GetOptions) (metav1.Object, error) {
	return r.Clientset.CoreV1().ConfigMaps(namespace).Get(ctx, name, opts)
}

func (r *ConfigMapResource) Create(ctx context.Context, resource metav1.Object, opts metav1.CreateOptions) (metav1.Object, error) {
	cm := resource.(*corev1.ConfigMap)
	return r.Clientset.CoreV1().ConfigMaps(cm.Namespace).Create(ctx, cm, opts)
}

func (r *ConfigMapResource) Update(ctx context.Context, resource metav1.Object, opts metav1.UpdateOptions) (metav1.Object, error) {
	cm := resource.(*corev1.ConfigMap)
	return r.Clientset.CoreV1().ConfigMaps(cm.Namespace).Update(ctx, cm, opts)
}

func (r *ConfigMapResource) Delete(ctx context.Context, name, namespace string, opts metav1.DeleteOptions) error {
	return r.Clientset.CoreV1().ConfigMaps(namespace).Delete(ctx, name, opts)
}

func (r *ConfigMapResource) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return r.Clientset.CoreV1().ConfigMaps("").Watch(ctx, opts)
}
