import { patchState, signalStoreFeature, withMethods } from '@ngrx/signals';
import { addEntities, removeAllEntities, removeEntity, setEntity, withEntities } from '@ngrx/signals/entities';
import { inject } from '@angular/core';
import { WorkloadService } from '../../services';
import { V1DaemonSet, V1Deployment, V1Pod } from '@kubernetes/client-node';
import { SortKubernetesResourcesByName } from '../../../utils';
import { DaemonSetEntityConfig, DeploymentEntityConfig, PodEntityConfig } from './selectors';

export function withWorkloads() {
    return signalStoreFeature(
        withEntities(DeploymentEntityConfig),
        withEntities(PodEntityConfig),
        withEntities(DaemonSetEntityConfig),
        withMethods((store, workloadService = inject(WorkloadService)) => ({
            async listPods(contextId: string) {
                let pods = await workloadService.listPods(contextId);
                pods = SortKubernetesResourcesByName(pods);
                patchState(store, addEntities(pods, PodEntityConfig));
            },
            addUpdatePod(pod: V1Pod) {
                patchState(store, setEntity(pod, PodEntityConfig));
            },
            deletePod(pod: V1Pod) {
                if (pod.metadata?.uid) {
                    patchState(store, removeEntity(pod.metadata.uid, PodEntityConfig));
                }
            },
            async stopPodTracking(contextId: string): Promise<void> {
                try {
                    await workloadService.stopPodSync(contextId);
                    patchState(store, removeAllEntities(PodEntityConfig));
                } catch (e) {
                    console.error('Failed to stop tracking grouped: ', e);
                }
            },
            async startPodTracking(contextId: string) {
                try {
                    await workloadService.startPodSync(contextId);
                } catch (e) {
                    console.error('failed to start tracking grouped: ', e);
                }
            },
            async listDeployments(contextId: string) {
                let deployments = await workloadService.listDeployments(contextId);
                deployments = SortKubernetesResourcesByName(deployments);
                patchState(store, addEntities(deployments, DeploymentEntityConfig));
            },
            addUpdateDeployment(deployment: V1Deployment) {
                patchState(store, setEntity(deployment, DeploymentEntityConfig));
            },
            deleteDeployment(deployment: V1Deployment) {
                if (deployment.metadata?.uid) {
                    patchState(store, removeEntity(deployment.metadata.uid, DeploymentEntityConfig));
                }
            },
            async stopDeploymentTracking(contextId: string): Promise<void> {
                try {
                    await workloadService.stopDeploymentSync(contextId);
                    patchState(store, removeAllEntities(DeploymentEntityConfig));
                } catch (e) {
                    console.error('Failed to stop tracking grouped: ', e);
                }
            },
            async startDeploymentTracking(contextId: string) {
                try {
                    await workloadService.startDeploymentSync(contextId);
                } catch (e) {
                    console.error('failed to start tracking grouped: ', e);
                }
            },
            async listDaemonSets(contextId: string) {
                let daemonSets = await workloadService.listDaemonSets(contextId);
                daemonSets = SortKubernetesResourcesByName(daemonSets);
                patchState(store, addEntities(daemonSets, DaemonSetEntityConfig));
            },
            addUpdateDaemonSet(daemonSet: V1DaemonSet) {
                patchState(store, setEntity(daemonSet, DaemonSetEntityConfig));
            },
            deleteDaemonSet(daemonSet: V1DaemonSet) {
                if (daemonSet.metadata?.uid) {
                    patchState(store, removeEntity(daemonSet.metadata.uid, DeploymentEntityConfig));
                }
            },
            async stopDaemonSetTracking(contextId: string): Promise<void> {
                try {
                    await workloadService.stopDaemonSync(contextId);
                    patchState(store, removeAllEntities(DaemonSetEntityConfig));
                } catch (e) {
                    console.error('Failed to stop tracking grouped: ', e);
                }
            },
            async startDaemonSetsTracking(contextId: string) {
                try {
                    await workloadService.startDaemonSync(contextId);
                } catch (e) {
                    console.error('failed to start tracking grouped: ', e);
                }
            },
        })),
    );
}
