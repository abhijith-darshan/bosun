package pkg

import (
	"bosun/pkg/internal/cluster"
	"bosun/pkg/internal/workloads"
	"context"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// App struct
type App struct {
	ctx           context.Context
	bosunClusters map[string]BosunCluster
}

type BosunCluster struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ShortName   string `json:"shortName"`
	Version     string `json:"version"`
	client      client.Client
	clientSet   *kubernetes.Clientset
	namespace   *cluster.BosunNamespace
	pod         *workloads.BosunPod
	deployment  *workloads.BosunDeployment
	daemonSet   *workloads.BosunDaemonSet
	statefulSet *workloads.BosunStatefulSet
	replicaSet  *workloads.BosunReplicaSet
	job         *workloads.BosunJob
	cronJob     *workloads.BosunCronJob
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}
