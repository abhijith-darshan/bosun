import { entityConfig, SelectEntityId } from '@ngrx/signals/entities';
import { V1DaemonSet, V1Deployment, V1Pod } from '@kubernetes/client-node';
import { type } from '@ngrx/signals';

const podSelect: SelectEntityId<V1Pod> = (pod) => {
    if (pod.metadata?.uid) {
        return pod.metadata.uid;
    }
    return '';
};

const deploymentSelect: SelectEntityId<V1Deployment> = (deployment) => {
    if (deployment.metadata?.uid) {
        return deployment.metadata.uid;
    }
    return '';
};

const daemonSetSelect: SelectEntityId<V1DaemonSet> = (daemonSet) => {
    if (daemonSet.metadata?.uid) {
        return daemonSet.metadata.uid;
    }
    return '';
};

export const PodEntityConfig = entityConfig({
    entity: type<V1Pod>(),
    collection: 'pods',
    selectId: podSelect,
});

export const DeploymentEntityConfig = entityConfig({
    entity: type<V1Deployment>(),
    collection: 'deployments',
    selectId: deploymentSelect,
});

export const DaemonSetEntityConfig = entityConfig({
    entity: type<V1DaemonSet>(),
    collection: 'daemonSets',
    selectId: daemonSetSelect,
});
