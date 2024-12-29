package workloads

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
	podUpdateEvent = "podUpdate"
	podDeleteEvent = "podDelete"
)

type BosunPod struct {
	client    client.Client
	clientSet *kubernetes.Clientset
	tracker   internal.InformerManager
}

func NewBosunPod(k8sClient client.Client, clientSet *kubernetes.Clientset) *BosunPod {
	informer := informers.NewSharedInformerFactoryWithOptions(clientSet, time.Minute*10).Core().V1().Pods().Informer()
	return &BosunPod{
		client:    k8sClient,
		clientSet: clientSet,
		tracker:   internal.NewInformerManager(informer),
	}
}

func (p *BosunPod) Get(ctx context.Context, name, namespace string) (*corev1.Pod, error) {
	pod := &corev1.Pod{}
	err := p.client.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, pod)
	if err != nil {
		return nil, err
	}
	return pod, nil
}

func (p *BosunPod) List(ctx context.Context) ([]corev1.Pod, error) {
	list := &corev1.PodList{}
	err := p.client.List(ctx, list)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (p *BosunPod) ListFromStore() []interface{} {
	return p.tracker.Store().List()
}

func (p *BosunPod) Update(ctx context.Context, pod *corev1.Pod) error {
	oldObj, err := p.Get(ctx, pod.GetName(), pod.GetNamespace())
	if err != nil {
		return err
	}
	patch := client.MergeFrom(oldObj)
	return p.client.Patch(ctx, pod, patch)
}

func (p *BosunPod) Delete(ctx context.Context, name, namespace string) error {
	pod, err := p.Get(ctx, name, namespace)
	if apiErrors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	err = p.client.Delete(ctx, pod)
	if apiErrors.IsGone(err) {
		return nil
	}
	return err
}

func (p *BosunPod) StartTracker(ctx context.Context) error {
	handlers := internal.EventHandlers{
		AddFunc:    internal.AddInformerEvent(ctx, podUpdateEvent),
		UpdateFunc: internal.UpdateInformerEvent(ctx, podUpdateEvent),
		DeleteFunc: internal.DeleteInformerEvent(ctx, podDeleteEvent),
	}
	err := p.tracker.Start(ctx, handlers)
	if err != nil {
		return err
	}
	runtime.LogPrint(ctx, "pod informer started")
	return nil
}

func (p *BosunPod) StopTracker(ctx context.Context) {
	if p.tracker == nil {
		return
	}
	p.tracker.Stop()
	p.tracker = nil
	runtime.LogPrintf(ctx, "pod informer stopped")
}

func (p *BosunPod) IsTrackerInit() bool {
	if p.tracker == nil {
		return false
	}
	return true
}
