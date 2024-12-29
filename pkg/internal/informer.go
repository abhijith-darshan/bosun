package internal

import (
	"context"
	"errors"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"k8s.io/client-go/tools/cache"
)

// InformerManager interface to manage start, stop, and events handling
type InformerManager interface {
	Start(context.Context, EventHandlers) error
	Stop()
	Store() cache.Store
}

// EventHandlers struct to store Add, update, and delete event handler functions
type EventHandlers struct {
	AddFunc       func(obj interface{})
	UpdateFunc    func(oldObj, newObj interface{})
	DeleteFunc    func(obj interface{})
	TransformFunc func(interface{}) (interface{}, error)
}

// genericInformer handles any type of informer
type genericInformer struct {
	informer   cache.SharedIndexInformer
	stopChan   chan struct{}
	hasStarted bool
}

// NewInformerManager creates and returns an InformerManager for the given informer and event handlers
func NewInformerManager(informer cache.SharedIndexInformer) InformerManager {
	return &genericInformer{
		informer: informer,
		stopChan: make(chan struct{}),
	}
}

// Start starts the informer and listens to events
func (g *genericInformer) Start(ctx context.Context, handlers EventHandlers) error {
	if g.hasStarted {
		runtime.LogPrintf(ctx, "informer already started")
		return nil
	}
	// Register the event handlers to the informer when it starts
	_, err := g.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    handlers.AddFunc,
		UpdateFunc: handlers.UpdateFunc,
		DeleteFunc: handlers.DeleteFunc,
	})
	if err != nil {
		runtime.LogErrorf(ctx, "informer add event handler error: %s", err.Error())
		return err
	}

	if handlers.TransformFunc != nil {
		err = g.informer.SetTransform(handlers.TransformFunc)
		if err != nil {
			runtime.LogErrorf(ctx, "informer set transform error: %s", err.Error())
			return err
		}
	}

	// start the informer in a separate goroutine
	go g.informer.Run(g.stopChan)

	// Wait for cache sync
	if !cache.WaitForCacheSync(ctx.Done(), g.informer.HasSynced) {
		return errors.New("timed out waiting for cache to sync")
	}
	g.hasStarted = true

	return nil
}

func (g *genericInformer) Store() cache.Store {
	return g.informer.GetStore()
}

// Stop stops the informer by closing the stop channel
func (g *genericInformer) Stop() {
	close(g.stopChan)
	g.hasStarted = false
}
