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

func TestConfigMapResource_Get(t *testing.T) {
	clientset := fake.NewSimpleClientset(&corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-configmap",
			Namespace: "default",
		},
	})

	cmr := &resource.ConfigMapResource{Clientset: clientset}

	cm, err := cmr.Get(context.Background(), "test-configmap", "default", metav1.GetOptions{})
	assert.NoError(t, err)
	assert.NotNil(t, cm)
	assert.Equal(t, "test-configmap", cm.GetName())
}

func TestConfigMapResource_Create(t *testing.T) {
	clientset := fake.NewSimpleClientset()

	cmr := &resource.ConfigMapResource{Clientset: clientset}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-configmap",
			Namespace: "default",
		},
	}

	createdCm, err := cmr.Create(context.Background(), cm, metav1.CreateOptions{})
	assert.NoError(t, err)
	assert.NotNil(t, createdCm)
	assert.Equal(t, "test-configmap", createdCm.GetName())
}

func TestConfigMapResource_Update(t *testing.T) {
	clientset := fake.NewSimpleClientset(&corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-configmap",
			Namespace: "default",
		},
	})

	cmr := &resource.ConfigMapResource{Clientset: clientset}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-configmap",
			Namespace: "default",
		},
		Data: map[string]string{
			"key": "value",
		},
	}

	updatedCm, err := cmr.Update(context.Background(), cm, metav1.UpdateOptions{})
	assert.NoError(t, err)
	assert.NotNil(t, updatedCm)
	assert.Equal(t, "value", updatedCm.(*corev1.ConfigMap).Data["key"])
}

func TestConfigMapResource_Delete(t *testing.T) {
	clientset := fake.NewSimpleClientset(&corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-configmap",
			Namespace: "default",
		},
	})

	cmr := &resource.ConfigMapResource{Clientset: clientset}

	err := cmr.Delete(context.Background(), "test-configmap", "default", metav1.DeleteOptions{})
	assert.NoError(t, err)

	_, err = cmr.Get(context.Background(), "test-configmap", "default", metav1.GetOptions{})
	assert.Error(t, err)
}

func TestConfigMapResource_Watch(t *testing.T) {
	clientset := fake.NewSimpleClientset()

	cmr := &resource.ConfigMapResource{Clientset: clientset}

	watcher, err := cmr.Watch(context.Background(), metav1.ListOptions{})
	assert.NoError(t, err)
	assert.NotNil(t, watcher)

	go func() {
		clientset.CoreV1().ConfigMaps("default").Create(context.Background(), &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-configmap",
				Namespace: "default",
			},
		}, metav1.CreateOptions{})
	}()

	event := <-watcher.ResultChan()
	assert.Equal(t, watch.Added, event.Type)
	cm := event.Object.(*corev1.ConfigMap)
	assert.Equal(t, "test-configmap", cm.GetName())
}
