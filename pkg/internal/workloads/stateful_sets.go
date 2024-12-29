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
	statefulSetUpdateEvent = "statefulSetUpdate"
	statefulSetDeleteEvent = "statefulSetDelete"
)

type BosunStatefulSet struct {
	client    client.Client
	clientSet *kubernetes.Clientset
	tracker   internal.InformerManager
}

func NewBosunStatefulSet(client client.Client, clientSet *kubernetes.Clientset) *BosunStatefulSet {
	informer := informers.NewSharedInformerFactoryWithOptions(clientSet, time.Minute*10).Apps().V1().StatefulSets().Informer()
	return &BosunStatefulSet{
		client:    client,
		clientSet: clientSet,
		tracker:   internal.NewInformerManager(informer),
	}
}

func (s *BosunStatefulSet) Get(ctx context.Context, name, namespace string) (*appsv1.StatefulSet, error) {
	ss := &appsv1.StatefulSet{}
	err := s.client.Get(ctx, client.ObjectKey{
		Name:      name,
		Namespace: namespace,
	}, ss)
	if err != nil {
		return nil, err
	}
	return ss, nil
}

func (s *BosunStatefulSet) List(ctx context.Context) ([]appsv1.StatefulSet, error) {
	list := &appsv1.StatefulSetList{}
	err := s.client.List(ctx, list)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (s *BosunStatefulSet) ListFromStore() []interface{} {
	return s.tracker.Store().List()
}

func (s *BosunStatefulSet) Update(ctx context.Context, statefulSet *appsv1.StatefulSet) error {
	oldObj, err := s.Get(ctx, statefulSet.GetName(), statefulSet.GetNamespace())
	if err != nil {
		return err
	}
	patch := client.MergeFrom(oldObj)
	return s.client.Patch(ctx, statefulSet, patch)
}

func (s *BosunStatefulSet) Delete(ctx context.Context, name, namespace string) error {
	statefulSet, err := s.Get(ctx, name, namespace)
	if apiErrors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	err = s.client.Delete(ctx, statefulSet)
	if apiErrors.IsGone(err) {
		return nil
	}
	return err
}

func (s *BosunStatefulSet) StartTracker(ctx context.Context) error {
	handlers := internal.EventHandlers{
		AddFunc:    internal.AddInformerEvent(ctx, statefulSetUpdateEvent),
		UpdateFunc: internal.UpdateInformerEvent(ctx, statefulSetUpdateEvent),
		DeleteFunc: internal.DeleteInformerEvent(ctx, statefulSetDeleteEvent),
	}
	err := s.tracker.Start(ctx, handlers)
	if err != nil {
		return err
	}
	runtime.LogPrint(ctx, "statefulSet informer started")
	return nil
}

func (s *BosunStatefulSet) StopTracker(ctx context.Context) {
	if s.tracker == nil {
		return
	}
	s.tracker.Stop()
	s.tracker = nil
	runtime.LogPrintf(ctx, "statefulSet informer stopped")
}

func (s *BosunStatefulSet) IsTrackerInit() bool {
	if s.tracker == nil {
		return false
	}
	return true
}
