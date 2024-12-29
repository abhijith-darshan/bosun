package cluster

import (
	"bosun/pkg/internal"
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	corev1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

const (
	namespaceUpdateEvent = "namespaceUpdate"
	namespaceDeleteEvent = "namespaceDelete"
)

type BosunNamespace struct {
	client    client.Client
	clientSet *kubernetes.Clientset
	tracker   internal.InformerManager
}

func NewBosunNamespace(k8sClient client.Client, clientSet *kubernetes.Clientset) *BosunNamespace {
	informer := informers.NewSharedInformerFactoryWithOptions(clientSet, time.Minute*10).Core().V1().Namespaces().Informer()
	return &BosunNamespace{
		client:    k8sClient,
		clientSet: clientSet,
		tracker:   internal.NewInformerManager(informer),
	}
}

func (n *BosunNamespace) Create(ctx context.Context, namespace *corev1.Namespace) error {
	return n.client.Create(ctx, namespace)
}

func (n *BosunNamespace) Get(ctx context.Context, name string) (*corev1.Namespace, error) {
	ns := &corev1.Namespace{}
	err := n.client.Get(ctx, client.ObjectKey{Name: name}, ns)
	if err != nil {
		return nil, err
	}
	return ns, nil
}

func (n *BosunNamespace) List(ctx context.Context) ([]corev1.Namespace, error) {
	list := &corev1.NamespaceList{}
	err := n.client.List(ctx, list)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (n *BosunNamespace) ListFromStore() []interface{} {
	return n.tracker.Store().List()
}

func (n *BosunNamespace) Update(ctx context.Context, namespace *corev1.Namespace) error {
	oldObj, err := n.Get(ctx, namespace.GetName())
	if err != nil {
		return err
	}
	patch := client.MergeFrom(oldObj)
	return n.client.Patch(ctx, namespace, patch)
}

func (n *BosunNamespace) Delete(ctx context.Context, name string) error {
	ns, err := n.Get(ctx, name)
	if apiErrors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	err = n.client.Delete(ctx, ns)
	if apiErrors.IsGone(err) {
		return nil
	}
	return err
}

func (n *BosunNamespace) StartTracker(ctx context.Context) error {
	handlers := internal.EventHandlers{
		AddFunc:    internal.AddInformerEvent(ctx, namespaceUpdateEvent),
		UpdateFunc: internal.UpdateInformerEvent(ctx, namespaceUpdateEvent),
		DeleteFunc: internal.DeleteInformerEvent(ctx, namespaceDeleteEvent),
	}
	err := n.tracker.Start(ctx, handlers)
	if err != nil {
		return err
	}
	runtime.LogPrint(ctx, "namespace informer started")
	return nil
}

func (n *BosunNamespace) StopTracker(ctx context.Context) {
	if n.tracker == nil {
		return
	}
	n.tracker.Stop()
	n.tracker = nil
	runtime.LogPrint(ctx, "namespace informer stopped")
}

func (n *BosunNamespace) IsTrackerInit() bool {
	if n.tracker == nil {
		return false
	}
	return true
}
