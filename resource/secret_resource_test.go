package resource_test

import (
	"context"
	"testing"

	"github.com/lab42/mirror/resource"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/fake"
)

func TestSecretResource_Get(t *testing.T) {
	clientset := fake.NewSimpleClientset(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-secret",
			Namespace: "default",
		},
	})

	sr := &resource.SecretResource{Clientset: clientset}

	secret, err := sr.Get(context.Background(), "test-secret", "default", metav1.GetOptions{})
	assert.NoError(t, err)
	assert.NotNil(t, secret)
	assert.Equal(t, "test-secret", secret.GetName())
}

func TestSecretResource_Create(t *testing.T) {
	clientset := fake.NewSimpleClientset()

	sr := &resource.SecretResource{Clientset: clientset}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-secret",
			Namespace: "default",
		},
	}

	createdSecret, err := sr.Create(context.Background(), secret, metav1.CreateOptions{})
	assert.NoError(t, err)
	assert.NotNil(t, createdSecret)
	assert.Equal(t, "test-secret", createdSecret.GetName())
}

func TestSecretResource_Update(t *testing.T) {
	clientset := fake.NewSimpleClientset(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-secret",
			Namespace: "default",
		},
	})

	sr := &resource.SecretResource{Clientset: clientset}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-secret",
			Namespace: "default",
		},
		Data: map[string][]byte{
			"key": []byte("value"),
		},
	}

	updatedSecret, err := sr.Update(context.Background(), secret, metav1.UpdateOptions{})
	assert.NoError(t, err)
	assert.NotNil(t, updatedSecret)
	assert.Equal(t, "value", string(updatedSecret.(*corev1.Secret).Data["key"]))
}

func TestSecretResource_Delete(t *testing.T) {
	clientset := fake.NewSimpleClientset(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-secret",
			Namespace: "default",
		},
	})

	sr := &resource.SecretResource{Clientset: clientset}

	err := sr.Delete(context.Background(), "test-secret", "default", metav1.DeleteOptions{})
	assert.NoError(t, err)

	_, err = sr.Get(context.Background(), "test-secret", "default", metav1.GetOptions{})
	assert.Error(t, err)
}

func TestSecretResource_Watch(t *testing.T) {
	clientset := fake.NewSimpleClientset()

	sr := &resource.SecretResource{Clientset: clientset}

	watcher, err := sr.Watch(context.Background(), metav1.ListOptions{})
	assert.NoError(t, err)
	assert.NotNil(t, watcher)

	go func() {
		clientset.CoreV1().Secrets("default").Create(context.Background(), &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-secret",
				Namespace: "default",
			},
		}, metav1.CreateOptions{})
	}()

	event := <-watcher.ResultChan()
	assert.Equal(t, watch.Added, event.Type)
	secret := event.Object.(*corev1.Secret)
	assert.Equal(t, "test-secret", secret.GetName())
}
