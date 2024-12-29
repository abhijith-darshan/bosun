package internal

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

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
