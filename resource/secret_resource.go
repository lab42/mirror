package resource

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type SecretResource struct {
	Clientset kubernetes.Interface
}

func (r *SecretResource) Get(ctx context.Context, name, namespace string, opts metav1.GetOptions) (metav1.Object, error) {
	return r.Clientset.CoreV1().Secrets(namespace).Get(ctx, name, opts)
}

func (r *SecretResource) Create(ctx context.Context, resource metav1.Object, opts metav1.CreateOptions) (metav1.Object, error) {
	secret := resource.(*corev1.Secret)
	return r.Clientset.CoreV1().Secrets(secret.Namespace).Create(ctx, secret, opts)
}

func (r *SecretResource) Update(ctx context.Context, resource metav1.Object, opts metav1.UpdateOptions) (metav1.Object, error) {
	secret := resource.(*corev1.Secret)
	return r.Clientset.CoreV1().Secrets(secret.Namespace).Update(ctx, secret, opts)
}

func (r *SecretResource) Delete(ctx context.Context, name, namespace string, opts metav1.DeleteOptions) error {
	return r.Clientset.CoreV1().Secrets(namespace).Delete(ctx, name, opts)
}

func (r *SecretResource) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return r.Clientset.CoreV1().Secrets("").Watch(ctx, opts)
}
