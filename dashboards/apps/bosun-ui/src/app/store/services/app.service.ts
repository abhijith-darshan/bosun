import { Injectable } from '@angular/core';
import { GetKnownResources, GetVersion, ReadKubeConfigs } from '../../../../wailsjs/go/pkg/App';
import { Bosun } from '../models';
import { pkg } from '../../../../wailsjs/go/models';
import { SelectItemGroup } from 'primeng/api';
import { GroupOrderMap, WorkloadOrderMap } from '../../utils';
import Resource = pkg.Resource;

@Injectable({
    providedIn: 'root',
})
export class AppService {
    constructor() {}

    async readKubeConfigs(): Promise<Bosun[]> {
        const bosunClusters = (await ReadKubeConfigs()) as Bosun[];
        for (const cluster of bosunClusters) {
            cluster.knownResources = [];
        }
        return [...bosunClusters];
    }

    async clusterLogin(contextId: string): Promise<string> {
        try {
            return await GetVersion(contextId);
        } catch (e: any) {
            throw new Error(e.toString());
        }
    }

    async listKnownResources(contextId: string): Promise<SelectItemGroup<Resource>[]> {
        try {
            const resources = await GetKnownResources(contextId);
            return this.applyKnownResourceTransformation(resources);
        } catch (e: any) {
            throw new Error(e.toString());
        }
    }

    private applyKnownResourceTransformation(resources: Resource[]): SelectItemGroup<Resource>[] {
        const groupedResourcesMap = new Map<string, SelectItemGroup<Resource>>();

        // Group resources by their `key` (group key)
        resources.forEach((resource) => {
            const groupKey = resource.key;

            // If the group does not exist, create it
            if (!groupedResourcesMap.has(groupKey)) {
                groupedResourcesMap.set(groupKey, {
                    label: groupKey,
                    value: groupKey.toLowerCase(),
                    items: [],
                });
            }

            // Get the group from the map and add the resource's kind and the resource itself
            const group = groupedResourcesMap.get(groupKey);
            if (group && !this.excludeItemForGroup(groupKey)) {
                group.items.push({
                    label: resource.displayName,
                    value: resource,
                });
            }
        });

        // Return the grouped resources as a dictionary
        const groups = Array.from(groupedResourcesMap.values()).sort((a, b) => GroupOrderMap[a.label] - GroupOrderMap[b.label]);
        const workloadIndex = groups.findIndex((g) => g.label === 'Workloads');
        groups[workloadIndex].items.sort((a, b) => {
            // Ensure we are working with valid strings and not undefined
            const aLabel = a.label ?? ''; // Default to empty string if undefined
            const bLabel = b.label ?? ''; // Default to empty string if undefined

            return WorkloadOrderMap[aLabel] - WorkloadOrderMap[bLabel];
        });
        return groups;
    }

    private excludeItemForGroup(groupKey: string): boolean {
        return groupKey === 'Namespaces' || groupKey === 'Nodes';
    }
}
