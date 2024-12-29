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
