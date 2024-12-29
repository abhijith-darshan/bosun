package pkg

import (
	"bosun/pkg/internal/cluster"
	"errors"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

/************************************************************
 * Namespaces API
************************************************************/

// StartNamespaceSync starts an Informer on namespaces for a specific cluster.
func (a *App) StartNamespaceSync(contextID string) error {
	kCtx, ok := a.bosunClusters[contextID]
	if !ok {
		return errors.New("context not found")
	}

	if !kCtx.namespace.IsTrackerInit() {
		ns := cluster.NewBosunNamespace(kCtx.client, kCtx.clientSet)
		*a.bosunClusters[contextID].namespace = *ns
	}
	// start an informer for namespaces for this cluster
	if err := kCtx.namespace.StartTracker(a.ctx); err != nil {
		runtime.LogErrorf(a.ctx, "error starting namespace informer: %s", err.Error())
		return err
	}
	return nil
}

func (a *App) ListNamespaces(contextID string) ([]interface{}, error) {
	kCtx, ok := a.bosunClusters[contextID]
	if !ok {
		return nil, errors.New("context not found")
	}
	return kCtx.namespace.ListFromStore(), nil
}

func (a *App) StopNamespaceSync(contextID string) error {
	kCtx, ok := a.bosunClusters[contextID]
	if !ok {
		return errors.New("context not found")
	}
	// Stop the informer for namespaces
	kCtx.namespace.StopTracker(a.ctx)
	return nil
}
