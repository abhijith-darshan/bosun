import { Table } from 'primeng/table';
import { V1Namespace } from '@kubernetes/client-node';

export enum ResourceEvents {
    NamespaceUpdate = 'namespaceUpdate',
    NamespaceDelete = 'namespaceDelete',
    PodUpdate = 'podUpdate',
    PodDelete = 'podDelete',
    DeploymentUpdate = 'deploymentUpdate',
    DeploymentDelete = 'deploymentDelete',
    DaemonSetUpdate = 'daemonSetUpdate',
    DaemonSetDelete = 'daemonSetDelete',
}

const groupOrder = ['Nodes', 'Workloads', 'Config'];

// Define the custom order
const workloadOrder = ['Pods', 'Deployments', 'DaemonSets', 'StatefulSets', 'ReplicaSets', 'Jobs', 'CronJobs'];

// Create a map that stores the index for each label in the custom order
export const WorkloadOrderMap: Record<string, number> = workloadOrder.reduce(
    (acc, label, index) => {
        acc[label] = index;
        return acc;
    },
    {} as Record<string, number>,
);

export const GroupOrderMap: Record<string, number> = groupOrder.reduce(
    (acc, label, index) => {
        acc[label] = index;
        return acc;
    },
    {} as Record<string, number>,
);

export function SortKubernetesResourcesByName(objects: any[]): any[] {
    objects.sort((a, b) => {
        if (a.metadata?.name === undefined && b.metadata?.name === undefined) return 0;
        if (a.metadata?.name === undefined) return 1;
        if (b.metadata?.name === undefined) return -1;
        return a.metadata.name.localeCompare(b.metadata.name);
    });
    return objects;
}

export function filterTableByNamespace(ns: V1Namespace, dt: Table) {
    if (ns) {
        dt.filterGlobal(ns.metadata?.name, 'equals');
    } else {
        dt.reset();
    }
}
