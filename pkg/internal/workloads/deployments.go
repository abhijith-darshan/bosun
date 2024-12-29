package workloads

import (
	"bosun/pkg/internal"
	"context"
	"errors"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	appsv1 "k8s.io/api/apps/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"slices"
	"time"
)

const (
	deploymentUpdateEvent = "deploymentUpdate"
	deploymentDeleteEvent = "deploymentDelete"
)

type BosunDeployment struct {
	client    client.Client
	clientSet *kubernetes.Clientset
	tracker   internal.InformerManager
}

func NewBosunDeployment(k8sClient client.Client, clientSet *kubernetes.Clientset) *BosunDeployment {
	informer := informers.NewSharedInformerFactoryWithOptions(clientSet, time.Minute*10).Apps().V1().Deployments().Informer()
	return &BosunDeployment{
		client:    k8sClient,
		clientSet: clientSet,
		tracker:   internal.NewInformerManager(informer),
	}
}

func (d *BosunDeployment) Get(ctx context.Context, name, namespace string) (*appsv1.Deployment, error) {
	deployment := &appsv1.Deployment{}
	err := d.client.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, deployment)
	if err != nil {
		return nil, err
	}
	return deployment, nil
}

func (d *BosunDeployment) List(ctx context.Context) ([]appsv1.Deployment, error) {
	list := &appsv1.DeploymentList{}
	err := d.client.List(ctx, list)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (d *BosunDeployment) ListFromStore() []interface{} {
	return d.tracker.Store().List()
}

func (d *BosunDeployment) Update(ctx context.Context, deployment *appsv1.Deployment) error {
	oldObj, err := d.Get(ctx, deployment.GetName(), deployment.GetNamespace())
	if err != nil {
		return err
	}
	patch := client.MergeFrom(oldObj)
	return d.client.Patch(ctx, deployment, patch)
}

func (d *BosunDeployment) Delete(ctx context.Context, name, namespace string) error {
	deployment, err := d.Get(ctx, name, namespace)
	if apiErrors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	err = d.client.Delete(ctx, deployment)
	if apiErrors.IsGone(err) {
		return nil
	}
	return err
}

func (d *BosunDeployment) StartTracker(ctx context.Context) error {
	handlers := internal.EventHandlers{
		AddFunc:       internal.AddInformerEvent(ctx, deploymentUpdateEvent),
		UpdateFunc:    internal.UpdateInformerEvent(ctx, deploymentUpdateEvent),
		DeleteFunc:    internal.DeleteInformerEvent(ctx, deploymentDeleteEvent),
		TransformFunc: transformDeployment,
	}
	err := d.tracker.Start(ctx, handlers)
	if err != nil {
		return err
	}
	runtime.LogPrint(ctx, "deployment informer started")
	return nil
}

func (d *BosunDeployment) StopTracker(ctx context.Context) {
	if d.tracker == nil {
		return
	}
	d.tracker.Stop()
	d.tracker = nil
	runtime.LogPrintf(ctx, "deployment informer stopped")
}

func (d *BosunDeployment) IsTrackerInit() bool {
	if d.tracker == nil {
		return false
	}
	return true
}

func transformDeployment(obj interface{}) (interface{}, error) {
	deployment, ok := obj.(*appsv1.Deployment)
	if !ok {
		return nil, errors.New("failed to assert deployment")
	}
	conditions := deployment.Status.Conditions
	if len(conditions) == 0 {
		return deployment, nil
	}
	slices.SortFunc(conditions, func(a appsv1.DeploymentCondition, b appsv1.DeploymentCondition) int {
		if a.Type < b.Type {
			return -1
		}
		if a.Type > b.Type {
			return 1
		}
		return 0
	})
	return deployment, nil
}
