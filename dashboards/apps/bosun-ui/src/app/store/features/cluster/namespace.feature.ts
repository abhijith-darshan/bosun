import { patchState, signalStoreFeature, withMethods, withState } from '@ngrx/signals';
import { addEntities, removeAllEntities, removeEntity, setEntity, withEntities } from '@ngrx/signals/entities';
import { inject } from '@angular/core';
import { NamespaceService } from '../../services';
import { SortKubernetesResourcesByName } from '../../../utils';
import { V1Namespace } from '@kubernetes/client-node';
import { NamespaceEntityConfig } from './selectors';
import { BosunNamespaceState } from '../../models';

const initialNamespaceState: BosunNamespaceState = {
    selectedNamespace: null,
};

export function withNamespaces() {
    return signalStoreFeature(
        withState(initialNamespaceState),
        withEntities(NamespaceEntityConfig),
        withMethods((store, namespaceService = inject(NamespaceService)) => ({
            async listNamespaces(contextId: string) {
                let namespaces = await namespaceService.listNamespaces(contextId);
                namespaces = SortKubernetesResourcesByName(namespaces);
                patchState(store, addEntities(namespaces, NamespaceEntityConfig));
            },
            addUpdateNamespace(namespace: V1Namespace) {
                patchState(store, setEntity(namespace, NamespaceEntityConfig));
            },
            deleteNamespace(namespace: V1Namespace) {
                if (namespace.metadata?.name) {
                    patchState(store, removeEntity(namespace.metadata.name, NamespaceEntityConfig));
                }
            },
            async stopNamespaceTracking(contextId: string): Promise<void> {
                try {
                    await namespaceService.stopNamespaceSync(contextId);
                    patchState(store, removeAllEntities(NamespaceEntityConfig));
                } catch (e) {
                    console.error('Failed to stop tracking grouped: ', e);
                }
            },
            async startNamespaceTracking(contextId: string) {
                try {
                    await namespaceService.startNamespaceSync(contextId);
                } catch (e) {
                    console.error('failed to start tracking grouped: ', e);
                }
            },
            setSelectedNamespace(namespace: V1Namespace) {
                patchState(store, { selectedNamespace: { ...namespace } });
            },
        })),
    );
}
