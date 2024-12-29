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
	daemonSetUpdateEvent = "daemonSetUpdate"
	daemonSetDeleteEvent = "daemonSetDelete"
)

type BosunDaemonSet struct {
	client    client.Client
	clientSet *kubernetes.Clientset
	tracker   internal.InformerManager
}

func NewBosunDaemonSet(k8sClient client.Client, clientSet *kubernetes.Clientset) *BosunDaemonSet {
	informer := informers.NewSharedInformerFactoryWithOptions(clientSet, time.Minute*10).Apps().V1().DaemonSets().Informer()
	return &BosunDaemonSet{
		client:    k8sClient,
		clientSet: clientSet,
		tracker:   internal.NewInformerManager(informer),
	}
}

func (d *BosunDaemonSet) Get(ctx context.Context, name, namespace string) (*appsv1.DaemonSet, error) {
	daemon := &appsv1.DaemonSet{}
	err := d.client.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, daemon)
	if err != nil {
		return nil, err
	}
	return daemon, nil
}

func (d *BosunDaemonSet) List(ctx context.Context) ([]appsv1.DaemonSet, error) {
	list := &appsv1.DaemonSetList{}
	err := d.client.List(ctx, list)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (d *BosunDaemonSet) ListFromStore() []interface{} {
	return d.tracker.Store().List()
}

func (d *BosunDaemonSet) Update(ctx context.Context, daemonSet *appsv1.DaemonSet) error {
	oldObj, err := d.Get(ctx, daemonSet.GetName(), daemonSet.GetNamespace())
	if err != nil {
		return err
	}
	patch := client.MergeFrom(oldObj)
	return d.client.Patch(ctx, daemonSet, patch)
}

func (d *BosunDaemonSet) Delete(ctx context.Context, name, namespace string) error {
	daemonSet, err := d.Get(ctx, name, namespace)
	if apiErrors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	err = d.client.Delete(ctx, daemonSet)
	if apiErrors.IsGone(err) {
		return nil
	}
	return err
}

func (d *BosunDaemonSet) StartTracker(ctx context.Context) error {
	handlers := internal.EventHandlers{
		AddFunc:    internal.AddInformerEvent(ctx, daemonSetUpdateEvent),
		UpdateFunc: internal.UpdateInformerEvent(ctx, daemonSetUpdateEvent),
		DeleteFunc: internal.DeleteInformerEvent(ctx, daemonSetDeleteEvent),
	}
	err := d.tracker.Start(ctx, handlers)
	if err != nil {
		return err
	}
	runtime.LogPrint(ctx, "daemonSet informer started")
	return nil
}

func (d *BosunDaemonSet) StopTracker(ctx context.Context) {
	if d.tracker == nil {
		return
	}
	d.tracker.Stop()
	d.tracker = nil
	runtime.LogPrintf(ctx, "daemonSet informer stopped")
}

func (d *BosunDaemonSet) IsTrackerInit() bool {
	if d.tracker == nil {
		return false
	}
	return true
}
