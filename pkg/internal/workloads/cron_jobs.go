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
	cronJobUpdateEvent = "cronJobUpdate"
	cronJobDeleteEvent = "cronJobDelete"
)

type BosunCronJob struct {
	client    client.Client
	clientSet *kubernetes.Clientset
	tracker   internal.InformerManager
}

func NewBosunCronJob(k8sClient client.Client, clientSet *kubernetes.Clientset) *BosunCronJob {
	informer := informers.NewSharedInformerFactoryWithOptions(clientSet, time.Minute*10).Batch().V1().CronJobs().Informer()
	return &BosunCronJob{
		client:    k8sClient,
		clientSet: clientSet,
		tracker:   internal.NewInformerManager(informer),
	}
}

func (cj *BosunCronJob) Get(ctx context.Context, name, namespace string) (*batchv1.CronJob, error) {
	cJob := &batchv1.CronJob{}
	err := cj.client.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, cJob)
	if err != nil {
		return nil, err
	}
	return cJob, nil
}

func (cj *BosunCronJob) List(ctx context.Context) ([]batchv1.CronJob, error) {
	list := &batchv1.CronJobList{}
	err := cj.client.List(ctx, list)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (cj *BosunCronJob) ListFromStore() []interface{} {
	return cj.tracker.Store().List()
}

func (cj *BosunCronJob) Update(ctx context.Context, cronJob *batchv1.CronJob) error {
	oldObj, err := cj.Get(ctx, cronJob.GetName(), cronJob.GetNamespace())
	if err != nil {
		return err
	}
	patch := client.MergeFrom(oldObj)
	return cj.client.Patch(ctx, cronJob, patch)
}

func (cj *BosunCronJob) Delete(ctx context.Context, name, namespace string) error {
	cronJob, err := cj.Get(ctx, name, namespace)
	if apiErrors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	err = cj.client.Delete(ctx, cronJob)
	if apiErrors.IsGone(err) {
		return nil
	}
	return err
}

func (cj *BosunCronJob) StartTracker(ctx context.Context) error {
	handlers := internal.EventHandlers{
		AddFunc:    internal.AddInformerEvent(ctx, cronJobUpdateEvent),
		UpdateFunc: internal.UpdateInformerEvent(ctx, cronJobUpdateEvent),
		DeleteFunc: internal.DeleteInformerEvent(ctx, cronJobDeleteEvent),
	}
	err := cj.tracker.Start(ctx, handlers)
	if err != nil {
		return err
	}
	runtime.LogPrint(ctx, "cron job informer started")
	return nil
}

func (cj *BosunCronJob) StopTracker(ctx context.Context) {
	if cj.tracker == nil {
		return
	}
	cj.tracker.Stop()
	cj.tracker = nil
	runtime.LogPrintf(ctx, "cron job informer stopped")
}

func (cj *BosunCronJob) IsTrackerInit() bool {
	if cj.tracker == nil {
		return false
	}
	return true
}
