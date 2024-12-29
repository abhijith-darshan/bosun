package pkg

import (
	"bosun/pkg/internal/workloads"
	"errors"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

/************************************************************
 * Pods API
************************************************************/

// StartPodSync starts an Informer on pods for a specific cluster.
func (a *App) StartPodSync(contextID string) error {
	kCtx, ok := a.bosunClusters[contextID]
	if !ok {
		return errors.New("context not found")
	}
	if !kCtx.pod.IsTrackerInit() {
		p := workloads.NewBosunPod(kCtx.client, kCtx.clientSet)
		*a.bosunClusters[contextID].pod = *p
	}
	// start an informer for pods for this cluster
	if err := kCtx.pod.StartTracker(a.ctx); err != nil {
		runtime.LogErrorf(a.ctx, "error starting pod informer: %s", err.Error())
		return err
	}
	return nil
}

func (a *App) ListPods(contextID string) ([]interface{}, error) {
	kCtx, ok := a.bosunClusters[contextID]
	if !ok {
		return nil, errors.New("context not found")
	}
	return kCtx.pod.ListFromStore(), nil
}

func (a *App) StopPodSync(contextID string) error {
	kCtx, ok := a.bosunClusters[contextID]
	if !ok {
		return errors.New("context not found")
	}
	// Stop the informer for pods
	kCtx.pod.StopTracker(a.ctx)
	return nil
}

/************************************************************
 * Deployments API
************************************************************/

// StartDeploymentsSync starts an Informer on deployments for a specific cluster.
func (a *App) StartDeploymentsSync(contextID string) error {
	kCtx, ok := a.bosunClusters[contextID]
	if !ok {
		return errors.New("context not found")
	}
	if !kCtx.deployment.IsTrackerInit() {
		d := workloads.NewBosunDeployment(kCtx.client, kCtx.clientSet)
		*a.bosunClusters[contextID].deployment = *d
	}
	// start an informer for deployments for this cluster
	if err := kCtx.deployment.StartTracker(a.ctx); err != nil {
		runtime.LogErrorf(a.ctx, "error starting pod informer: %s", err.Error())
		return err
	}
	return nil
}

func (a *App) ListDeployments(contextID string) ([]interface{}, error) {
	kCtx, ok := a.bosunClusters[contextID]
	if !ok {
		return nil, errors.New("context not found")
	}
	return kCtx.deployment.ListFromStore(), nil
}

func (a *App) StopDeploymentsSync(contextID string) error {
	kCtx, ok := a.bosunClusters[contextID]
	if !ok {
		return errors.New("context not found")
	}
	// Stop the informer for deployments
	kCtx.deployment.StopTracker(a.ctx)
	return nil
}

/************************************************************
 * Daemon Sets API
************************************************************/

// StartDaemonSetsSync starts an Informer on daemonSets for a specific cluster.
func (a *App) StartDaemonSetsSync(contextID string) error {
	kCtx, ok := a.bosunClusters[contextID]
	if !ok {
		return errors.New("context not found")
	}
	if !kCtx.daemonSet.IsTrackerInit() {
		d := workloads.NewBosunDaemonSet(kCtx.client, kCtx.clientSet)
		*a.bosunClusters[contextID].daemonSet = *d
	}
	// start an informer for daemonSets for this cluster
	if err := kCtx.daemonSet.StartTracker(a.ctx); err != nil {
		runtime.LogErrorf(a.ctx, "error starting pod informer: %s", err.Error())
		return err
	}
	return nil
}

func (a *App) ListDaemonSets(contextID string) ([]interface{}, error) {
	kCtx, ok := a.bosunClusters[contextID]
	if !ok {
		return nil, errors.New("context not found")
	}
	return kCtx.daemonSet.ListFromStore(), nil
}

func (a *App) StopDaemonSetsSync(contextID string) error {
	kCtx, ok := a.bosunClusters[contextID]
	if !ok {
		return errors.New("context not found")
	}
	// Stop the informer for daemonSets
	kCtx.daemonSet.StopTracker(a.ctx)
	return nil
}

/************************************************************
 * Stateful Sets API
************************************************************/

// StartStatefulSetsSync starts an Informer on statefulSets for a specific cluster.
func (a *App) StartStatefulSetsSync(contextID string) error {
	kCtx, ok := a.bosunClusters[contextID]
	if !ok {
		return errors.New("context not found")
	}
	if !kCtx.statefulSet.IsTrackerInit() {
		s := workloads.NewBosunStatefulSet(kCtx.client, kCtx.clientSet)
		*a.bosunClusters[contextID].statefulSet = *s
	}
	// start an informer for statefulSets for this cluster
	if err := kCtx.statefulSet.StartTracker(a.ctx); err != nil {
		runtime.LogErrorf(a.ctx, "error starting pod informer: %s", err.Error())
		return err
	}
	return nil
}

func (a *App) ListStatefulSets(contextID string) ([]interface{}, error) {
	kCtx, ok := a.bosunClusters[contextID]
	if !ok {
		return nil, errors.New("context not found")
	}
	return kCtx.statefulSet.ListFromStore(), nil
}

func (a *App) StopStatefulSetsSync(contextID string) error {
	kCtx, ok := a.bosunClusters[contextID]
	if !ok {
		return errors.New("context not found")
	}
	// Stop the informer for statefulSets
	kCtx.statefulSet.StopTracker(a.ctx)
	return nil
}

/************************************************************
 * Replica Sets API
************************************************************/

// StartReplicaSetsSync starts an Informer on replicaSets for a specific cluster.
func (a *App) StartReplicaSetsSync(contextID string) error {
	kCtx, ok := a.bosunClusters[contextID]
	if !ok {
		return errors.New("context not found")
	}
	if !kCtx.replicaSet.IsTrackerInit() {
		r := workloads.NewBosunReplicaSet(kCtx.client, kCtx.clientSet)
		*a.bosunClusters[contextID].replicaSet = *r
	}
	// start an informer for statefulSets for this cluster
	if err := kCtx.replicaSet.StartTracker(a.ctx); err != nil {
		runtime.LogErrorf(a.ctx, "error starting pod informer: %s", err.Error())
		return err
	}
	return nil
}

func (a *App) ListReplicaSets(contextID string) ([]interface{}, error) {
	kCtx, ok := a.bosunClusters[contextID]
	if !ok {
		return nil, errors.New("context not found")
	}
	return kCtx.replicaSet.ListFromStore(), nil
}

func (a *App) StopReplicaSetsSync(contextID string) error {
	kCtx, ok := a.bosunClusters[contextID]
	if !ok {
		return errors.New("context not found")
	}
	// Stop the informer for replicaSets
	kCtx.replicaSet.StopTracker(a.ctx)
	return nil
}

/************************************************************
 * Jobs API
************************************************************/

// StartJobSync starts an Informer on jobs for a specific cluster.
func (a *App) StartJobSync(contextID string) error {
	kCtx, ok := a.bosunClusters[contextID]
	if !ok {
		return errors.New("context not found")
	}
	if !kCtx.job.IsTrackerInit() {
		j := workloads.NewBosunJob(kCtx.client, kCtx.clientSet)
		*a.bosunClusters[contextID].job = *j
	}
	// start an informer for jobs for this cluster
	if err := kCtx.job.StartTracker(a.ctx); err != nil {
		runtime.LogErrorf(a.ctx, "error starting job informer: %s", err.Error())
		return err
	}
	return nil
}

func (a *App) ListJobs(contextID string) ([]interface{}, error) {
	kCtx, ok := a.bosunClusters[contextID]
	if !ok {
		return nil, errors.New("context not found")
	}
	return kCtx.job.ListFromStore(), nil
}

func (a *App) StopJobSync(contextID string) error {
	kCtx, ok := a.bosunClusters[contextID]
	if !ok {
		return errors.New("context not found")
	}
	// Stop the informer for jobs
	kCtx.job.StopTracker(a.ctx)
	return nil
}

/************************************************************
 * Cron Jobs API
************************************************************/

// StartCronJobSync starts an Informer on cron jobs for a specific cluster.
func (a *App) StartCronJobSync(contextID string) error {
	kCtx, ok := a.bosunClusters[contextID]
	if !ok {
		return errors.New("context not found")
	}
	if !kCtx.cronJob.IsTrackerInit() {
		cj := workloads.NewBosunCronJob(kCtx.client, kCtx.clientSet)
		*a.bosunClusters[contextID].cronJob = *cj
	}
	// start an informer for cron jobs for this cluster
	if err := kCtx.cronJob.StartTracker(a.ctx); err != nil {
		runtime.LogErrorf(a.ctx, "error starting cron job informer: %s", err.Error())
		return err
	}
	return nil
}

func (a *App) ListCronJobs(contextID string) ([]interface{}, error) {
	kCtx, ok := a.bosunClusters[contextID]
	if !ok {
		return nil, errors.New("context not found")
	}
	return kCtx.cronJob.ListFromStore(), nil
}

func (a *App) StopCronJobSync(contextID string) error {
	kCtx, ok := a.bosunClusters[contextID]
	if !ok {
		return errors.New("context not found")
	}
	// Stop the informer for cron jobs
	kCtx.cronJob.StopTracker(a.ctx)
	return nil
}
