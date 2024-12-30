import { patchState, signalStore, withHooks, withMethods, withState } from '@ngrx/signals';
import { AppService } from './services';
import { inject } from '@angular/core';
import { DynamicDialogRef } from 'primeng/dynamicdialog';
import { Bosun, BosunAppState } from './models';
import { withNamespaces, withWorkloads } from './features';

const initialBosunAppState: BosunAppState = {
    currentContext: { id: '', name: '', shortName: '', version: '', knownResources: [] },
    clusters: [],
    isClustersLoading: true,
    clustersLoadingErrored: false,
    isClusterLoginLoading: true,
    clusterLoginErrored: false,
    clusterLoginError: '',
};

export const AppStore = signalStore(
    { providedIn: 'root', protectedState: false },
    withState(initialBosunAppState),
    withNamespaces(),
    withWorkloads(),
    withMethods((store, appService = inject(AppService)) => ({
        async getClusters(): Promise<void> {
            try {
                patchState(store, { isClustersLoading: true });
                const clusters = await appService.readKubeConfigs();
                patchState(store, { clusters: [...clusters], isClustersLoading: false, clustersLoadingErrored: false });
            } catch (e) {
                patchState(store, {
                    clustersLoadingErrored: true,
                    clustersLoadingError: e?.toString(),
                });
            }
        },
        async loginToCluster(cluster: Bosun, ref: DynamicDialogRef): Promise<void> {
            try {
                patchState(store, {
                    currentContext: cluster,
                    isClusterLoginLoading: true,
                    clusterLoginErrored: false,
                    clusterLoginError: '',
                });
                const version = await appService.clusterLogin(cluster.id);
                const clusters = store.clusters().map((c) => {
                    if (c.id === cluster.id) {
                        patchState(store, { currentContext: { ...cluster, version } });
                        return { ...c, version };
                    }
                    return c;
                });
                /***
                 * performance bottleneck can be observed here for clusters with large number of namespaces
                 ***/
                // get known resources
                await this.listKnownClusterResources();
                // start tracking the namespaces
                await store.startNamespaceTracking(cluster.id);
                // list namespaces
                await store.listNamespaces(cluster.id);
                patchState(store, { clusters, isClusterLoginLoading: false });
                ref.close(true);
            } catch (e: any) {
                patchState(store, {
                    isClusterLoginLoading: false,
                    clusterLoginErrored: true,
                    clusterLoginError: e?.message,
                });
            }
        },
        async listKnownClusterResources(): Promise<void> {
            try {
                const resources = await appService.listKnownResources(store.currentContext().id);
                patchState(store, { currentContext: { ...store.currentContext(), knownResources: [...resources] } });
            } catch (e: any) {
                console.error('Failed to get known resources:', e);
            }
        },
    })),
    withHooks({
        onInit(store) {
            store.getClusters();
        },
    }),
);
