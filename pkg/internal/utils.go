package internal

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	CategoryWorkload   = "Workloads"
	CategoryConfig     = "Config"
	CategoryNamespaces = "Namespaces"
	CategoryNodes      = "Nodes"
)

const (
	GroupApps                  = "apps/v1"
	GroupBatch                 = "batch/v1"
	GroupV1                    = "v1"
	GroupAdmissionRegistration = "admissionregistration.k8s.io/v1"
)

type BosunKind string

const (
	KindNodes                          BosunKind = "Node"
	KindNamespaces                     BosunKind = "Namespace"
	KindDeployment                     BosunKind = "Deployment"
	KindControllerRevision             BosunKind = "ControllerRevision"
	KindDaemonSet                      BosunKind = "DaemonSet"
	KindReplicaSet                     BosunKind = "ReplicaSet"
	KindStatefulSet                    BosunKind = "StatefulSet"
	KindPod                            BosunKind = "Pod"
	KindCronJob                        BosunKind = "CronJob"
	KindJob                            BosunKind = "Job"
	KindSecret                         BosunKind = "Secret"
	KindConfigMap                      BosunKind = "ConfigMap"
	KindMutatingWebhookConfiguration   BosunKind = "MutatingWebhookConfiguration"
	KindValidatingWebhookConfiguration BosunKind = "ValidatingWebhookConfiguration"
)

func GetDisplayName(kind BosunKind) string {
	switch kind {
	case KindNodes:
		return "Nodes"
	case KindNamespaces:
		return "Namespaces"
	case KindDeployment:
		return "Deployments"
	case KindControllerRevision:
		return "ControllerRevisions"
	case KindDaemonSet:
		return "DaemonSets"
	case KindReplicaSet:
		return "ReplicaSets"
	case KindStatefulSet:
		return "StatefulSets"
	case KindPod:
		return "Pods"
	case KindCronJob:
		return "CronJobs"
	case KindJob:
		return "Jobs"
	case KindSecret:
		return "Secrets"
	case KindConfigMap:
		return "ConfigMaps"
	case KindMutatingWebhookConfiguration:
		return "Mutating Webhook Configs"
	case KindValidatingWebhookConfiguration:
		return "Validating Webhook Configs"
	default:
		return string(kind)
	}
}

func AddInformerEvent(ctx context.Context, eventName string) func(obj interface{}) {
	return func(obj interface{}) {
		runtime.EventsEmit(ctx, eventName, obj)
	}
}

func UpdateInformerEvent(ctx context.Context, eventName string) func(oldObj, newObj interface{}) {
	return func(oldObj, newObj interface{}) {
		runtime.EventsEmit(ctx, eventName, newObj)
	}
}

func DeleteInformerEvent(ctx context.Context, eventName string) func(obj interface{}) {
	return func(obj interface{}) {
		runtime.EventsEmit(ctx, eventName, obj)
	}
}
