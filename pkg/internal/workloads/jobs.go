package workloads

import (
	"bosun/pkg/internal"
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	batchv1 "k8s.io/api/batch/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

const (
	jobUpdateEvent = "jobUpdate"
	jobDeleteEvent = "jobDelete"
)

type BosunJob struct {
	client    client.Client
	clientSet *kubernetes.Clientset
	tracker   internal.InformerManager
}

func NewBosunJob(k8sClient client.Client, clientSet *kubernetes.Clientset) *BosunJob {
	informer := informers.NewSharedInformerFactoryWithOptions(clientSet, time.Minute*10).Batch().V1().Jobs().Informer()
	return &BosunJob{
		client:    k8sClient,
		clientSet: clientSet,
		tracker:   internal.NewInformerManager(informer),
	}
}

func (j *BosunJob) Get(ctx context.Context, name, namespace string) (*batchv1.Job, error) {
	job := &batchv1.Job{}
	err := j.client.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, job)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (j *BosunJob) List(ctx context.Context) ([]batchv1.Job, error) {
	list := &batchv1.JobList{}
	err := j.client.List(ctx, list)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (j *BosunJob) ListFromStore() []interface{} {
	return j.tracker.Store().List()
}

func (j *BosunJob) Update(ctx context.Context, job *batchv1.Job) error {
	oldObj, err := j.Get(ctx, job.GetName(), job.GetNamespace())
	if err != nil {
		return err
	}
	patch := client.MergeFrom(oldObj)
	return j.client.Patch(ctx, job, patch)
}

func (j *BosunJob) Delete(ctx context.Context, name, namespace string) error {
	job, err := j.Get(ctx, name, namespace)
	if apiErrors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	err = j.client.Delete(ctx, job)
	if apiErrors.IsGone(err) {
		return nil
	}
	return err
}

func (j *BosunJob) StartTracker(ctx context.Context) error {
	handlers := internal.EventHandlers{
		AddFunc:    internal.AddInformerEvent(ctx, jobUpdateEvent),
		UpdateFunc: internal.UpdateInformerEvent(ctx, jobUpdateEvent),
		DeleteFunc: internal.DeleteInformerEvent(ctx, jobDeleteEvent),
	}
	err := j.tracker.Start(ctx, handlers)
	if err != nil {
		return err
	}
	runtime.LogPrint(ctx, "job informer started")
	return nil
}

func (j *BosunJob) StopTracker(ctx context.Context) {
	if j.tracker == nil {
		return
	}
	j.tracker.Stop()
	j.tracker = nil
	runtime.LogPrintf(ctx, "job informer stopped")
}

func (j *BosunJob) IsTrackerInit() bool {
	if j.tracker == nil {
		return false
	}
	return true
}
