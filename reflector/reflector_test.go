package reflector_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	"github.com/lab42/mirror/reflector"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

// MockResource is a mock implementation of the Resource interface.
type MockResource struct {
	mock.Mock
}

// Get mocks the Get method of Resource.
func (m *MockResource) Get(ctx context.Context, name, namespace string, opts metav1.GetOptions) (metav1.Object, error) {
	args := m.Called(ctx, name, namespace, opts)
	return args.Get(0).(metav1.Object), args.Error(1)
}

// Create mocks the Create method of Resource.
func (m *MockResource) Create(ctx context.Context, resource metav1.Object, opts metav1.CreateOptions) (metav1.Object, error) {
	args := m.Called(ctx, resource, opts)
	return args.Get(0).(metav1.Object), args.Error(1)
}

// Update mocks the Update method of Resource.
func (m *MockResource) Update(ctx context.Context, resource metav1.Object, opts metav1.UpdateOptions) (metav1.Object, error) {
	args := m.Called(ctx, resource, opts)
	return args.Get(0).(metav1.Object), args.Error(1)
}

// Delete mocks the Delete method of Resource.
func (m *MockResource) Delete(ctx context.Context, name, namespace string, opts metav1.DeleteOptions) error {
	args := m.Called(ctx, name, namespace, opts)
	return args.Error(0)
}

// Watch mocks the Watch method of Resource.
func (m *MockResource) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(watch.Interface), args.Error(1)
}

// TestReflector_Run tests the Run method of Reflector.
func TestReflector_Run(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	mockResource := new(MockResource)
	reflector := reflector.NewReflector(reflector.ReflectorConfig{ClientSet: kubernetes.NewForConfigOrDie(clientset), Resource: mockResource, ResourceType: reflector.ConfigMap})

	ctx := context.Background()

	// Mock watch result
	watchResult := make(chan watch.Event)
	mockWatcher := &MockWatcher{resultChan: watchResult}
	mockResource.On("Watch", ctx, mock.Anything).Return(mockWatcher, nil)

	// Mock watch events
	watchEvents := []watch.Event{
		{Type: watch.Added, Object: &corev1.ConfigMap{}},
		{Type: watch.Modified, Object: &corev1.ConfigMap{}},
		{Type: watch.Deleted, Object: &corev1.ConfigMap{}},
	}

	go func() {
		for _, event := range watchEvents {
			watchResult <- event
		}
		close(watchResult)
	}()

	err := reflector.Run(ctx)
	assert.NoError(t, err)

	// Verify that the expected methods were called
	mockResource.AssertCalled(t, "Watch", ctx, mock.Anything)
}

// MockWatcher is a mock implementation of watch.Interface.
type MockWatcher struct {
	mock.Mock
	resultChan <-chan watch.Event
}

// ResultChan returns the result channel.
func (m *MockWatcher) ResultChan() <-chan watch.Event {
	return m.resultChan
}

// Stop stops the watcher.
func (m *MockWatcher) Stop() {
	// Do nothing for now
}
