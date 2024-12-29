package workloads

import (
	"bosun/pkg/internal"
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	appsv1 "k8s.io/api/apps/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

const (
	replicaSetUpdateEvent = "replicaSetUpdate"
	replicaSetDeleteEvent = "replicaSetDelete"
)

type BosunReplicaSet struct {
	client    client.Client
	clientSet *kubernetes.Clientset
	tracker   internal.InformerManager
}

func NewBosunReplicaSet(client client.Client, clientSet *kubernetes.Clientset) *BosunReplicaSet {
	informer := informers.NewSharedInformerFactoryWithOptions(clientSet, time.Minute*10).Apps().V1().ReplicaSets().Informer()
	return &BosunReplicaSet{
		client:    client,
		clientSet: clientSet,
		tracker:   internal.NewInformerManager(informer),
	}
}

func (r *BosunReplicaSet) Get(ctx context.Context, name, namespace string) (*appsv1.ReplicaSet, error) {
	rs := &appsv1.ReplicaSet{}
	err := r.client.Get(ctx, client.ObjectKey{
		Name:      name,
		Namespace: namespace,
	}, rs)
	if err != nil {
		return nil, err
	}
	return rs, nil
}

func (r *BosunReplicaSet) List(ctx context.Context) ([]appsv1.ReplicaSet, error) {
	list := &appsv1.ReplicaSetList{}
	err := r.client.List(ctx, list)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (r *BosunReplicaSet) ListFromStore() []interface{} {
	return r.tracker.Store().List()
}

func (r *BosunReplicaSet) Update(ctx context.Context, replicaSet *appsv1.ReplicaSet) error {
	oldObj, err := r.Get(ctx, replicaSet.GetName(), replicaSet.GetNamespace())
	if err != nil {
		return err
	}
	patch := client.MergeFrom(oldObj)
	return r.client.Patch(ctx, replicaSet, patch)
}

func (r *BosunReplicaSet) Delete(ctx context.Context, name, namespace string) error {
	replicaSet, err := r.Get(ctx, name, namespace)
	if apiErrors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	err = r.client.Delete(ctx, replicaSet)
	if apiErrors.IsGone(err) {
		return nil
	}
	return err
}

func (r *BosunReplicaSet) StartTracker(ctx context.Context) error {
	handlers := internal.EventHandlers{
		AddFunc:    internal.AddInformerEvent(ctx, replicaSetUpdateEvent),
		UpdateFunc: internal.UpdateInformerEvent(ctx, replicaSetUpdateEvent),
		DeleteFunc: internal.DeleteInformerEvent(ctx, replicaSetDeleteEvent),
	}
	err := r.tracker.Start(ctx, handlers)
	if err != nil {
		return err
	}
	runtime.LogPrint(ctx, "replicaSet informer started")
	return nil
}

func (r *BosunReplicaSet) StopTracker(ctx context.Context) {
	if r.tracker == nil {
		return
	}
	r.tracker.Stop()
	r.tracker = nil
	runtime.LogPrintf(ctx, "replicaSet informer stopped")
}

func (r *BosunReplicaSet) IsTrackerInit() bool {
	if r.tracker == nil {
		return false
	}
	return true
}
