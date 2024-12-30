import { ResolveFn } from '@angular/router';
import { inject } from '@angular/core';
import { AppStore } from './app.store';

/*export const InitiateNamespacesTracking: ResolveFn<any> = async () => {
    const store = inject(AppStore);
    await store.listKnownClusterResources();
    store.startNamespaceTracking(store.currentContext().id).then();
};*/

export const InitiatePodsTracking: ResolveFn<any> = async () => {
    const store = inject(AppStore);
    await store.startPodTracking(store.currentContext().id);
};

export const InitiateDeploymentsTracking: ResolveFn<any> = async () => {
    const store = inject(AppStore);
    await store.startDeploymentTracking(store.currentContext().id);
};

export const InitiateDaemonSetsTracking: ResolveFn<any> = async () => {
    const store = inject(AppStore);
    await store.startDaemonSetsTracking(store.currentContext().id);
};
